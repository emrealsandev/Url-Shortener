package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func APILimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Dakikada en fazla 20 istek.
		Max: 20,

		// Limit penceresinin süresi.
		Expiration: 1 * time.Minute,

		// Hangi isteğin kime ait olduğunu belirleyen anahtar.
		// c.IP() kullanmak, her bir IP adresini ayrı ayrı limitlememizi sağlar.
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		// LimitExceeded alanına gerek yok. Middleware, limit aşıldığında
		// otomatik olarak 429 status kodu ve "Too Many Requests" mesajını döner.
	})
}

func RedirectLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		// Saniyede en fazla 5 istek (dakikada 300 istek).
		// Aynı IP'den çok hızlı ve tekrar tekrar tıklanmasını engeller.
		Max: 5,

		// Burst (ani) trafiği engellemek için pencereyi kısa tutuyoruz.
		Expiration: 1 * time.Second,

		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	})
}
