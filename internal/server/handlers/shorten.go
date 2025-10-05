package handlers

import (
	"errors"
	"net/http"
	"url-shortener/internal/short"

	"github.com/gofiber/fiber/v2"
)

type ShortenHandler struct{ Svc *short.Service }

type shortenReq struct {
	URL         string  `json:"url"`
	CustomAlias *string `json:"custom_alias,omitempty"`
}

func (h ShortenHandler) Serve(c *fiber.Ctx) error {
	var req shortenReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).SendString("bad_request")
	}

	code, shortURL, err := h.Svc.Shorten(c.Context(), req.URL, req.CustomAlias)
	if err != nil {
		switch {
		case errors.Is(err, short.ErrInvalidURL):
			return c.Status(http.StatusBadRequest).SendString("invalid_url")
		case errors.Is(err, short.ErrConflict):
			return c.Status(http.StatusConflict).SendString("conflict")
		default:
			return c.Status(http.StatusInternalServerError).SendString("internal")
		}
	}
	return c.JSON(fiber.Map{"code": code, "short_url": shortURL})
}
