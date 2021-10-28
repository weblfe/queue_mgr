package router

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/weblfe/queue_mgr/docs"
	"github.com/weblfe/queue_mgr/middlewares"
	"github.com/weblfe/queue_mgr/transport/http"
	"github.com/weblfe/queue_mgr/utils"
	"strings"
	_ "strings"
)

// Http 注册路由
func Http(app *fiber.App) {

	// 查询接口
	var (
		monitorWare = monitor.New()
		routerApi   = http.NewRouterApi(app)
		managerApi  = http.NewManagerApi()
		promWare    = middlewares.CreatePromWare()
	)

	// 跨域
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin, Content-Type, Accept, Content-Language, Accept-Language, Authorization",
		AllowCredentials: true,
	}))

	// requestID
	app.Use(requestid.New(requestid.Config{Header: "X-Request-ID"}))
	docs.SwaggerInfo.Host = strings.TrimPrefix(utils.GetEnvVal("APP_URL", docs.SwaggerInfo.Host), docs.GetDefaultSchema())
	// 数据监控
	app.Get("/", monitorWare)
	app.Get("/queue_mgr", monitorWare)
	app.Get("/metrics", promWare)

	var router = app.Group("/queue_mgr")
	// expose prometheus metrics 接口
	router.All("/metrics", promWare)
	// 数据监控
	router.Get("/dashboard", monitorWare)
	// swag
	// 是否开启swagger docs
	if utils.GetEnvBool("APP_ENABLE_DOCS") {
		router.Get("/swagger/*", swagger.Handler)
	}
	// 路由信息列表
	router.Get("/routers", routerApi.ListRouter)

	// --- QueueManager-API ---
	// 罗列消费器列列表
	router.Get("/consumers", managerApi.ListConsumers)
	// 罗列消费队列列表
	router.Get("/queues", managerApi.ListQueues)

	// 控制消费队列状态
	router.Post("/state/update", managerApi.Control)
	// 查询队列消费器状态
	router.Get("/state", managerApi.State)
	// 给队列绑定消费协程
	router.Post("/bind", managerApi.Bind)

	// 创建可消费队列信息
	router.Post("/queue/create", managerApi.CreateQueue)
	// 创建队列消费器
	router.Post("/consumer/create", managerApi.CreateConsumer)
}
