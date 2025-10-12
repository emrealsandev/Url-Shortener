package short

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
	"url-shortener/internal/cache"
	"url-shortener/internal/config"
	"url-shortener/internal/logger"
	"url-shortener/internal/repo"
	"url-shortener/internal/security"
	"url-shortener/pkg/base62"
)

var (
	ErrInvalidURL = errors.New("invalid_url")
	ErrConflict   = errors.New("conflict")
	ErrExpired    = errors.New("expired")
	ErrNotFound   = errors.New("not_found")
	ErrSequence   = errors.New("sequence_error")
	ErrSystem     = errors.New("system_error")
)

type Service struct {
	repo    repo.Repository
	cache   cache.Cache
	baseURL string
	logger  logger.Logger
}

func NewService(r repo.Repository, c cache.Cache, baseURL string, logger logger.Logger) *Service {
	return &Service{repo: r, cache: c, baseURL: baseURL, logger: logger}
}

func (s *Service) Shorten(ctx context.Context, inputURL string, customAlias *string, settings repo.Settings) (string, string, error) {

	target, err := security.NormalizeUrl(inputURL)
	if err != nil {
		return "", "", ErrInvalidURL
	}

	value, hasError, errorMsg := s.cache.GetCodeByURLKey(ctx, target)
	if hasError {
		s.logger.Error(errorMsg.Error())
		return "", "", ErrSystem
	}

	// rediste varsa onu dön
	if value != "" {
		s.logger.Info("cache hit")
		return value, s.baseURL + "/" + value, nil
	}

	code, _ := s.repo.GetCodeByUrl(target)

	if code != "" {
		s.logger.Info("code exist")
		s.processCacheAfterShorten(ctx, code, target, settings)
		return code, s.baseURL + "/" + code, nil
	}

	if customAlias != nil && *customAlias != "" {
		code = *customAlias
	} else {
		seq, err := s.GetSeqNum(ctx)
		if err != nil {
			return "", "", ErrSequence
		}

		salt, err := strconv.ParseUint(strings.ReplaceAll(config.Get().SequenceSalt, "_", ""), 0, 64)
		if err != nil {
			return "", "", ErrSequence
		}

		code = base62.Encode(seq ^ salt)
	}

	var exp *time.Time
	if !settings.IsZero() && settings.TtlTime > 0 {
		e := time.Now().Add(time.Duration(settings.TtlTime) * time.Hour).UTC()
		exp = &e
	}

	u := repo.URL{Code: code, Target: target, CreatedAt: time.Now().UTC(), ExpiresAt: exp, Disabled: false}
	if err := s.repo.Insert(u); err != nil {
		// repo duplicate → ErrConflict
		return "", "", ErrConflict
	}

	s.processCacheAfterShorten(ctx, code, target, settings)

	return code, s.baseURL + "/" + code, nil
}

func (s *Service) Resolve(ctx context.Context, code string, settings repo.Settings) (string, error) {

	value, hasError, errorMsg := s.cache.GetURLByCode(ctx, code)
	if hasError {
		s.logger.Error(errorMsg.Error())
		return "", ErrSystem
	}

	if value != "" {
		return value, nil
	}

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

	s.processCacheAfterShorten(ctx, code, u.Target, settings)
	return u.Target, nil
}

func (s *Service) GetSeqNum(ctx context.Context) (uint64, error) {
	seq, err := s.repo.FindOneAndUpdate(ctx)
	if err != nil {
		return 0, err
	}
	return seq, nil
}

func (s *Service) processCacheAfterShorten(ctx context.Context, code string, target string, settings repo.Settings) {
	exp := time.Duration(5) * time.Minute

	if !settings.IsZero() && settings.RedisTtlTime > 0 {
		exp = time.Duration(settings.RedisTtlTime) * time.Minute
	}

	_ = s.cache.SetURLByCode(ctx, code, target, exp)
	_ = s.cache.SetCodeByURLKey(ctx, target, code, exp)
}
