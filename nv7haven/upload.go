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

const fileDir = "/home/container/file%d%s"

func (n *Nv7Haven) upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	var id int
	var nums map[int]bool
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
	err = res.Scan(&ext)
	if err != nil {
		return err
	}
	file := fmt.Sprintf(fileDir, num, ext)
	err = c.SendFile(file)
	if err != nil {
		return err
	}
	return os.Remove(file)
}
