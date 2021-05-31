package gdo

import (
	"bufio"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type eventRequest struct {
	Event string
	UID   string
}

func (g *GDO) sendEvent(c *fiber.Ctx) error {
	var body eventRequest
	err := c.BodyParser(&body)
	if err != nil {
		return err
	}

	g.Lock.RLock()
	stream, exists := g.Streams[body.UID]
	g.Lock.RUnlock()

	// SSE (wait for response)
	ctx := c.Context()
	ctx.SetContentType("text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")
	ctx.Response.Header.Set("Transfer-Encoding", "chunked")
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
	ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if !exists {
			w.WriteString("gdo: client not connected")
		} else {
			stream.SendEvent(body.Event)
			w.WriteString(stream.GetEvent())
		}
		w.Flush()
	}))

	return nil
}
