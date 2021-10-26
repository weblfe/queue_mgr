package starter

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/repo"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	quit chan os.Signal // 监听信道
)

func init() {
	// 命令行启动
	GetArgumentsStarter().Init()
	// 加载环境配置
	GetEnvStarter().Init()
}

// StartUp 应用初始化
func StartUp() {
	// 环境变量 服务
	GetEnvStarter().StartUp()
	// 配置 服务
	GetSettingStarter().StartUp()
	// 日志 服务
	GetLoggerStarter().StartUp()
	// 数控库 服务
	GetDataBaseStarter().StartUp()
	// 依赖服务注册
	GetServiceRegisterStarter().StartUp()
	// 后台任务组件 注册
	GetScheduleStarter().StartUp()
	// 应用 主程服务
	GetAppStarter().StartUp()
}

// Run 应用 运行
func Run() {
	// 应用启动
	_ = repo.GetPoolRepo().Add(GetAppStarter().Run)
	shutdown()
}

// 关闭
func shutdown() {
	var repos = repo.GetPoolRepo()
	if quit == nil {
		quit = make(chan os.Signal)
	}
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infoln("stop signal")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// 关闭
	defer repos.GetPool().Release()
	defer cancel()
	// 释放携程池
	// 等待5秒
	select {
	case <-ctx.Done():
		log.Infoln("stopped")
		return
	}
}
