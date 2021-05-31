package gdo

import (
	"bufio"
	"errors"

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

	if !exists {
		return errors.New("gdo: client not connected")
	}

	stream.SendEvent(body.Event)

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
		w.WriteString(stream.GetEvent())
		w.Flush()
	}))

	return nil
}
