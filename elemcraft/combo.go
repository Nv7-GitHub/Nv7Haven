package elemcraft

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

func StripRecipe(recipe [][]int) [][]int {
	newR := recipe
	// Forwards
	for _, row := range recipe {
		allZeros := true
		for _, i := range row {
			if i != -1 {
				allZeros = false
				break
			}
		}
		if allZeros {
			newR = newR[1:]
		} else {
			break
		}
	}
	recipe = newR

	// Backwards
	newR = recipe
	for i := len(recipe) - 1; i >= 0; i-- {
		allZeros := true
		for _, i := range recipe[i] {
			if i != -1 {
				allZeros = false
				break
			}
		}
		if allZeros {
			newR = newR[:len(newR)-1]
		} else {
			break
		}
	}

	// Strip cols
	stripCnt := 0
	for i := 0; i < len(recipe[0]); i++ {
		allZeros := true
		for j := range recipe {
			if recipe[j][i] != -1 {
				allZeros = false
				break
			}
		}
		if allZeros {
			stripCnt++
		} else {
			break
		}
	}
	for i := range recipe {
		recipe[i] = recipe[i][stripCnt:]
	}

	// Strip cols backwards
	if len(recipe[0]) == 0 { // Empty
		return recipe
	}
	stripCnt = 0
	for i := len(recipe[0]) - 1; i >= 0; i-- {
		allZeros := true
		for j := range recipe {
			if recipe[j][i] != -1 {
				allZeros = false
				break
			}
		}
		if allZeros {
			stripCnt++
		} else {
			break
		}
	}
	for i := range recipe {
		recipe[i] = recipe[i][:len(recipe[i])-stripCnt]
	}

	return newR
}

func RecipeToString(recipe [][]int) string {
	out := &strings.Builder{}
	for _, row := range recipe {
		for i, v := range row {
			out.WriteString(strconv.Itoa(v))
			if i != len(row)-1 {
				out.WriteString(",")
			}
		}
		out.WriteString("|")
	}
	return out.String()
}

func (e *ElemCraft) Combo(c echo.Context) error {
	// Get recipe
	var comb [][]int
	dec := json.NewDecoder(c.Request().Body)
	err := dec.Decode(&comb)
	if err != nil {
		return err
	}
	comb = StripRecipe(comb)

	// Check if in inv
	u := c.Get(apis.ContextUserKey).(*models.User)
	var inv []int
	err = json.Unmarshal([]byte(u.Profile.GetDataValue("inv").(types.JsonRaw)), &inv)
	if err != nil {
		return err
	}
	invMap := make(map[int]struct{})
	for _, i := range inv {
		invMap[i] = struct{}{}
	}
	for _, row := range comb {
		for _, i := range row {
			if i == -1 {
				continue
			}
			if _, ok := invMap[i]; !ok {
				return fmt.Errorf("elemcraft: missing element %d", i)
			}
		}
	}

	// Check
	r, err := e.app.Dao().FindCollectionByNameOrId("recipes")
	if err != nil {
		return err
	}
	res, err := e.app.Dao().FindFirstRecordByData(r, "recipe", RecipeToString(comb))
	if err != nil {
		return c.JSON(200, nil) // Not found
	}

	// Get result
	e.app.Dao().ExpandRecord(res, []string{"result"}, func(c *models.Collection, ids []string) ([]*models.Record, error) {
		return e.app.Dao().FindRecordsByIds(c, ids, nil)
	})
	el := res.GetExpand()["result"].(*models.Record)
	id := el.GetIntDataValue("index") - 1

	// Check if in inv
	if _, ok := invMap[id]; !ok { // Not in inv, add
		inv = append(inv, id)
		invRaw, err := json.Marshal(inv)
		if err != nil {
			return err
		}
		u.Profile.SetDataValue("inv", types.JsonRaw(invRaw))
		err = e.app.Dao().SaveRecord(u.Profile)
		if err != nil {
			return err
		}
	}

	return c.JSON(200, id)
}
