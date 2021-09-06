package anarchy

import (
	"context"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

const recentsLength = 30

func (a *Anarchy) GetRecents(ctx context.Context, _ *emptypb.Empty) (*pb.AnarchyRecents, error) {
	recents := make([]*pb.AnarchyRecentCombination, 0, recentsLength)
	res, err := a.db.Query("SELECT elem1, elem2, elem3 FROM anarchy_recents ORDER BY createdon DESC LIMIT ?", recentsLength)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		item := &pb.AnarchyRecentCombination{}
		err = res.Scan(&item.Elem1, &item.Elem2, &item.Elem3)
		if err != nil {
			return nil, err
		}

		recents = append(recents, item)
	}

	return &pb.AnarchyRecents{
		Recents: recents,
	}, nil
}

func (a *Anarchy) WaitForNextRecent(_ *emptypb.Empty, stream pb.Anarchy_WaitForNextRecentServer) error {
	a.recents.L.Lock()
	a.recents.Wait()
	a.recents.L.Unlock()

	err := stream.Send(&emptypb.Empty{})
	if err != nil {
		return err
	}

	return nil
}
