package nv7haven

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) syncDb(c *fiber.Ctx) error {
	w := zip.NewWriter(c)
	defer w.Close()

	return filepath.Walk("data", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		f, err := w.Create(path)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		return nil
	})
}
