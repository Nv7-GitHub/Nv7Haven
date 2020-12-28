package elemental

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (e *Elemental) createSuggestion(c *fiber.Ctx) error {

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
	if !(existing.Votes >= maxVotes) && (int(time.Now().Weekday()) != 5) {
		return c.SendString("This element still needs more votes!")
	}

	// Get combos
	comboData, err := e.getSuggestions(elem1)
	if err != nil {
		return err
	}
	combos := comboData[elem2]

	// Delete hanging elements
	for _, val := range combos {
		_, err = e.db.Exec("DELETE FROM suggestions WHERE name=?", val)
		if err != nil {
			return err
		}
	}

	// Delete combos
	delete(comboData, elem2)
	data, err := json.Marshal(comboData)
	if err != nil {
		return err
	}
	_, err = e.db.Exec("UPDATE suggestion_combos SET combos=? WHERE name=?", data, elem1)
	if err != nil {
		return err
	}

	// New Recent Combo
	var recents []RecentCombination
	data, err = e.fdb.Get("recent")
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
	e.fdb.SetData("recent", recents)

	res, err := e.db.Query("SELECT COUNT(1) FROM elements WHERE name=?", existing.Name)
	defer res.Close()
	if err != nil {
		return err
	}
	var count int
	res.Next()
	res.Scan(&count)
	if count == 0 {
		color := fmt.Sprintf("%s_%f_%f", existing.Color.Base, existing.Color.Saturation, existing.Color.Lightness)
		_, err = e.db.Exec("INSERT INTO elements VALUES( ?, ?, ?, ?, ?, ?, ?, ? )", existing.Name, color, mark, elem1, elem2, existing.Creator, pioneer, int(time.Now().Unix())*1000)
		if err != nil {
			return err
		}
	}

	// Create combo
	err = e.addCombo(elem1, elem2, existing.Name)
	if err != nil {
		return err
	}
	return nil
}
