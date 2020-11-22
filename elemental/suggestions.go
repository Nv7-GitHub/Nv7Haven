package elemental

import (
	"encoding/json"
	"net/url"

	"github.com/gofiber/fiber/v2"
)

func getSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	data, err := db.Get("suggestions/" + id)
	if err != nil {
		return err
	}
	if string(data) == "null" {
		return c.SendString("null")
	}
	var suggestion Suggestion
	err = json.Unmarshal(data, &suggestion)
	if err != nil {
		return err
	}
	return c.JSON(suggestion)
}
