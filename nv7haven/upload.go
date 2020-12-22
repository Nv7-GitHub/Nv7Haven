package nv7haven

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

const numDir = "/home/container/num.txt"
const fileDir = "/home/container/file%d%s"

func (n *Nv7Haven) fileNum() (int, error) {
	info, err := os.Stat(numDir)
	if !(os.IsNotExist(err) && !info.IsDir()) {
		numDat, err := ioutil.ReadFile(numDir)
		num, err := strconv.Atoi(string(numDat))
		if err != nil {
			return 0, err
		}
		return num, nil
	}
	return 0, nil
}

func (n *Nv7Haven) incrementFileNum() error {
	fileNum, err := n.fileNum()
	if err != nil {
		return err
	}
	fileNum++
	file, err := os.OpenFile(numDir, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(strconv.Itoa(fileNum)))
	if err != nil {
		return err
	}
	return nil
}

func (n *Nv7Haven) upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	fileNum, err := n.fileNum()
	if err != nil {
		return err
	}
	extension := strings.Split(file.Filename, ".")
	ext := ""
	if len(extension) > 1 {
		ext = "." + extension[len(extension)-1]
	}
	_, err = n.sql.Exec("INSERT INTO upload VALUES ( ?, ?, ? )", fileNum, ext, time.Now().Unix()+86400)
	if err != nil {
		return err
	}
	go n.checkDates()
	err = n.incrementFileNum()
	if err != nil {
		return err
	}
	return c.SaveFile(file, fmt.Sprintf(fileDir, fileNum, ext))
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
