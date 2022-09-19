package elemcraft

import (
	"strconv"

	"github.com/labstack/echo/v5"
)

type Element struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Color       int    `json:"color"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Created     int    `json:"created"`
}

func (e *ElemCraft) GetElement(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}
	els, err := e.app.Dao().FindCollectionByNameOrId("elements")
	if err != nil {
		return err
	}
	el, err := e.app.Dao().FindFirstRecordByData(els, "index", id+1)
	if err != nil {
		return err
	}
	creator, err := e.app.Dao().FindUserById(el.GetStringDataValue("creator"))
	if err != nil {
		return err
	}

	return c.JSON(200, Element{
		ID:          int(el.GetFloatDataValue("index")) - 1,
		Name:        el.GetStringDataValue("name"),
		Color:       int(el.GetIntDataValue("color")),
		Description: el.GetStringDataValue("description"),
		Creator:     creator.Profile.GetStringDataValue("name"),
		Created:     int(el.GetIntDataValue("created")),
	})
}
