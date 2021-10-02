package anarchy

import (
	"context"
	"encoding/json"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (a *Anarchy) GetInv(_ context.Context, uid *wrapperspb.StringValue) (*pb.AnarchyInventory, error) {
	var found []string
	var data string
	err := a.db.QueryRow("SELECT inv FROM anarchy_inv WHERE uid=?", uid.Value).Scan(&data)
	if err != nil {
		found = []string{"Air", "Earth", "Fire", "Water"}
	} else {
		err = json.Unmarshal([]byte(data), &found)
		if err != nil {
			return &pb.AnarchyInventory{}, err
		}
	}
	return &pb.AnarchyInventory{
		Found: found,
	}, nil
}

func (a *Anarchy) AddFound(ctx context.Context, req *pb.AnarchyUserRequest) (*emptypb.Empty, error) {
	var cnt int
	var found *pb.AnarchyInventory
	err := a.db.QueryRow("SELECT COUNT(1) FROM anarchy_inv WHERE uid=?", req.Uid).Scan(&cnt)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	if cnt == 0 {
		// Create if not exists
		found = &pb.AnarchyInventory{
			Found: []string{"Air", "Earth", "Fire", "Water"},
		}
		// Updated foundby for those
		for _, elem := range found.Found {
			lock.RLock()
			el := a.cache[elem]
			lock.RUnlock()
			el.FoundBy++
			lock.Lock()
			a.cache[elem] = el
			lock.Unlock()

			_, err = a.db.Exec("UPDATE anarchy_elements SET foundby=? WHERE name=?", el.FoundBy, elem)
			if err != nil {
				return &emptypb.Empty{}, err
			}
		}
		// Add found to DB
		dat, err := json.Marshal(found.Found)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		_, err = a.db.Exec("INSERT INTO anarchy_inv VALUES ( ?, ? )", req.Uid, string(dat))
		if err != nil {
			return &emptypb.Empty{}, err
		}
	} else {
		found, err = a.GetInv(ctx, &wrapperspb.StringValue{Value: req.Uid})
		if err != nil {
			return &emptypb.Empty{}, err
		}
	}

	// Add to inv
	inv := found.Found
	for _, val := range inv {
		if val == req.Element {
			return &emptypb.Empty{}, nil
		}
	}
	inv = append(inv, req.Element)

	dat, err := json.Marshal(inv)
	if err != nil {
		return &emptypb.Empty{}, err
	}
	_, err = a.db.Exec("UPDATE anarchy_inv SET inv=? WHERE uid=?", string(dat), req.Uid)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	// Increment foundby
	lock.RLock()
	el := a.cache[req.Element]
	lock.RUnlock()
	el.FoundBy++
	lock.Lock()
	a.cache[req.Element] = el
	lock.Unlock()

	_, err = a.db.Exec("UPDATE anarchy_elements SET foundby=? WHERE name=?", el.FoundBy, req.Element)
	return &emptypb.Empty{}, err
}
