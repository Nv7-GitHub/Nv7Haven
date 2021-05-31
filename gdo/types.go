package gdo

import "sync"

type EventStream struct {
	Cond  *sync.Cond
	UID   string
	Event string
}

func (e *EventStream) GetEvent() string {
	e.Cond.L.Lock()
	defer e.Cond.L.Unlock()
	e.Cond.Wait()

	return e.Event
}

func (e *EventStream) SendEvent(event string) {
	e.Cond.L.Lock()
	e.Event = event
	e.Cond.Broadcast()
	e.Cond.L.Unlock()
}

func NewEventStream(uid string) *EventStream {
	return &EventStream{
		UID:   uid,
		Cond:  sync.NewCond(&sync.Mutex{}),
		Event: "",
	}
}
