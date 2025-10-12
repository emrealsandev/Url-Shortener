package config

import (
	"errors"
	"sync"
	"url-shortener/internal/cache"
	"url-shortener/internal/repo"

	"github.com/redis/go-redis/v9"
)

type Provider struct {
	repo  repo.Repository
	cache cache.Cache
	key   string
	mu    sync.Mutex
	local repo.Settings
}

const DEFAULT_SETTINGS_TTL_REDIS_MINUTE = 5

func NewProvider(r repo.Repository, c cache.Cache, key string) *Provider {
	return &Provider{repo: r, cache: c, key: key}
}

func (p *Provider) Get() (repo.Settings, error) {
	if p.cache != nil {
		var settingsFromRedis repo.Settings
		err := p.cache.GetHash(p.key, &settingsFromRedis)

		if err == nil && !settingsFromRedis.IsZero() {
			p.local = settingsFromRedis
			return p.local, nil
		}

		if err != nil && !errors.Is(err, redis.Nil) {
			return repo.Settings{}, err
		}
	}

	settingsFromDB, err := p.repo.GetAllSettings()
	if err != nil {
		return repo.Settings{}, err
	}

	if p.cache != nil && !settingsFromDB.IsZero() {
		_ = p.cache.SetHash(p.key, settingsFromDB, DEFAULT_SETTINGS_TTL_REDIS_MINUTE)
	}

	p.local = *settingsFromDB

	return p.local, nil
}
