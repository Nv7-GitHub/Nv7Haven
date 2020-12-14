package nv7haven

import "github.com/gofiber/fiber/v2"

func (d *Nv7Haven) getIP(c *fiber.Ctx) error {
	return c.SendString(c.IPs()[0])
}
