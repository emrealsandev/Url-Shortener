package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"
	"url-shortener/pkg/utils"

	"github.com/davecgh/go-spew/spew"
	"github.com/redis/go-redis/v9"
)

type Redis struct{ Rdb *redis.Client }

func NewRedis(addr, pass string, db int) *Redis {
	return &Redis{Rdb: redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: db})}
}

func (c *Redis) GetURLByCode(ctx context.Context, code string) (string, bool, error) {
	v, err := c.Rdb.Get(ctx, "c:"+code).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	return v, err != nil, err
}
func (c *Redis) SetURLByCode(ctx context.Context, code, target string, ttl time.Duration) error {
	return c.Rdb.Set(ctx, "c:"+code, target, ttl).Err()
}
func (c *Redis) DelURLByCode(ctx context.Context, code string) error {
	return c.Rdb.Del(ctx, "c:"+code).Err()
}

func (c *Redis) GetCodeByURLKey(ctx context.Context, urlKey string) (string, bool, error) {
	v, err := c.Rdb.Get(ctx, "u:"+urlKey).Result()
	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	return v, err != nil, err
}
func (c *Redis) SetCodeByURLKey(ctx context.Context, urlKey, code string, ttl time.Duration) error {
	return c.Rdb.Set(ctx, "u:"+urlKey, code, ttl).Err()
}

func (c *Redis) IsKeyExists(ctx context.Context, key string) int64 {
	return c.Rdb.Exists(ctx, key).Val()
}

func (c *Redis) GetHash(hashKey string, dest any) error {
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return errors.New("dest must be a non-nil pointer")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	hash := c.Rdb.HGetAll(ctx, hashKey)
	if hash.Err() != nil || len(hash.Val()) == 0 {
		return redis.Nil
	}

	elem := v.Elem()

	switch elem.Kind() {
	case reflect.Struct:
		// Hedef bir struct ise, map'i struct'a doldur
		return utils.MapToStruct(hash.Val(), dest)

	case reflect.Map:
		// Hedef bir map ise, türünü kontrol et ve doldur
		mapType := elem.Type()
		// Sadece map[string]string türünü destekleyelim çünkü HGetAll bunu döner.
		if mapType.Key().Kind() != reflect.String || mapType.Elem().Kind() != reflect.String {
			return fmt.Errorf("destination map must be of type map[string]string, but got %s", mapType)
		}

		// Gelen Redis map'ini doğrudan ata
		// reflect.ValueOf ile hashVal'i bir reflect.Value'ya çeviriyoruz.
		elem.Set(reflect.ValueOf(hash.Val()))
		return nil

	default:
		return fmt.Errorf("unsupported destination type: dest must be a pointer to a struct or a map[string]string")
	}
}

func (c *Redis) SetHash(hashKey string, src any, ttl int16) error {
	if src == nil {
		return errors.New("source cannot be nil")
	}

	var dataMap map[string]any

	val := reflect.ValueOf(src)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		m, err := utils.StructToMap(src)
		if err != nil {
			return fmt.Errorf("failed to convert struct to map: %w", err)
		}
		dataMap = m

	case reflect.Map:
		// Eğer src zaten bir map ise, türünü kontrol edip kullanalım
		// HSet için `map[string]any` formatına getirmek en güvenlisi.
		convertedMap := make(map[string]any)
		iter := val.MapRange()
		for iter.Next() {
			key, ok := iter.Key().Interface().(string)
			if !ok {
				return errors.New("map keys must be of type string")
			}
			convertedMap[key] = iter.Value().Interface()
		}
		dataMap = convertedMap

	default:
		return fmt.Errorf("unsupported source type: src must be a struct or a map, got %T", src)
	}

	if len(dataMap) == 0 {
		return nil // Veya errors.New("source data is empty")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pipe := c.Rdb.TxPipeline()
	pipe.HSet(ctx, hashKey, dataMap)
	if ttl > 0 {
		spew.Dump(ttl)
		ttlDuration := time.Duration(ttl) * time.Minute
		pipe.Expire(ctx, hashKey, ttlDuration)
	}
	_, err := pipe.Exec(ctx)
	return err
}
