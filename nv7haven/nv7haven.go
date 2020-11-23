package nv7haven

import "github.com/gofiber/fiber/v2"

// InitNv7Haven initializes the handlers for Nv7Haven
func InitNv7Haven(app *fiber.App) error {
	app.Get("/hella/:input", calcHella)
	return nil
}
