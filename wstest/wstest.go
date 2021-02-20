package main

import (
	"fmt"
	"os"

	"github.com/rgamba/evtwebsocket"
)

var cnchan = make(chan bool)

func main() {
	msg := evtwebsocket.Msg{
		Body: []byte("Test message"),
		Callback: func(resp []byte, _ *evtwebsocket.Conn) {
			// This function executes when the server responds
			fmt.Printf("Got response: %s\n", resp)
			cnchan <- true
		},
	}
	c := evtwebsocket.Conn{
		OnConnected: func(w *evtwebsocket.Conn) {
			fmt.Println("Connected")
			w.Send(msg)
		},
		OnMessage: func(msg []byte, _ *evtwebsocket.Conn) {
			fmt.Printf("Received message: %s\n", msg)
			cnchan <- true
		},
		OnError: func(err error) {
			fmt.Printf("** ERROR **\n%s\n", err.Error())
		},
	}
	// Connect
	c.Dial("ws://localhost:"+os.Getenv("PORT")+"/ws/123?v=1.0", "")
	<-cnchan
}
