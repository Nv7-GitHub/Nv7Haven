package elemental

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (e *Elemental) CreateSugg(ctx context.Context, req *pb.CreateRequest) (*emptypb.Empty, error) {
	suc, msg := e.CreateSuggestion(req.Mark, req.Pioneer, req.Elem1, req.Elem2, req.Id)
	if !suc {
		return &emptypb.Empty{}, errors.New(msg)
	}
	return &emptypb.Empty{}, nil
}

func (e *Elemental) incrementUses(id string) error {
	elem, err := e.GetElement(id)
	if err != nil {
		return err
	}
	elem.Uses++
	lock.Lock()
	e.cache[elem.Name] = elem
	lock.Unlock()
	_, err = e.db.Exec("UPDATE elements SET uses=? WHERE name=?", elem.Uses, elem.Name)
	if err != nil {
		return err
	}
	return nil
}

// CreateSuggestion creates a suggestion
func (e *Elemental) CreateSuggestion(mark string, pioneer string, elem1 string, elem2 string, id string) (bool, string) {
	existing, err := e.getSugg(id)
	if err != nil {
		return false, err.Error()
	}
	if !(existing.Votes >= maxVotes) && (time.Now().Weekday() != anarchyDay) {
		return false, "This element still needs more votes!"
	}

	// Get combos
	combos, err := e.GetSuggestions(elem1, elem2)
	if err != nil {
		return false, err.Error()
	}

	// Delete hanging elements
	for _, val := range combos {
		_, err = e.db.Exec("DELETE FROM suggestions WHERE name=?", val)
		if err != nil {
			return false, err.Error()
		}
	}

	// Delete combos
	_, err = e.db.Exec("DELETE FROM sugg_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?)", elem1, elem2, elem2, elem1)
	if err != nil {
		return false, err.Error()
	}

	res, err := e.db.Query("SELECT COUNT(1) FROM elements WHERE name=?", existing.Name)
	if err != nil {
		return false, err.Error()
	}
	defer res.Close()

	parent1, err := e.GetElement(elem1)
	if err != nil {
		return false, err.Error()
	}
	parent2, err := e.GetElement(elem2)
	if err != nil {
		return false, err.Error()
	}
	complexity := max(parent1.Complexity, parent2.Complexity) + 1

	err = e.incrementUses(elem1)
	if err != nil {
		return false, err.Error()
	}
	if elem2 != elem1 {
		err = e.incrementUses(elem2)
		if err != nil {
			return false, err.Error()
		}
	}

	var count int
	res.Next()
	res.Scan(&count)
	if count == 0 {
		color := fmt.Sprintf("%s_%f_%f", existing.Color.Base, existing.Color.Saturation, existing.Color.Lightness)
		_, err = e.db.Exec("INSERT INTO elements VALUES( ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", existing.Name, color, mark, elem1, elem2, existing.Creator, pioneer, int(time.Now().Unix())*1000, complexity, 0, 0)
		if err != nil {
			return false, err.Error()
		}
	}

	// Create combo
	err = e.addCombo(elem1, elem2, existing.Name)
	if err != nil {
		return false, err.Error()
	}

	// New Recent Combo
	err = e.NewRecent(&pb.RecentCombination{
		Elem1: elem1,
		Elem2: elem2,
		Elem3: id,
	})
	if err != nil {
		return false, err.Error()
	}

	return true, ""
}
