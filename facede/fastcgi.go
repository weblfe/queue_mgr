package facede

import (
	"github.com/gofiber/fiber/v2"
	"sync"
	"time"
)

type Handler interface {
	// Cancel 取消处理
	Cancel()
	// Type 类型
	Type() string
	// Handle 处理请求
	Handle(ctx *fiber.Ctx) error
	// SetTimeout 设置超时时长
	SetTimeout(duration time.Duration)
	// Register 注册处理器池
	Register(pool *sync.Pool)
}
