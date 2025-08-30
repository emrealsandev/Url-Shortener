package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"url-shortener/internal/short"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Options struct {
	Port    string
	BaseURL string
	Repo    short.Repository
}

type Server struct {
	app *fiber.App
	opt Options
}

func New(opt Options) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Prefork:      false,
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(compress.New())

	// Domain service (şimdilik cache yok)
	svc := short.NewService(opt.Repo, nil, opt.BaseURL)

	// Routes
	registerRoutes(app, svc)

	return &Server{app: app, opt: opt}
}

func (s *Server) Start(ctx context.Context) error {
	addr := ":" + s.opt.Port
	log.Println("listening on", addr)

	// Fiber listen’i ayrı goroutine’de; context iptaliyle kapanalım
	errCh := make(chan error, 1)
	go func() { errCh <- s.app.Listen(addr) }()

	select {
	case <-ctx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		_ = s.app.ShutdownWithContext(shutCtx)
		return nil
	case err := <-errCh:
		// Fiber shutdown’da da error dönebilir; loglayıp geri ver
		return fmt.Errorf("fiber: %w", err)
	}
}
