package http

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type routerApi struct {
	app   *fiber.App
	cache *[][]map[string]interface{}
}

func NewRouterApi(app *fiber.App) *routerApi {
	return &routerApi{
		app:   app,
		cache: new([][]map[string]interface{}),
	}
}

func (api *routerApi) ListRouter(ctx *fiber.Ctx) error {
	if api.cache == nil || len(*api.cache) <= 0 {
		var (
			data, _ = json.MarshalIndent(api.app.Stack(), "", "  ")
		)
		_ = json.Unmarshal(data, api.cache)
	}
	return ctx.JSON(api.cache)
}
