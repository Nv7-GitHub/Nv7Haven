package admin

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"
)

type ConfigRequest struct {
	Guild     string `json:"guild"`
	UID       string `json:"uid"`
	NewConfig string `json:"newConfig"`
}

func (a *Admin) Config(c *fiber.Ctx) error {
	var req ConfigRequest
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
	return c.JSON(v.Config)
}

func (a *Admin) SetConfig(c *fiber.Ctx) error {
	var req ConfigRequest
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

	v.Config.PlayChannels = nil
	err = json.Unmarshal([]byte(req.NewConfig), &v.Config)
	if err != nil {
		return err
	}

	err = v.SaveConfig()
	return err
}
