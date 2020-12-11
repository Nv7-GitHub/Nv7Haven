package elemental

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

const recentsLength int = 30

// RecentCombination represents the data for a recent combination
type RecentCombination struct {
	Recipe [2]string
	Result string
}

func (e *Elemental) getRecents(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	var recents []RecentCombination
	data, err := e.db.Get("recent")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &recents)
	if err != nil {
		return err
	}
	return c.JSON(recents)
}
