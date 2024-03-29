package elemcraft

import (
	"encoding/json"
	"strings"

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

	// Check if element exists
	els, err := e.app.Dao().FindCollectionByNameOrId("elements")
	if err != nil {
		return err
	}
	resEl := dbx.NullStringMap{}
	err = e.app.Dao().RecordQuery(els).AndWhere(dbx.HashExp{"UPPER(name)": strings.ToUpper(req.Name)}).Limit(1).One(&resEl)
	if err == nil {
		el := models.NewRecordFromNullStringMap(els, resEl)
		req.Name = el.GetStringDataValue("name")
		req.Color = el.GetIntDataValue("color")
		req.Description = el.GetStringDataValue("description")
	}

	// Check if suggestion exists
	var id struct{ ID string }
	err = e.app.DB().Select("id").From("suggestions").Where(&dbx.HashExp{"recipe": recStr, "name": req.Name}).One(&id)
	if err == nil { // Exists
		return e.Vote(c, id.ID, user)
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
	return e.Vote(c, sugg.Id, user)
}
