package server

import (
	"github.com/emrealsandev/Url-Shortener/internal/config"
	"github.com/emrealsandev/Url-Shortener/internal/server/docs"
	handlers2 "github.com/emrealsandev/Url-Shortener/internal/server/handlers"
	"github.com/emrealsandev/Url-Shortener/internal/server/middleware"
	"github.com/emrealsandev/Url-Shortener/internal/short"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, svc *short.Service, settingsProvider *config.Provider) {

	// Serve static files (frontend)
	app.Static("/static", "./web/static")

	// Serve index.html at root
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./web/index.html")
	})

	// api
	api := app.Group("/v1")

	api.Use(
		middleware.Settings(settingsProvider),
		middleware.APILimiter(),
	)

	// health
	api.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })
	api.Get("/readyz", func(c *fiber.Ctx) error { return c.SendString("ready") })

	// swagger docs
	api.Get("/docs", docs.SwaggerUI)
	api.Get("/docs/swagger.json", docs.SwaggerJSON)

	api.Post("/shorten", handlers2.ShortenHandler{Svc: svc}.Serve)

	// v1 altında olmadığı için api grubuna dahil değil.
	app.Get("/:code",
		middleware.Settings(settingsProvider),
		middleware.RedirectLimiter(),
		handlers2.RedirectHandler{Svc: svc}.Serve)
}
