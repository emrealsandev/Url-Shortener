package middleware

import (
	"url-shortener/internal/config"
	"url-shortener/internal/repo"

	"github.com/gofiber/fiber/v2"
)

func Settings(provider *config.Provider) fiber.Handler {
	return func(c *fiber.Ctx) error {
		settings, err := provider.Get()
		if err != nil {
			// Ayarları çekemiyorsak, bu kritik bir hatadır.
			// Burada kendi logger'ınızı kullanarak loglama yapabilirsiniz.
			return c.Status(fiber.StatusInternalServerError).SendString("could not load application settings")
		}

		// Ayarları isteğin context'ine ("locals") koyuyoruz.
		c.Locals(repo.COLLECTION_SETTINGS, settings)

		// Bir sonraki middleware veya handler'a geç.
		return c.Next()
	}
}
