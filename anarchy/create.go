package anarchy

import (
	"context"
	"errors"
	"time"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"github.com/finnbear/moderation"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (a *Anarchy) CreateElement(ctx context.Context, req *pb.AnarchyElementCreate) (*emptypb.Empty, error) {
	// Check if innapropriate
	if moderation.IsInappropriate(req.Elem3) {
		return &emptypb.Empty{}, errors.New("no innapropriate suggestions are allowed")
	}

	// Check for exist
	lock.RLock()
	_, exists := a.cache[req.Elem1]
	lock.RUnlock()
	if !exists {
		return &emptypb.Empty{}, errors.New("anarchy: elem1 doesn't exist")
	}
	lock.RLock()
	_, exists = a.cache[req.Elem1]
	lock.RUnlock()
	if !exists {
		return &emptypb.Empty{}, errors.New("anarchy: elem2 doesn't exist")
	}

	// Check for combination
	res, err := a.GetCombination(ctx, &pb.AnarchyCombination{Elem1: req.Elem1, Elem2: req.Elem2})
	if err == nil && res.Exists {
		return &emptypb.Empty{}, errors.New("anarchy: combination already exists")
	}

	// Get creator name
	var creator string
	err = a.db.QueryRow("SELECT name FROM users WHERE uid=?", req.Uid).Scan(&creator)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	// Calc stats
	lock.RLock()
	par1 := a.cache[req.Elem1]
	lock.RUnlock()
	lock.RLock()
	par2 := a.cache[req.Elem2]
	lock.RUnlock()
	var complexity int
	if par1.Complexity > par2.Complexity {
		complexity = int(par1.Complexity)
	} else {
		complexity = int(par2.Complexity)
	}
	complexity++

	// Create element if not exists
	lock.RLock()
	el, exists := a.cache[req.Elem3]
	lock.RUnlock()

	var createdon int64
	if !exists {
		createdon = time.Now().Unix()
		_, err = a.db.Exec("INSERT INTO anarchy_elements VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", req.Elem3, req.Color, req.Comment, req.Elem1, req.Elem2, creator, createdon, complexity, 0, 0)
		if err != nil {
			return &emptypb.Empty{}, err
		}
	} else {
		createdon = el.CreatedOn
	}

	// Update cache
	el = &pb.AnarchyElement{
		Name:       req.Elem3,
		Color:      req.Color,
		Comment:    req.Comment,
		Parents:    []string{req.Elem1, req.Elem2},
		Creator:    creator,
		CreatedOn:  createdon,
		Complexity: int64(complexity),
		Uses:       0,
		FoundBy:    0,
	}
	lock.Lock()
	a.cache[req.Elem3] = el
	lock.Unlock()

	// Create Combo
	_, err = a.db.Exec("INSERT INTO anarchy_combos VALUES ( ?, ?, ? )", req.Elem1, req.Elem2, req.Elem3)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	// Update usedin
	for _, par := range el.Parents {
		lock.RLock()
		parEl := a.cache[par]
		lock.RUnlock()

		parEl.Uses++

		lock.Lock()
		a.cache[par] = parEl
		lock.Unlock()

		_, err = a.db.Exec("UPDATE anarchy_elements SET uses=? WHERE name=?", parEl.Uses, par)
		if err != nil {
			return &emptypb.Empty{}, err
		}
	}

	// Create recent
	_, err = a.db.Exec("INSERT INTO anarchy_recents VALUES (?, ?, ?, ?)", req.Elem1, req.Elem2, req.Elem3, time.Now().Unix())
	if err != nil {
		return &emptypb.Empty{}, err
	}

	a.recents.L.Lock()
	a.recents.Broadcast()
	a.recents.L.Unlock()

	return &emptypb.Empty{}, nil
}
