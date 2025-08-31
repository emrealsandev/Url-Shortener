package server

import (
	"url-shortener/internal/handlers"
	"url-shortener/internal/short"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, svc *short.Service) {
	// health
	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })
	app.Get("/readyz", func(c *fiber.Ctx) error { return c.SendString("ready") })

	// api
	api := app.Group("/v1")
	api.Post("/shorten", handlers.ShortenHandler{Svc: svc}.Serve)

	// redirect
	app.Get("/:code", handlers.RedirectHandler{Svc: svc}.Serve)
}
