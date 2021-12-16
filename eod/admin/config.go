package admin

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ConfigRequest struct {
	Guild string `json:"guild"`
	UID   string `json:"uid"`
}

func (a *Admin) Config(c *fiber.Ctx) error {
	var req ConfigRequest
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	v, res := a.data.GetDB(req.Guild)
	if !res.Exists {
		return errors.New(res.Message)
	}
	return c.JSON(v.Config)
}
