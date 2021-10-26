package starter

import (
	"github.com/weblfe/queue_mgr/config"
	"github.com/weblfe/queue_mgr/repo"
	"github.com/weblfe/queue_mgr/service"
)

type serviceRegisterStarter struct {
	baseStarterConstructor
}

// GetServiceRegisterStarter 获取服务 注册启动器
func GetServiceRegisterStarter() *serviceRegisterStarter {
	var starter = serviceRegisterStarter{}
	starter.baseStarterConstructor = newStarterConstructor()
	starter.name = "serviceStarter"
	return &starter
}

// StartUp 启动 初始化 全局配置
func (starter *serviceRegisterStarter) StartUp() {
	starter.init(starter.boot)
}

// 内部 服务 注册
func (starter *serviceRegisterStarter) boot() {
	// 服务包内部初始化
	service.Boot()
	var servRepo = repo.GetServerRegisterRepository()
	servRepo.Load(service.GetUserServiceImpl())

	// 注册数据库 服务
	repo.GetDatabaseRepository().InitConnection(config.GetAppConfig().GetDbKv())

}
