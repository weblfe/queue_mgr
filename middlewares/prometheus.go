package middlewares

import (
		"github.com/gofiber/fiber/v2"
		"github.com/weblfe/queue_mgr/repo"
		"github.com/weblfe/queue_mgr/utils"
)

func CreatePromWare() fiber.Handler {
		var (
				prometheus  = repo.GetPrometheusRepo()
				handlerServ = prometheus.GetHttpHandler()
		)
		return func(ctx *fiber.Ctx) error {
				var responseWriter, request, err = utils.CreateHttpNative(ctx)
				if err != nil {
						return err
				}
				handlerServ.ServeHTTP(responseWriter, request)
				return nil
		}
}
