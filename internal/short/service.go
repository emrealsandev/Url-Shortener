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
	repo    Repository
	cache   cache.Cache
	baseURL string
	logger  logger.Logger
}

func NewService(r Repository, c cache.Cache, baseURL string, logger logger.Logger) *Service {
	return &Service{repo: r, cache: c, baseURL: baseURL, logger: logger}
}

func (s *Service) Shorten(ctx context.Context, inputURL string, customAlias *string, ttlHours *int) (string, string, error) {

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

	code, _ := s.checkEligibleCodeExists(target)

	if code != "" {
		s.logger.Info("code exist")
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
	if ttlHours != nil && *ttlHours > 0 {
		e := time.Now().Add(time.Duration(*ttlHours) * time.Hour).UTC()
		exp = &e
	}

	u := URL{Code: code, Target: target, CreatedAt: time.Now().UTC(), ExpiresAt: exp, Disabled: false}
	if err := s.repo.Insert(u); err != nil {
		// repo duplicate → ErrConflict
		return "", "", ErrConflict
	}

	s.processCacheAfterShorten(ctx, code, target)

	return code, s.baseURL + "/" + code, nil
}

func (s *Service) Resolve(ctx context.Context, code string) (string, error) {

	value, hasError, errorMsg := s.cache.GetURLByCode(ctx, code)
	if hasError {
		s.logger.Error(errorMsg.Error())
		return "", ErrSystem
	}

	// rediste varsa onu dön
	if value != "" {
		return value, nil
	}

	// buraya redis gelecek
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

func (s *Service) processCacheAfterShorten(ctx context.Context, code string, target string) {
	_ = s.cache.SetURLByCode(ctx, code, target, time.Duration(24)*time.Hour)
	_ = s.cache.SetCodeByURLKey(ctx, target, code, time.Duration(24)*time.Hour)
}

func (s *Service) checkEligibleCodeExists(target string) (string, error) {
	return s.repo.GetCodeByUrl(target)
}
