package http

import "github.com/gofiber/fiber/v2"

type Controller struct{}

func (c *Controller) getTransport(ctx *fiber.Ctx) *httpTransport {
	var tran = createTransport(ctx)
	return &tran
}
