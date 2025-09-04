package short

import (
	"context"
	"errors"
	"time"
	"url-shortener/internal/security"
	"url-shortener/pkg/base62"
)

var (
	ErrInvalidURL = errors.New("invalid_url")
	ErrConflict   = errors.New("conflict")
	ErrExpired    = errors.New("expired")
	ErrNotFound   = errors.New("not_found")
	ErrSequence   = errors.New("sequence_error")
)

type Cache interface {
	// ileride eklenecek (Get/Set)
}

type Service struct {
	repo    Repository
	cache   Cache
	baseURL string
}

func NewService(r Repository, c Cache, baseURL string) *Service {
	return &Service{repo: r, cache: c, baseURL: baseURL}
}

func (s *Service) Shorten(ctx context.Context, inputURL string, customAlias *string, ttlHours *int) (string, string, error) {
	target, err := security.NormalizeUrl(inputURL)
	if err != nil {
		return "", "", ErrInvalidURL
	}

	var code string
	if customAlias != nil && *customAlias != "" {
		code = *customAlias
	} else {
		seq, err := s.GetSeqNum(ctx)
		if err != nil {
			return "", "", ErrSequence
		}
		code = base62.Encode(seq)
	}

	var exp *time.Time
	if ttlHours != nil && *ttlHours > 0 {
		e := time.Now().Add(time.Duration(*ttlHours) * time.Hour).UTC()
		exp = &e
	}

	u := URL{Code: code, Target: target, CreatedAt: time.Now().UTC(), ExpiresAt: exp, Disabled: false}
	if err := s.repo.Insert(u); err != nil {
		// repo duplicate → ErrConflict
		return "", "", ErrConflict
	}
	return code, s.baseURL + "/" + code, nil
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {
	u, err := s.repo.GetByCode(code)
	if err != nil || u == nil {
		return "", ErrNotFound
	}

	if u.Disabled {
		return "", ErrNotFound
	}

	if u.ExpiresAt != nil && u.ExpiresAt.Before(time.Now().UTC()) {
		return "", ErrExpired
	}
	return u.Target, nil
}

func (s *Service) GetSeqNum(ctx context.Context) (uint64, error) {
	seq, err := s.repo.FindOneAndUpdate(ctx)
	if err != nil {
		return 0, err
	}
	return seq, nil
}
