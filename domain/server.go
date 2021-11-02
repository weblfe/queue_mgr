package domain

import (
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/models"
	"github.com/weblfe/queue_mgr/repo"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type (
	Queue string
	State string

	serverDomainImpl struct {
		safe          sync.RWMutex
		refreshTicker *time.Ticker
		quit          chan os.Signal
		// 队列状态树
		trees map[State]*stateTree
		// 队列信息存储
		queues map[Queue]entity.QueueState
	}

	// 队列状态数
	stateTree struct {
		safe    sync.RWMutex
		items   []*QueueInfo
		indexes map[string]int
		state   entity.QueueState
	}
)

func NewServDomain() *serverDomainImpl {
	var servImpl = new(serverDomainImpl)
	return servImpl.init()
}

func (serv *serverDomainImpl) init() *serverDomainImpl {
	serv.safe = sync.RWMutex{}
	serv.quit = make(chan os.Signal)
	serv.trees = make(map[State]*stateTree)
	serv.refreshTicker = time.NewTicker(time.Second)
	serv.queues = make(map[Queue]entity.QueueState)
	return serv
}

func (serv *serverDomainImpl) Register(queue *models.QueueInfo) bool {
	if queue == nil {
		return false
	}
	var bind = queue.GetBinding()
	return serv.add(queue, bind)
}

func (serv *serverDomainImpl) add(base *models.QueueInfo, bind *models.QueryBindInfo) bool {
	if base == nil {
		return false
	}
	serv.safe.Lock()
	defer serv.safe.Unlock()
	var (
		queue = &QueueInfo{
			Base: base,
			Bind: bind,
		}
		name = queue.Queue()
	)
	if name == "" {
		return false
	}
	if state, ok := serv.queues[Queue(name)]; ok {
		if !state.Is(base.Status) {
			return false
		}
		if state != entity.QueueState(base.Status) {
			treeContainer(serv.trees).Remove(state, name)
		}
	}
	var (
		res   bool
		state entity.QueueState
	)
	if !state.Is(base.Status) {
		return false
	}
	state = entity.QueueState(base.Status)
	res = treeContainer(serv.trees).Register(state, queue)
	if res {
		serv.queues[Queue(name)] = state
	}
	return res
}

// Observe 开始前观察
func (serv *serverDomainImpl) Observe() error {
	return repo.GetPoolRepo().Add(serv.refresh)
}

// 刷新
func (serv *serverDomainImpl) refresh() {
	signal.Notify(serv.quit, syscall.SIGINT, syscall.SIGTERM)
	for {

		select {
		case <-serv.refreshTicker.C:
			serv.up()
		case <-serv.quit:
			serv.refreshTicker.Stop()
			return
		}
	}
}

// 刷新队列
func (serv *serverDomainImpl) up() {
	// 发现新
	serv.discover()
	// 运行的
	serv.run()
	// 空闲
	serv.idle()

}

func (serv *serverDomainImpl) discover() {

}

func (serv *serverDomainImpl) idle() {

}

func (serv *serverDomainImpl) run() {

}
