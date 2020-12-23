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
	
	
	var recents []RecentCombination
	data, err := e.fdb.Get("recent")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &recents)
	if err != nil {
		return err
	}
	return c.JSON(recents)
}
