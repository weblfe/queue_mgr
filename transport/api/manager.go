package api

import (
	"github.com/gofiber/fiber/v2"
)

// QueueManagerApi app队列管理服务接口集合
type QueueManagerApi interface {

	// CreateConsumer godoc
	// @Summary 创建队列消费器
	// @Tags QueueMgrServ
	// @Description create queue consumer
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param consumer formData string true "consumer/消费器名"
	// @Param type formData string true "type/消费器类型" Enums("FastCGI","Native","Shell","Api","Grpc","Proxy","Plugins")
	// @Param properties formData string false "properties/消费器相关参数(json)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /consumer/create [post]
	CreateConsumer(ctx *fiber.Ctx) error

	// CreateQueue godoc
	// @Summary 创建可消费队列信息
	// @Tags QueueMgrServ
	// @Description add queue  info
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param queue formData string true "queue/队列名"
	// @Param driver formData string true "queue/队列链接驱动器类型" Enums("AMQP","MQTT","HTTP","WS","PLUGINS")
	// @Param properties formData string false "properties/可消费队列相关参数(json)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /queue/create [post]
	CreateQueue(ctx *fiber.Ctx) error

	// Bind godoc
	// @Summary 给队列绑定消费协程
	// @Tags QueueMgrServ
	// @Description bind queue consumer
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param consumer formData string false "consumer/消费器名"
	// @Param queue formData string false "queue/队列名"
	// @Param properties formData string false "properties/绑定消费器相关参数(json)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /bind [post]
	Bind(ctx *fiber.Ctx) error

	// State godoc
	// @Summary 查询队列消费器状态
	// @Tags QueueMgrServ
	// @Description query queue consumer state
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param queue query string true "queue/消费队列名"
	// @Param tag query string false "tag/消费进程标签"
	// @Param state query int false "state/消费进程状态" Enums(1,2,3)
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /state [get]
	State(ctx *fiber.Ctx) error

	// Control godoc
	// @Summary 控制队列消费器状态
	// @Tags QueueMgrServ
	// @Description change queue consumer state or processor number
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param queue formData string true "queue/消费队列名"
	// @Param state formData int false "state/消费进程状态" Enums(1,2,3)
	// @Param tag formData string false "tag/消费进程标签"
	// @Param scale formData int false "scale/消费队列协程数扩缩容" default(0)
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /state/update [post]
	Control(ctx *fiber.Ctx) error

	// ListConsumers godoc
	// @Summary 罗列消费器列信息
	// @Tags QueueMgrServ
	// @Description query queue consumer lists
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param page formData int false "page/页码" default(1)
	// @Param count query int false "count/分页量" default(10)
	// @Param state query int false "state/消费进程状态" Enums(0,1,2,3)
	// @Param sort query string false "sort/排序参数" default("created_at:desc")
	// @Param queue query string false "name/限定队列名(模糊匹配eg: test*)"
	// @Param name  query string false "name/限定消费器名(模糊匹配eg: test*)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /consumers [get]
	ListConsumers(ctx *fiber.Ctx) error

	// ListQueues godoc
	// @Summary 罗列消费队列信息
	// @Tags QueueMgrServ
	// @Description query queue consumer lists
	// @Accept  x-www-form-urlencoded
	// @Produce  json
	// @Param Authorization header string true "access jwt token"
	// @Param page formData int false "page/页码" default(1)
	// @Param count query int false "count/分页量" default(10)
	// @Param state query int false "state/消费进程状态" Enums(0,1,2,3)
	// @Param sort query string false "sort/排序参数" default("created_at:desc")
	// @Param queue query string false "name/限定队列名(模糊匹配eg: test*)"
	// @Success 200 {object} entity.JsonResponse
	// @Failure 400,404 {object} entity.JsonResponse
	// @Failure 500 {object} entity.JsonResponse
	// @Failure default {object} entity.JsonResponse
	// @Router /queues [get]
	ListQueues(ctx *fiber.Ctx) error

}
