package remodrive

import (
	"fmt"
	"io"

	"github.com/Nv7-Github/Nv7Haven/pb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (r *RemoDrive) Drive(stream pb.RemoDrive_DriveServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}

		lock.RLock()
		room, exists := r.Rooms[msg.Room]
		lock.RUnlock()
		if !exists {
			return fmt.Errorf("remodrive: room %s doesn't exist", msg.Room)
		}

		room.Msgs <- msg
	}
}

func (r *RemoDrive) Host(room *wrapperspb.StringValue, stream pb.RemoDrive_HostServer) error {
	msgs := make(chan *pb.DriverMessage)

	lock.Lock()
	r.Rooms[room.Value] = Room{
		Msgs: msgs,
	}
	lock.Unlock()

	for msg := range msgs {
		if err := stream.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
