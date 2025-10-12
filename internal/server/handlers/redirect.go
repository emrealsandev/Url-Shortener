package handlers

import (
	"net/http"
	"url-shortener/internal/repo"

	"url-shortener/internal/short"

	"github.com/gofiber/fiber/v2"
)

type RedirectHandler struct{ Svc *short.Service }

func (h RedirectHandler) Serve(c *fiber.Ctx) error {

	settings, ok := c.Locals("settings").(repo.Settings)
	if !ok {
		return c.Status(http.StatusInternalServerError).SendString("internal: failed to retrieve settings from context")
	}

	code := c.Params("code")
	target, err := h.Svc.Resolve(c.Context(), code, settings)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	return c.Redirect(target, http.StatusFound)
}
