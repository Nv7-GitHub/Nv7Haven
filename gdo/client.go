package gdo

import (
	"bufio"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func (g *GDO) clientConnect(c *fiber.Ctx) error {
	uid := c.Params("uid")

	g.Lock.RLock()
	stream, exists := g.Streams[uid]
	g.Lock.RUnlock()

	if !exists {
		stream = NewEventStream(uid)

		g.Lock.Lock()
		g.Streams[uid] = stream
		g.Lock.Unlock()
	}

	// SSE (give client stream next event)
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

func (g *GDO) clientResponse(c *fiber.Ctx) error {
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
	return nil
}

func (g *GDO) clientDisconnect(c *fiber.Ctx) error {
	uid := c.Body()

	g.Lock.RLock()
	delete(g.Streams, string(uid))
	g.Lock.RUnlock()

	return nil
}
