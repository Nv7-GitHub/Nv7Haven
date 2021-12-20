package nv7haven

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) downloadEoDDB(c *fiber.Ctx) error {
	if string(c.Body()) != os.Getenv("PASSWORD") {
		return errors.New("eodb: invalid password")
	}

	out := zip.NewWriter(c)
	defer out.Close()

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	return addFiles(out, filepath.Join(home, "go/src/github.com/Nv7-Github/Nv7haven/data/eod"), "")
}

func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			newFolder := filepath.Join(baseInZip, file.Name())
			newBase := filepath.Join(basePath, file.Name())
			err := addFiles(w, newBase, newFolder)
			if err != nil {
				return err
			}
		} else {
			f, err := os.Open(filepath.Join(basePath, file.Name()))
			if err != nil {
				return err
			}
			g, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				return err
			}
			_, err = io.Copy(g, f)
			if err != nil {
				return err
			}
			f.Close()
		}
	}

	return nil
}
