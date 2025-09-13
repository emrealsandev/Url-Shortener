package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct{ Rdb *redis.Client }

func NewRedis(addr, pass string, db int) *Redis {
	return &Redis{Rdb: redis.NewClient(&redis.Options{Addr: addr, Password: pass, DB: db})}
}

func (c *Redis) GetURLByCode(ctx context.Context, code string) (string, bool, error) {
	v, err := c.Rdb.Get(ctx, "c:"+code).Result()
	if err == redis.Nil {
		return "", false, nil
	}
	return v, err == nil, err
}
func (c *Redis) SetURLByCode(ctx context.Context, code, target string, ttl time.Duration) error {
	return c.Rdb.Set(ctx, "c:"+code, target, ttl).Err()
}
func (c *Redis) DelURLByCode(ctx context.Context, code string) error {
	return c.Rdb.Del(ctx, "c:"+code).Err()
}

func (c *Redis) GetCodeByURLKey(ctx context.Context, urlKey string) (string, bool, error) {
	v, err := c.Rdb.Get(ctx, "u:"+urlKey).Result()
	if err == redis.Nil {
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
