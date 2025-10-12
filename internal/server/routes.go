package server

import (
	"url-shortener/internal/config"
	handlers2 "url-shortener/internal/server/handlers"
	"url-shortener/internal/server/middleware"
	"url-shortener/internal/short"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, svc *short.Service, settingsProvider *config.Provider) {
	// health
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })
	app.Get("/readyz", func(c *fiber.Ctx) error { return c.SendString("ready") })

	// api
	api := app.Group("/v1")

	api.Use(
		middleware.Settings(settingsProvider),
		middleware.APILimiter(),
	)

	api.Post("/shorten", handlers2.ShortenHandler{Svc: svc}.Serve)

	// redirect
	app.Get("/:code",
		middleware.Settings(settingsProvider),
		middleware.RedirectLimiter(),
		handlers2.RedirectHandler{Svc: svc}.Serve)
}
