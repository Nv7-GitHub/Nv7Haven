package remodrive

import "fmt"

const maxBufLen = 10

func (r *RemoDrive) CloseRoomByName(name string) error {
	lock.RLock()
	room, exists := r.Rooms[name]
	lock.RUnlock()
	if !exists {
		return fmt.Errorf("remodrive: room %s doesn't exist", name)
	}

	close(room.Msgs)

	lock.Lock()
	delete(r.Rooms, name)
	lock.Unlock()

	return nil
}

func (r *RemoDrive) NewRoom(name string) error {
	r.CloseRoomByName(name)

	msgs := make(chan string, maxBufLen)

	lock.Lock()
	r.Rooms[name] = Room{
		Msgs: msgs,
	}
	lock.Unlock()

	return nil
}
