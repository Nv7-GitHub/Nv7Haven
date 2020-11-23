package elemental

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

const minVotes = -1
const maxVotes = 4

func getSugg(id string) (Suggestion, error) {
	data, err := db.Get("suggestions/" + id)
	if err != nil {
		return Suggestion{}, err
	}
	if string(data) == "null" {
		return Suggestion{}, errors.New("null")
	}
	var suggestion Suggestion
	err = json.Unmarshal(data, &suggestion)
	if err != nil {
		return Suggestion{}, err
	}
	return suggestion, nil
}

func getSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	suggestion, err := getSugg(id)
	if err != nil {
		if err.Error() == "null" {
			return c.SendString("null")
		}
		return err
	}
	return c.JSON(suggestion)
}

func getSuggestionCombos(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	elem1, err := url.PathUnescape(c.Params("elem1"))
	if err != nil {
		return err
	}
	elem2, err := url.PathUnescape(c.Params("elem2"))
	if err != nil {
		return err
	}
	comboData, err := db.Get("suggestionMap/" + elem1 + "/" + elem2)
	if err != nil {
		return err
	}
	var data []string
	err = json.Unmarshal(comboData, &data)
	if err != nil {
		return err
	}
	return c.JSON(data)
}

func downVoteSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	existing, err := getSugg(id)
	if err != nil {
		return err
	}
	for _, voted := range existing.Voted {
		if voted == uid {
			return c.SendString("You already voted!")
		}
	}
	existing.Votes--
	if existing.Votes < minVotes {
		db.SetData("suggestions/"+id, nil)
	}
	existing.Voted = append(existing.Voted, uid)
	err = db.SetData("suggestions/"+id, existing)
	if err != nil {
		return err
	}
	return nil
}

func upVoteSuggestion(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Headers", "*")
	id, err := url.PathUnescape(c.Params("id"))
	if err != nil {
		return err
	}
	uid := c.Params("uid")
	existing, err := getSugg(id)
	if err != nil {
		return err
	}
	for _, voted := range existing.Voted {
		if voted == uid {
			return c.SendString("You already voted!")
		}
	}
	existing.Votes++
	existing.Voted = append(existing.Voted, uid)
	err = db.SetData("suggestions/"+id, existing)
	if err != nil {
		return err
	}
	if existing.Votes >= maxVotes {
		return c.SendString("create")
	}
	return nil
}

func createSuggestion(c *fiber.Ctx) error {
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
	existing, err := getSugg(id)
	if err != nil {
		return err
	}
	if !(existing.Votes >= maxVotes) {
		return c.SendString("This element still needs more votes!")
	}

	// Get combos
	comboData, err := db.Get("suggestionMap/" + elem1 + "/" + elem2)
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
		err = db.SetData("suggestions/"+val, nil)
		if err != nil {
			return err
		}
	}

	// Delete combos
	err = db.SetData("suggestionMap/"+elem1+"/"+elem2, nil)
	if err != nil {
		return err
	}

	// New Recent Combo
	var recents []RecentCombination
	data, err := db.Get("recent")
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
	recents = append(recents, combo)
	if len(recents) > recentsLength {
		recents = recents[len(recents)-recentsLength-1:]
	}
	fdb, err := fireapp.Database(ctx)
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
		CreatedOn: int(time.Now().Unix()),
	}

	_, err = store.Collection("elements").Doc(newElement.Name).Set(ctx, newElement)
	if err != nil {
		return err
	}

	// Create combo
	existingCombos, err := store.Collection("combos").Doc(elem1).Get(ctx)
	if !(existingCombos == nil || !existingCombos.Exists()) && err != nil {
		return err
	}
	if existingCombos == nil || !existingCombos.Exists() {
		_, err = store.Collection("combos").Doc(elem1).Set(ctx, map[string]string{elem2: newElement.Name})
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
		_, err = store.Collection("combos").Doc(elem1).Set(ctx, elemCombos)
		if err != nil {
			return err
		}
	}
	return nil
}
