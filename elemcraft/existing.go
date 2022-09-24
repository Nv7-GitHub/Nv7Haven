package elemcraft

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
)

type ExistingSuggestion struct {
	Name        string  `json:"name"`
	Color       float64 `json:"color"`
	Description string  `json:"description"`
}

func (e *ElemCraft) Existing(c echo.Context) error {
	// Get recipe
	var recipe [][]int
	dec := json.NewDecoder(c.Request().Body)
	err := dec.Decode(&recipe)
	if err != nil {
		return err
	}
	recipe = StripRecipe(recipe)
	recStr := RecipeToString(recipe)

	// Get existing
	var res []ExistingSuggestion
	err = e.app.DB().Select("name", "color", "description").From("suggestions").Where(&dbx.HashExp{"recipe": recStr}).All(&res)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Return
	return c.JSON(200, res)
}
