package elemental

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) createSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	ctx := context.Background()

	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}
	mark, err := url.PathUnescape(c.Params("mark"))
	if err != nil {
		return err
	}
	pioneer, err := url.PathUnescape(c.Params("pioneer"))
	if err != nil {
		return err
	}
	existing, err := e.getSugg(id)
	if err != nil {
		return err
	}
	if !(existing.Votes >= maxVotes) {
		return c.SendString("This element still needs more votes!")
	}

	// Get combos
	comboData, err := e.db.Get("suggestionMap/" + elem1 + "/" + elem2)
	if err != nil {
		return err
	}
	var combos []string
	err = json.Unmarshal(comboData, &combos)
	if err != nil {
		return err
	}

	// Delete hanging elements
	for _, val := range combos {
		err = e.db.SetData("suggestions/"+val, nil)
		if err != nil {
			return err
		}
	}

	// Delete combos
	err = e.db.SetData("suggestionMap/"+elem1+"/"+elem2, nil)
	if err != nil {
		return err
	}

	// New Recent Combo
	var recents []RecentCombination
	data, err := e.db.Get("recent")
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &recents)
	if err != nil {
		return err
	}
	combo := RecentCombination{
		Recipe: [2]string{elem1, elem2},
		Result: id,
	}
	recents = append([]RecentCombination{combo}, recents...)
	if len(recents) > recentsLength {
		recents = recents[:recentsLength-1]
	}
	fdb, err := e.fireapp.Database(ctx)
	if err != nil {
		return err
	}
	ref := fdb.NewRef("recent")
	err = ref.Set(ctx, recents)
	if err != nil {
		return err
	}

	// Create Element
	newElement := Element{
		Name:      existing.Name,
		Color:     fmt.Sprintf("%s_%f_%f", existing.Color.Base, existing.Color.Saturation, existing.Color.Lightness),
		Creator:   existing.Creator,
		Pioneer:   pioneer,
		Parents:   []string{elem1, elem2},
		Comment:   mark,
		CreatedOn: int(time.Now().Unix()) * 1000,
	}

	elementExists, _ := e.store.Collection("elements").Doc(newElement.Name).Get(ctx)
	if !elementExists.Exists() {
		_, err = e.store.Collection("elements").Doc(newElement.Name).Set(ctx, newElement)
		if err != nil {
			return err
		}
	}

	// Create combo
	existingCombos, err := e.store.Collection("combos").Doc(elem1).Get(ctx)
	if !(existingCombos == nil || !existingCombos.Exists()) && err != nil {
		return err
	}
	if existingCombos == nil || !existingCombos.Exists() {
		_, err = e.store.Collection("combos").Doc(elem1).Set(ctx, map[string]string{elem2: newElement.Name})
		if err != nil {
			return err
		}
	} else {
		var elemCombos map[string]string
		err = existingCombos.DataTo(&elemCombos)
		if err != nil {
			return err
		}
		elemCombos[elem2] = newElement.Name
		_, err = e.store.Collection("combos").Doc(elem1).Set(ctx, elemCombos)
		if err != nil {
			return err
		}
	}
	return nil
}
