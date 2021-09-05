package anarchy

import (
	"context"
	"encoding/json"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	minBatchCount = 60
	batchCount    = 10
)

// GetElement gets an element from the database
func (a *Anarchy) GetElem(ctx context.Context, name *wrapperspb.StringValue) (*pb.AnarchyElement, error) {
	lock.RLock()
	val := a.cache[name.Value]
	lock.RUnlock()
	return val, nil
}

func (a *Anarchy) GetCombination(ctx context.Context, combo *pb.AnarchyCombination) (*pb.AnarchyCombinationResult, error) {
	var cnt int
	err := a.db.QueryRow("SELECT COUNT(1) FROM anarchy_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?)", combo.Elem1, combo.Elem2, combo.Elem2, combo.Elem1).Scan(&cnt)
	if err != nil {
		return &pb.AnarchyCombinationResult{}, err
	}
	if cnt == 0 {
		return &pb.AnarchyCombinationResult{Exists: false}, nil
	}

	var elem3 string
	err = a.db.QueryRow("SELECT elem3 FROM anarchy_combos WHERE (elem1=? AND elem2=?) OR (elem1=? AND elem2=?)", combo.Elem1, combo.Elem2, combo.Elem2, combo.Elem1).Scan(&elem3)
	if err != nil {
		return &pb.AnarchyCombinationResult{}, err
	}
	return &pb.AnarchyCombinationResult{Exists: true, Data: elem3}, nil
}

type empty struct{}

func (a *Anarchy) GetAll(uid *wrapperspb.StringValue, stream pb.Anarchy_GetAllServer) error {
	res, err := a.db.Query("SELECT inv FROM users WHERE uid=?", uid.Value)
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

	recents, err := a.GetRecents(context.Background(), &emptypb.Empty{})
	if err != nil {
		return err
	}

	req := make(map[string]empty)
	for _, val := range found {
		req[val] = empty{}
	}
	for _, rec := range recents.Recents {
		req[rec.Elem1] = empty{}
		req[rec.Elem2] = empty{}
		req[rec.Elem3] = empty{}
	}

	total := int64(len(req))

	getAllBatchSize := minBatchCount
	if total/batchCount > int64(getAllBatchSize) {
		getAllBatchSize = int(total / batchCount)
	}

	batch := make([]*pb.AnarchyElement, 0, getAllBatchSize)

	i := 0
	for k := range req {
		lock.RLock()
		elem := a.cache[k]
		lock.RUnlock()
		batch = append(batch, elem)
		if i == getAllBatchSize {
			err = stream.Send(&pb.AnarchyGetAllChunk{
				Count:    total,
				Elements: batch,
			})
			if err != nil {
				return err
			}
			batch = make([]*pb.AnarchyElement, 0, getAllBatchSize)
			i = 0
		}
		i++
	}
	if len(batch) > 0 {
		err = stream.Send(&pb.AnarchyGetAllChunk{
			Count:    total,
			Elements: batch,
		})
		return err
	}
	return nil
}
