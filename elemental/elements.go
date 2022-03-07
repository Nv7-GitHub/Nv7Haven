package elemental

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type empty struct{}

var lock = &sync.RWMutex{}

// GetElement gets an element from the database
func (e *Elemental) GetElement(elemName string) (*pb.Element, error) {
	lock.RLock()
	val, exists := e.cache[elemName]
	lock.RUnlock()
	if !exists {
		return e.RefreshElement(elemName)
	}
	return val, nil
}

// RefreshElement gets an element from the database and refreshes the local cache with that element
func (e *Elemental) RefreshElement(elemName string) (*pb.Element, error) {
	elem := &pb.Element{}
	res, err := e.db.Query("SELECT * FROM elements WHERE name=?", elemName)
	if err != nil {
		return &pb.Element{}, err
	}
	defer res.Close()
	elem.Parents = make([]string, 2)
	res.Next()
	err = res.Scan(&elem.Name, &elem.Color, &elem.Comment, &elem.Parents[0], &elem.Parents[1], &elem.Creator, &elem.Pioneer, &elem.CreatedOn, &elem.Complexity, &elem.Uses, &elem.FoundBy)
	if err != nil {
		return &pb.Element{}, err
	}
	if (elem.Parents[0] == "") && (elem.Parents[1] == "") {
		elem.Parents = make([]string, 0)
	}

	lock.Lock()
	e.cache[elemName] = elem
	lock.Unlock()
	return elem, nil
}

func (e *Elemental) GetElem(_ context.Context, inp *wrapperspb.StringValue) (*pb.Element, error) {
	return e.GetElement(inp.Value)
}

func (e *Elemental) GetCombination(_ context.Context, inp *pb.Combination) (*pb.CombinationResult, error) {
	comb, suc, err := e.GetCombo(inp.Elem1, inp.Elem2)

	return &pb.CombinationResult{
		Data:   comb,
		Exists: suc,
	}, err
}

// GetCombo gets a combination
func (e *Elemental) GetCombo(elem1, elem2 string) (string, bool, error) {
	res, err := e.db.Query("SELECT COUNT(1) FROM elem_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?) LIMIT 1", elem1, elem2, elem2, elem1)
	if err != nil {
		return "", false, err
	}
	defer res.Close()
	var count int
	res.Next()
	res.Scan(&count)
	if count == 0 {
		return "", false, nil
	}

	res, err = e.db.Query("SELECT elem3 FROM elem_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?) LIMIT 1", elem1, elem2, elem2, elem1)
	if err != nil {
		return "", false, err
	}
	defer res.Close()
	var elem3 string
	res.Next()
	err = res.Scan(&elem3)
	if err != nil {
		return "", false, err
	}

	return elem3, true, nil
}

const (
	minBatchCount = 60
	batchCount    = 10
)

func (e *Elemental) GetAll(uid *wrapperspb.StringValue, stream pb.Elemental_GetAllServer) error {
	res, err := e.db.Query("SELECT found FROM users WHERE uid=?", uid.Value)
	if err != nil {
		return err
	}
	defer res.Close()
	var data string
	res.Next()
	err = res.Scan(&data)
	if err != nil {
		return err
	}
	var found []string
	err = json.Unmarshal([]byte(data), &found)
	if err != nil {
		return err
	}

	recents, err := e.GetRecents()
	if err != nil {
		return err
	}

	req := make(map[string]empty)
	for _, val := range found {
		req[val] = empty{}
	}
	for _, rec := range recents {
		req[rec.Elem1] = empty{}
		req[rec.Elem2] = empty{}
		req[rec.Elem3] = empty{}
	}

	total := int64(len(req))

	getAllBatchSize := minBatchCount
	if total/batchCount > int64(getAllBatchSize) {
		getAllBatchSize = int(total / batchCount)
	}

	batch := make([]*pb.Element, 0, getAllBatchSize)

	i := 0
	for k := range req {
		elem, err := e.GetElement(k)
		if err != nil {
			return err
		}
		batch = append(batch, elem)
		if i == getAllBatchSize {
			err = stream.Send(&pb.GetAllChunk{
				Count:    total,
				Elements: batch,
			})
			if err != nil {
				return err
			}
			batch = make([]*pb.Element, 0, getAllBatchSize)
			i = 0
		}
		i++
	}
	if len(batch) > 0 {
		err = stream.Send(&pb.GetAllChunk{
			Count:    total,
			Elements: batch,
		})
		return err
	}
	return nil
}
