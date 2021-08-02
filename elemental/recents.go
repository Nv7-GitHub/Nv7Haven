package elemental

import (
	"context"
	"time"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const recentsLength = 30

func (e *Elemental) GetRecents() ([]*pb.RecentCombination, error) {
	recents := make([]*pb.RecentCombination, 0, recentsLength)
	res, err := e.db.Query("SELECT elem1, elem2, elem3 FROM recents ORDER BY createdon DESC LIMIT ?", recentsLength)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		item := &pb.RecentCombination{}
		err = res.Scan(&item.Elem1, &item.Elem2, &item.Elem3)
		if err != nil {
			return nil, err
		}

		recents = append(recents, item)
	}

	return recents, nil
}

func (e *Elemental) NewRecent(recent *pb.RecentCombination) error {
	_, err := e.db.Exec("INSERT INTO recents VALUES (?, ?, ?, ?)", recent.Elem1, recent.Elem2, recent.Elem3, time.Now().Unix())
	if err != nil {
		return err
	}

	e.recents.L.Lock()
	e.recents.Broadcast()
	e.recents.L.Unlock()
	return nil
}

func (e *Elemental) GetRec(ctx context.Context, _ *emptypb.Empty) (*pb.Recents, error) {
	recents, err := e.GetRecents()
	if err != nil {
		return &pb.Recents{}, err
	}
	return &pb.Recents{Recents: recents}, nil
}

func (e *Elemental) WaitForNextRecent(_ *emptypb.Empty, stream pb.Elemental_WaitForNextRecentServer) error {
	e.recents.L.Lock()
	e.recents.Wait()
	e.recents.L.Unlock()

	err := stream.Send(&emptypb.Empty{})
	if err != nil {
		return err
	}

	return nil
}
