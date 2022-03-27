package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/debug"

	"github.com/Nv7-Github/Nv7Haven/eod/logs"
	"github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

var monitors = [][]string{{"measure_temp"}, {"measure_volts"}, {"get_mem", "arm"} /*, {"get_mem", "gpu"}, {"get_throttled"}*/} // Commented out part 1 gets VRAM, commented out part 2 gets if throttled

func systemHandlers(app *fiber.App) {
	if runtime.GOOS == "linux" {
		app.Get("/temp", func(c *fiber.Ctx) error {
			for _, m := range monitors {
				cmd := exec.Command("vcgencmd", m...)
				cmd.Stdout = c
				err := cmd.Run()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	app.Get("/freememory", func(c *fiber.Ctx) error {
		debug.FreeOSMemory()
		return nil
	})
	app.Get("/kill/:password", func(c *fiber.Ctx) error {
		if c.Params("password") == os.Getenv("MOD_PASSWORD") {
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

	mysql.SetLogger(&Logger{})
	app.Get("/mysqlogs", func(c *fiber.Ctx) error {
		file, err := os.Open("mysqlogs.txt")
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

type Logger struct{}

func (l *Logger) Print(args ...any) {
	log.SetOutput(logs.MysqLogs)
	log.Print(args...)
}
