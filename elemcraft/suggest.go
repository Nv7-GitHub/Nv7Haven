package elemcraft

import (
	"encoding/json"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
)

const MinVotes = 3

type SuggestionRequest struct {
	Recipe      [][]int `json:"recipe"`
	Name        string  `json:"name"`
	Color       int     `json:"color"`
	Description string  `json:"description"`
}

func (e *ElemCraft) Suggest(c echo.Context) error {
	// Get req & user
	var req SuggestionRequest
	dec := json.NewDecoder(c.Request().Body)
	err := dec.Decode(&req)
	if err != nil {
		return err
	}
	user := c.Get(apis.ContextUserKey).(*models.User)
	recStr := RecipeToString(StripRecipe(req.Recipe))

	// TODO: Check if element exists

	// Check if suggestion exists
	var id struct{ ID string }
	err = e.app.DB().Select("id").From("suggestions").Where(&dbx.HashExp{"recipe": recStr, "name": req.Name}).One(&id)
	if err == nil { // Exists
		return e.Vote(id.ID, user)
	}

	// Create suggestion
	suggs, err := e.app.Dao().FindCollectionByNameOrId("suggestions")
	if err != nil {
		return err
	}
	sugg := models.NewRecord(suggs)
	sugg.Load(map[string]any{
		"recipe":      recStr,
		"name":        req.Name,
		"color":       req.Color,
		"description": req.Description,
		"creator":     user.Id,
	})
	err = e.app.Dao().SaveRecord(sugg)
	if err != nil {
		return err
	}

	// Vote
	return e.Vote(sugg.Id, user)
}
