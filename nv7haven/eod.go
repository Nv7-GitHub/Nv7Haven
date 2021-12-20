package nv7haven

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func (n *Nv7Haven) getEoDDB(c *fiber.Ctx) error {
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
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				return err
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"

			err := addFiles(w, newBase, baseInZip+file.Name()+"/")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
