package elemental

import (
	"encoding/json"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) foundElement(c *fiber.Ctx) error {
	return nil
}

func (e *Elemental) getFound(c *fiber.Ctx) error {
	
	
	res, err := e.db.Query("SELECT found FROM users WHERE uid=?", c.Params("uid"))
	if err != nil {
		return err
	}
	defer res.Close()
	var data string
	res.Next()
	err = res.Scan(&data)
	if err != nil {
		return err
	}
	return c.SendString(data)
}

func (e *Elemental) newFound(c *fiber.Ctx) error {
	
	
	var found []string
	res, err := e.db.Query("SELECT found FROM users WHERE uid=?", c.Params("uid"))
	if err != nil {
		return err
	}
	defer res.Close()
	var data string
	res.Next()
	res.Scan(&data)
	err = json.Unmarshal([]byte(data), &found)
	if err != nil {
		return err
	}

	elem, err := url.PathUnescape(c.Params("elem"))
	if err != nil {
		return err
	}
	for _, val := range found {
		if val == elem {
			return nil
		}
	}
	found = append(found, elem)

	dat, err := json.Marshal(found)
	if err != nil {
		return err
	}
	data = string(dat)
	_, err = e.db.Exec("UPDATE users SET found=? WHERE uid=?", data, c.Params("uid"))
	if err != nil {
		return err
	}
	return nil
}
