package handlers

import (
	"net/http"

	"url-shortener/internal/short"

	"github.com/gofiber/fiber/v2"
)

type RedirectHandler struct{ Svc *short.Service }

func (h RedirectHandler) Serve(c *fiber.Ctx) error {
	code := c.Params("code")
	target, err := h.Svc.Resolve(c.Context(), code)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	return c.Redirect(target, http.StatusFound)
}
