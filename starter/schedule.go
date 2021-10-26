package starter

import (
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/domain"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/repo"
	"github.com/weblfe/queue_mgr/utils"
)

// 任务执行
type scheduleStarter struct {
	scheduler facede.Scheduler
	baseStarterConstructor
}

var (
	schedulerStarter = newScheduleStarter()
)

func GetScheduleStarter() *scheduleStarter {
	return schedulerStarter
}

func newScheduleStarter() *scheduleStarter {
	var starter = new(scheduleStarter)
	starter.baseStarterConstructor = newStarterConstructor()
	starter.name = "scheduleStarter"
	return starter
}

func (starter *scheduleStarter) StartUp() {
	starter.init(starter.boot)
}

func (starter *scheduleStarter) boot() {
	// 业务功能内部初始化
	// domain.Boot()
	if starter.scheduler == nil {
		starter.scheduler = domain.NewScheduler()
	}

	/*var number = utils.GetEnvInt("SCHEDULE_NUMBER", 3)
	if number <= 0 {
		number = 3
	}*/
	if !utils.GetEnvBool("SCHEDULE_ON") {
		return
	}
	var container = repo.GetContainerRepo()
	starter.scheduler.SetDiscover(nil)
	// 注册后台 任务对象
	container.Register("scheduler", starter.scheduler)
	// 任务监听
	if err := repo.RegisterProcessor(starter.scheduler.Start); err != nil {
		panic("scheduleStarter Processor Register Error: " + err.Error())
	}
	// 任务发现
	if err := repo.RegisterProcessor(starter.scheduler.Watch); err != nil {
		panic("scheduleStarter Processor Register Error: " + err.Error())
	}
	log.Infoln("scheduleStarter started")
}
