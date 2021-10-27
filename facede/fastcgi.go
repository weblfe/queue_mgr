package facede

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"sync"
	"time"
)

type Handler interface {
	// Cancel 取消处理
	Cancel()
	// Type 类型
	Type() string
	// Call 处理请求
	Call(ctx *fiber.Ctx) error
	// Proxy 代理
	Proxy(res http.ResponseWriter, req *http.Request) error
	// SetTimeout 设置超时时长
	SetTimeout(duration time.Duration)
	// Register 注册处理器池
	Register(pool *sync.Pool)
}
