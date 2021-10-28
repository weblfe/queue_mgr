package http

import "github.com/gofiber/fiber/v2"

type ManagerApi struct {
	Controller
}

func NewManagerApi() *ManagerApi {
	var api = new(ManagerApi)
	return api
}

// CreateConsumer 创建队列消费器
func (mgr *ManagerApi) CreateConsumer(ctx *fiber.Ctx) error {
	panic("implement me")
}

// CreateQueue 创建可消费队列信息
func (mgr *ManagerApi) CreateQueue(ctx *fiber.Ctx) error {
	panic("implement me")
}

// Bind 给队列绑定消费协程
func (mgr *ManagerApi) Bind(ctx *fiber.Ctx) error {
	panic("implement me")
}

// State 查询队列消费器状态
func (mgr *ManagerApi) State(ctx *fiber.Ctx) error {
	panic("implement me")
}

// Control 控制消费队列状态
func (mgr *ManagerApi) Control(ctx *fiber.Ctx) error {
	panic("implement me")
}

// ListConsumers 罗列消费器列列表
func (mgr *ManagerApi) ListConsumers(ctx *fiber.Ctx) error {
	panic("implement me")
}

// ListQueues 罗列消费队列列表
func (mgr *ManagerApi) ListQueues(ctx *fiber.Ctx) error {
	panic("implement me")
}
