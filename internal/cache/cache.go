package cache

import (
	"context"
	"time"
)

type Cache interface {
	GetURLByCode(ctx context.Context, code string) (string, bool, error)
	SetURLByCode(ctx context.Context, code, target string, ttl time.Duration) error
	DelURLByCode(ctx context.Context, code string) error

	GetCodeByURLKey(ctx context.Context, urlKey string) (string, bool, error)
	SetCodeByURLKey(ctx context.Context, urlKey, code string, ttl time.Duration) error
	IsKeyExists(ctx context.Context, key string) int64
	GetHash(hashKey string, dest any) error
	SetHash(hashKey string, src any, ttl int16) error
}
