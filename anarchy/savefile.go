package anarchy

import (
	"context"
	"encoding/json"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (a *Anarchy) GetInv(ctx context.Context, uid *wrapperspb.StringValue) (*pb.AnarchyInventory, error) {
	var found []string
	var data string
	err := a.db.QueryRow("SELECT inv FROM users WHERE uid=?", uid.Value).Scan(&data)
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
	}, err
}

func (a *Anarchy) AddFound(ctx context.Context, req *pb.AnarchyUserRequest) (*emptypb.Empty, error) {
	found, err := a.GetInv(ctx, &wrapperspb.StringValue{Value: req.Uid})
	if err != nil {
		// Create if not exists
		found = &pb.AnarchyInventory{
			Found: []string{"Air", "Earth", "Fire", "Water"},
		}
		dat, err := json.Marshal(found.Found)
		if err != nil {
			return &emptypb.Empty{}, err
		}
		_, err = a.db.Exec("INSERT INTO anarchy_inv VALUES ( ?, ? )", req.Uid, string(dat))
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

	_, err = a.db.Exec("UPDATE anarchy_elements SET foundby=? WHERE element=?", el.FoundBy, req.Element)
	return &emptypb.Empty{}, err
}
