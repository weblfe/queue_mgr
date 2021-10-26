package domain

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/repo"
	"github.com/weblfe/queue_mgr/utils"
	"sort"
	"sync"
	"time"
)

type schedulerDomain struct {
	ctx          context.Context
	cancel       context.CancelFunc
	channel      chan *entity.Crontab
	queue        []*entity.Crontab
	cacheIndex   map[string]bool
	locker       sync.RWMutex
	ticker       *time.Ticker
	notify       chan int
	state        int
	logger       *logrus.Logger
	debug        bool
	discoverFunc func(ch chan<- *entity.Crontab, t time.Time)
}

const (
	StateRunning  = 1
	StateDiscover = 2
	StateQueue    = 3
	StateStop     = 4
	StateRefresh  = 5
)

func NewScheduler(duration ...time.Duration) *schedulerDomain {
	duration = append(duration, utils.GetEnvDuration("APP_SCHEDULE_INTERVAL", 3*time.Minute))
	var (
		scheduler   = new(schedulerDomain)
		ctx, cancel = context.WithCancel(context.Background())
	)
	scheduler.ctx = ctx
	scheduler.cancel = cancel
	scheduler.channel = make(chan *entity.Crontab, 10)
	scheduler.ticker = time.NewTicker(duration[0])
	scheduler.cacheIndex = make(map[string]bool)
	scheduler.locker = sync.RWMutex{}
	scheduler.notify = make(chan int, 2)
	scheduler.debug = utils.GetEnvBool("APP_DEBUG")
	scheduler.logger = repo.GetLogger("scheduler")
	return scheduler
}

func (scheduler *schedulerDomain) SetDiscover(discover func(ch chan<- *entity.Crontab, t time.Time)) facede.Scheduler {
	if scheduler.discoverFunc == nil {
		scheduler.discoverFunc = discover
	}
	return scheduler
}

func (scheduler *schedulerDomain) Start() {
	scheduler.logger.Infoln("scheduler Start ...")
	for {
		select {
		case <-scheduler.ctx.Done():
			close(scheduler.channel)
		case cron := <-scheduler.channel:
			var ok = cron.Parse()
			if scheduler.debug {
				scheduler.logger.WithFields(logrus.Fields{
					"ID":       cron.ID(),
					"Data":     cron.Data,
					"At":       cron.At,
					"Parse":    ok,
					"callback": cron.Check(),
				}).Infoln("scheduler cron ...")
			}
			if ok {
				scheduler.execute(cron)
			} else {
				scheduler.cache(cron)
			}
		}
	}
}

func (scheduler *schedulerDomain) Watch() {
	scheduler.logger.Infoln("scheduler Watching ...")
	for {
		select {
		case <-scheduler.ctx.Done():
			close(scheduler.channel)
		// 定时发现
		case <-scheduler.ticker.C:
			scheduler.discover()
		// 通知触发更新
		case v := <-scheduler.notify:
			if v == StateRefresh && scheduler.state != StateDiscover && scheduler.state != StateStop {
				scheduler.logger.Infoln("scheduler Refresh ...")
				scheduler.discover()
			}
		// 缓存 延迟任务检查
		default:
			scheduler.dispatch()
		}
	}
}

// 定时任务执行发现
func (scheduler *schedulerDomain) dispatch() {
	scheduler.locker.Lock()
	defer scheduler.locker.Unlock()
	var (
		index = -1
		debug = scheduler.debug
		size  = len(scheduler.queue)
	)
	if size == 0 {
		return
	}
	// 检查任务列表
	for i, v := range scheduler.queue {
		var match = v.Parse()
		if match {
			index = i
			scheduler.channel <- v
			scheduler.remove(v.ID())
			if debug {
				scheduler.logger.WithFields(logrus.Fields{
					"index": i,
					"ID":    v.ID(),
				}).Infoln("任务触发")
			}
		}
		if debug && !match {
			scheduler.logger.WithFields(logrus.Fields{
				"index":    i,
				"ID":       v.ID(),
				"Duration": v.Duration().String(),
				"DateTime": v.DateTime(),
			}).Infoln("任务触发剩余时长")
		}
	}
	// 任务移除
	if index >= 0 {
		if size <= index+1 {
			scheduler.queue = nil
		} else {
			scheduler.queue = scheduler.queue[index+1:]
		}
		if scheduler.debug {
			scheduler.logger.WithFields(logrus.Fields{
				"index":      index,
				"size":       size,
				"queue_size": len(scheduler.queue),
			}).Infoln("任务剔除")
		}
	}

}

// 移除
func (scheduler *schedulerDomain) remove(id string) {
	delete(scheduler.cacheIndex, id)
}

// 发现新增任务
func (scheduler *schedulerDomain) discover() {
	if scheduler.discoverFunc != nil {
		if scheduler.debug {
			scheduler.logger.WithField("queue_size", len(scheduler.queue)).Infoln("任务发现...")
		}
		scheduler.state = StateDiscover
		scheduler.discoverFunc(scheduler.channel, time.Now())
		scheduler.state = StateRunning
		if scheduler.debug {
			scheduler.logger.WithField("queue_size", len(scheduler.queue)).Infoln("任务结束...")
		}
	}
}

// 派发执行任务
func (scheduler *schedulerDomain) execute(cron *entity.Crontab) {
	if scheduler.debug {
		scheduler.logger.WithField("id", cron.ID()).Infoln("执行定时任务")
	}
	if cron != nil && cron.Check() {
		// 携程异步
		var err = repo.GetPoolRepo().Add(func() {
			if err := cron.Execute(); err != nil {
				scheduler.logger.WithField("id", cron.ID()).Errorln("error:", err)
			}
		})
		// 携程执行异常
		if err != nil {
			scheduler.logger.WithField("id", cron.ID()).Errorln("execute worker error:", err)
		}
	}
}

// 缓存任务
func (scheduler *schedulerDomain) cache(cron *entity.Crontab) {
	if cron == nil || !cron.Check() {
		return
	}
	var id = cron.ID()
	if _, ok := scheduler.cacheIndex[id]; ok {
		return
	}
	scheduler.locker.Lock()
	defer scheduler.locker.Unlock()
	scheduler.queue = append(scheduler.queue, cron)
	sort.Sort(entity.CrontabItems(scheduler.queue))
	scheduler.cacheIndex[id] = true
	if scheduler.debug {
		scheduler.logger.WithField("id", cron.ID()).Infoln("任务暂存")
	}
}

// Stop 停止任务
func (scheduler *schedulerDomain) Stop() {
	if scheduler.cancel != nil && scheduler.state != StateStop {
		scheduler.cancel()
		scheduler.cancel = nil
	}
	scheduler.state = StateStop
}

func (scheduler *schedulerDomain) Add(task *entity.Crontab) {
	if scheduler.channel != nil && task != nil && task.Check() {
		scheduler.channel <- task
	}
}

func (scheduler *schedulerDomain) Refresh() {
	if scheduler.discoverFunc != nil && scheduler.state != StateDiscover {
		scheduler.notify <- StateRefresh
	}
}

func (scheduler *schedulerDomain) Push(task *entity.Crontab) {
	scheduler.Add(task)
}
