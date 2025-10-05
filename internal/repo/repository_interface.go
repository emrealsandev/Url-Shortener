package repo

import (
	"context"
)

type Repository interface {
	Insert(url URL) error
	GetByCode(code string) (*URL, error)
	FindOneAndUpdate(ctx context.Context) (uint64, error)
	GetCodeByUrl(urlKey string) (string, error)
	GetAllSettings() (*Settings, error)
}
