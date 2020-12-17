package nv7haven

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}
func dashbaordHandler(w http.ResponseWriter, r *http.Request) {
	chn := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second)
		chn <- true
	}()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	timeout := time.After(1 * time.Second)
	select {
	case ev := <-chn:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(ev)
		fmt.Fprintf(w, "data: %v\n\n", buf.String())
		fmt.Printf("data: %v\n", buf.String())
	case <-timeout:
		fmt.Fprintf(w, ": nothing to sent\n\n")
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
