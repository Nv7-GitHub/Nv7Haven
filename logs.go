package main

import (
	"io"
	"os"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

func systemHandlers(app *fiber.App) {
	app.Get("/freememory", func(c *fiber.Ctx) error {
		debug.FreeOSMemory()
		return nil
	})
	app.Get("/kill/:password", func(c *fiber.Ctx) error {
		if c.Params("password") == os.Getenv("PASSWORD") {
			os.Exit(2)
		}
		return nil
	})
	app.Get("/logs", func(c *fiber.Ctx) error {
		file, err := os.Open("logs.txt")
		if err != nil {
			return err
		}
		_, err = io.Copy(c, file)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	})
	app.Get("/createlogs", func(c *fiber.Ctx) error {
		file, err := os.Open("createlogs.txt")
		if err != nil {
			return err
		}
		_, err = io.Copy(c, file)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	})
	app.Get("/discordlogs", func(c *fiber.Ctx) error {
		file, err := os.Open("discordlogs.txt")
		if err != nil {
			return err
		}
		_, err = io.Copy(c, file)
		if err != nil {
			return err
		}
		file.Close()
		return nil
	})
}
