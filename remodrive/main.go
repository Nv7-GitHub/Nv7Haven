package remodrive

import (
	"context"
	"fmt"
	"sync"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var lock = &sync.RWMutex{}

type RemoDrive struct {
	*pb.UnimplementedRemoDriveServer

	Rooms map[string]Room
}

type Room struct {
	Msgs chan *pb.DriverMessage
}

func InitRemoDrive(grpc *grpc.Server) {
	rd := &RemoDrive{}
	rd.Rooms = make(map[string]Room)
	pb.RegisterRemoDriveServer(grpc, rd)
}

func (r *RemoDrive) CloseRoom(ctx context.Context, roomName *wrapperspb.StringValue) (*emptypb.Empty, error) {
	lock.RLock()
	room, exists := r.Rooms[roomName.Value]
	lock.RUnlock()
	if !exists {
		return &emptypb.Empty{}, fmt.Errorf("remodrive: room %s doesn't exist", roomName.Value)
	}

	close(room.Msgs)

	lock.Lock()
	delete(r.Rooms, roomName.Value)
	lock.Unlock()

	return &emptypb.Empty{}, nil
}
