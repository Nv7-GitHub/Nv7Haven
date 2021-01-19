package single

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type packData struct {
	ID          string
	Title       string
	Description string
	Data        string
	UID         string
}

func (s *Single) upload(c *fiber.Ctx) error {
	var dat packData
	err := json.Unmarshal(c.Body(), &dat)
	if err != nil {
		return err
	}
	res := s.db.QueryRow("SELECT COUNT(1) FROM single WHERE uid=? AND id=? LIMIT 1", dat.UID, dat.ID)
	var num int
	err = res.Scan(&num)
	if err != nil {
		return err
	}
	if num == 0 {
		_, err = s.db.Exec("INSERT INTO single VALUES ( ?, ?, ?, ?, ?, ?, ? )", dat.ID, dat.Title, dat.Description, dat.UID, time.Now().Unix(), 0, "{}")
		if err != nil {
			return err
		}
	} else {
		_, err = s.db.Exec("UPDATE single SET createdOn=?, title=?, description=? WHERE id=? AND uid=?", time.Now().Unix(), dat.Title, dat.Description, dat.ID, dat.UID)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(fmt.Sprintf("/home/container/packs/%s_%s.pack", dat.UID, dat.ID), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	err = file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(dat.Data))
	if err != nil {
		return err
	}
	return nil
}
