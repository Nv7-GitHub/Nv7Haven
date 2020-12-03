package nv7haven

import "github.com/gofiber/fiber/v2"

func getIP(c *fiber.Ctx) error {
	forwarded := c.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return c.SendString(forwarded)
	}
	return c.SendString(c.IP())
}
