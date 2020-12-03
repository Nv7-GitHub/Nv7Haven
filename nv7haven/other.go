package nv7haven

import "github.com/gofiber/fiber/v2"

func getIP(c *fiber.Ctx) error {
	return c.JSON(c.IPs())
}
