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
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	var found []string
	data, err := e.db.Get("users/" + c.Params("uid") + "/found")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &found)
	if err != nil {
		return err
	}
	return c.JSON(found)
}

func (e *Elemental) newFound(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	var found []string
	data, err := e.db.Get("users/" + c.Params("uid") + "/found")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &found)
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
	err = e.db.SetData("users/"+c.Params("uid")+"/found", found)
	if err != nil {
		return err
	}
	return nil
}
