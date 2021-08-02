package elemental

import (
	"context"
	"encoding/json"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// GetFound gets a user's found elements, based on their UID
func (e *Elemental) GetFound(uid string) ([]string, error) {
	var found []string
	res, err := e.db.Query("SELECT found FROM users WHERE uid=?", uid)
	if err != nil {
		return []string{}, err
	}
	defer res.Close()
	var data string
	res.Next()
	res.Scan(&data)
	err = json.Unmarshal([]byte(data), &found)
	if err != nil {
		return []string{}, err
	}
	return found, nil
}

func (e *Elemental) GetInv(ctx context.Context, uid *wrapperspb.StringValue) (*pb.Inventory, error) {
	res, err := e.GetFound(uid.Value)
	return &pb.Inventory{
		Found: res,
	}, err
}

func (e *Elemental) AddFound(ctx context.Context, req *pb.SuggestionRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, e.NewFound(req.Element, req.Uid)
}

// NewFound adds an element to a user's savefile
func (e *Elemental) NewFound(elem string, uid string) error {
	var found []string
	res, err := e.db.Query("SELECT found FROM users WHERE uid=?", uid)
	if err != nil {
		return err
	}
	defer res.Close()
	var data string
	res.Next()
	res.Scan(&data)
	err = json.Unmarshal([]byte(data), &found)
	if err != nil {
		return err
	}

	for _, val := range found {
		if val == elem {
			return nil
		}
	}
	found = append(found, elem)

	dat, err := json.Marshal(found)
	if err != nil {
		return err
	}
	data = string(dat)
	_, err = e.db.Exec("UPDATE users SET found=? WHERE uid=?", data, uid)
	if err != nil {
		return err
	}

	// increment foundby
	el, err := e.GetElement(elem)
	if err != nil {
		return err
	}
	el.FoundBy++
	lock.Lock()
	e.cache[el.Name] = el
	lock.Unlock()
	_, err = e.db.Exec("UPDATE elements SET foundby=? WHERE name=?", el.FoundBy, el.Name)
	if err != nil {
		return err
	}
	return nil
}
