package facede

import (
	"github.com/weblfe/queue_mgr/entity"
	"time"
)

type Scheduler interface {
	Watch()
	Start()
	Stop()
	Refresh()
	SetDiscover(discover func(ch chan<- *entity.Crontab, t time.Time)) Scheduler
	Add(cron *entity.Crontab)
}
