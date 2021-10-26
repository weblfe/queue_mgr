package api

import "github.com/gofiber/fiber/v2"

type RouterApi interface {
	// ListRouter godoc
	// @Summary 罗列服务接口列表
	// @Tags AppCdnServ
	// @Description List Service Api Routers
	// @Produce  json
	// @Success 200 {object} entity.KvMap
	// @Failure 400,404 {object} entity.KvMap
	// @Failure 500 {object} entity.KvMap
	// @Failure default {object} entity.KvMap
	// @Router /routers [get]
	ListRouter(ctx *fiber.Ctx) error
}
