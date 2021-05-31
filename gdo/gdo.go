package gdo

import "github.com/gofiber/fiber/v2"

type GDO struct {
}

func (g *GDO) routing(app *fiber.App) {

}

func InitGDO(app *fiber.App) {
	gdo := GDO{}
	gdo.routing(app)
}
