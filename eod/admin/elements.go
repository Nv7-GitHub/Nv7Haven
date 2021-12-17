package admin

import (
	"encoding/json"
	"errors"

	"github.com/Nv7-Github/Nv7Haven/eod/types"
	"github.com/gofiber/fiber/v2"
)

type ElementRequest struct {
	Guild string `json:"guild"`
	UID   string `json:"uid"`
	ID    int    `json:"id"`
}

func (a *Admin) Element(c *fiber.Ctx) error {
	var req ElementRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	if !a.CheckUID(req.UID) {
		return errors.New("admin: invalid uid")
	}

	v, res := a.data.GetDB(req.Guild)
	if !res.Exists {
		return errors.New(res.Message)
	}
	el, res := v.GetElement(req.ID)
	if !res.Exists {
		return errors.New(res.Message)
	}
	return c.JSON(el)
}

type ElementSetRequest struct {
	Guild    string `json:"guild"`
	UID      string `json:"uid"`
	ElemData string `json:"elemData"`
}

func (a *Admin) SetElement(c *fiber.Ctx) error {
	var req ElementSetRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}
	if !a.CheckUID(req.UID) {
		return errors.New("admin: invalid uid")
	}

	v, res := a.data.GetDB(req.Guild)
	if !res.Exists {
		return errors.New(res.Message)
	}

	var el types.Element
	err = json.Unmarshal([]byte(req.ElemData), &el)
	if err != nil {
		return err
	}

	err = v.SaveElement(el)
	return err
}
