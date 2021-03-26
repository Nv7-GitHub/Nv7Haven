package nv7haven

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const fileDir = "files/file%d%s"

func (n *Nv7Haven) upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	var id int
	nums := make(map[int]bool)
	var thing int
	res, err := n.sql.Query("SELECT id FROM upload WHERE expiry<=?", time.Now().Unix())
	if err != nil {
		return err
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(&thing)
		if err != nil {
			return err
		}
		nums[thing] = true
	}
	exists := true
	for exists {
		id = rand.Intn(10000)
		_, exists = nums[id]
	}

	extension := strings.Split(file.Filename, ".")
	ext := ""
	if len(extension) > 1 {
		ext = "." + extension[len(extension)-1]
	}
	_, err = n.sql.Exec("INSERT INTO upload VALUES ( ?, ?, ? )", id, ext, time.Now().Unix()+86400)
	if err != nil {
		return err
	}
	go n.checkDates()
	if err != nil {
		return err
	}

	if _, err := os.Stat("files"); os.IsNotExist(err) {
		err = os.Mkdir("files", 0777)
		if err != nil {
			return err
		}
	}
	err = c.SaveFile(file, fmt.Sprintf(fileDir, id, ext))
	if err != nil {
		return err
	}
	return c.SendString(strconv.Itoa(id))
}

func (n *Nv7Haven) checkDates() {
	res, err := n.sql.Query("SELECT id, extension FROM upload WHERE expiry<=?", time.Now().Unix())
	if err != nil {
		log.Println(err)
	}
	defer res.Close()
	var num int
	var ext string
	for res.Next() {
		err = res.Scan(&num, &ext)
		if err != nil {
			log.Println(err)
		}
		err = os.Remove(fmt.Sprintf(fileDir, num, ext))
		if err != nil {
			log.Println(err)
		}
	}
	_, err = n.sql.Exec("DELETE FROM upload WHERE expiry<=?", time.Now().Unix())
	if err != nil {
		log.Println(err)
	}
}

func (n *Nv7Haven) getFile(c *fiber.Ctx) error {
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	id = strings.Split(id, ".")[0]
	num, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	res, err := n.sql.Query("SELECT extension FROM upload WHERE id=? LIMIT 1", num)
	if err != nil {
		return err
	}
	defer res.Close()
	var ext string
	res.Next()
	err = res.Scan(&ext)
	if err != nil {
		return err
	}

	id, err = url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}

	if (len(strings.Split(id, ".")) < 2) && (strings.Contains(ext, ".")) {
		return c.Redirect(c.Path() + ext)
	}
	err = c.SendFile(fmt.Sprintf(fileDir, num, ext))
	if err != nil {
		return err
	}
	return err
}
