package elemcraft

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

func (e *ElemCraft) Vote(id string, user *models.User) error {
	// Check if exists
	var cnt struct{ Cnt int }
	err := e.app.DB().Select("COUNT(*) as cnt").From("votes").Where(&dbx.HashExp{"suggestion": id, "user": user.Id}).One(&cnt)
	if err != nil {
		return err
	}
	if cnt.Cnt > 0 {
		return errors.New("elemcraft: already voted")
	}

	// Add vote
	suggs, err := e.app.Dao().FindCollectionByNameOrId("suggestions")
	if err != nil {
		return err
	}
	rec := models.NewRecord(suggs)
	rec.Load(map[string]any{
		"suggestion": id,
		"user":       user.Id,
	})
	err = e.app.Dao().SaveRecord(rec)
	if err != nil {
		return err
	}

	// Check if enough votes
	err = e.app.DB().Select("COUNT(*) as cnt").From("votes").Where(&dbx.HashExp{"suggestion": id}).One(&cnt)
	if err != nil {
		return err
	}
	if cnt.Cnt >= MinVotes {
		sugg, err := e.app.Dao().FindRecordById(suggs, id, nil)
		if err != nil {
			return err
		}
		els, err := e.app.Dao().FindCollectionByNameOrId("elements")
		if err != nil {
			return err
		}
		recipes, err := e.app.Dao().FindCollectionByNameOrId("recipes")
		if err != nil {
			return err
		}

		// 1. Create element
		var id int
		err = e.app.DB().Select("MAX(index)").From("elements").One(&id)
		if err != nil {
			return err
		}
		id += 1

		el := models.NewRecord(els)
		el.Load(map[string]any{
			"index":       id,
			"name":        sugg.GetStringDataValue("name"),
			"color":       sugg.GetIntDataValue("color"),
			"description": sugg.GetStringDataValue("description"),
			"creator":     sugg.GetStringDataValue("creator"),
		})
		err = e.app.Dao().SaveRecord(el)
		if err != nil {
			return err
		}

		// 2. Create recipe
		recipe := models.NewRecord(recipes)
		recipe.Load(map[string]any{
			"recipe": sugg.GetStringDataValue("recipe"),
			"result": el.Id,
		})
		err = e.app.Dao().SaveRecord(recipe)
		if err != nil {
			return err
		}

		// 3. Delete suggestion (this deletes votes)
		err = e.app.Dao().DeleteRecord(sugg)
		if err != nil {
			return err
		}
	}

	return nil
}
