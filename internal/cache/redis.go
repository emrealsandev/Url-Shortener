package cache

import (
	"context"
	"errors"
	"time"
	"url-shortener/internal/repo"
	"url-shortener/pkg/utils"

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

func (c *Redis) GetHash(hashKey string) *repo.Settings {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	settingsHash := c.Rdb.HGetAll(ctx, hashKey)
	if settingsHash.Err() != nil || len(settingsHash.Val()) == 0 {
		return nil
	}

	settings := repo.Settings{}
	err := utils.MapToStruct(settingsHash.Val(), &settings)
	if err != nil {
		return nil
	}
	return &settings
}

func (c *Redis) SetHash(hashKey string, src interface{}) error {
	if src == nil {
		return nil
	}
	m, err := utils.StructToMap(src)
	if err != nil {
		return err
	}
	if len(m) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pipe := c.Rdb.TxPipeline()
	pipe.HSet(ctx, hashKey, m)
	pipe.Expire(ctx, hashKey, 5*time.Minute)

	_, err = pipe.Exec(ctx)
	return err
}
