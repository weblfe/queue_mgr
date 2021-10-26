package domain

import (
	"github.com/weblfe/queue_mgr/entity"
	"runtime"
	"sync"
)

// 队列管理器器
type queueMgrDomainImpl struct {
	safe        sync.RWMutex
	processPool map[string]sync.Pool
	queueTree   map[string]*entity.QueueOptions
}

func NewQueueMgrDomain() *queueMgrDomainImpl {
	var domain = new(queueMgrDomainImpl)
	return domain.init()
}

func (domain *queueMgrDomainImpl) init() *queueMgrDomainImpl {
	domain.safe = sync.RWMutex{}
	domain.processPool = make(map[string]sync.Pool)
	domain.queueTree = make(map[string]*entity.QueueOptions)
	runtime.SetFinalizer(domain, (*queueMgrDomainImpl).destroy)
	return domain
}

func (domain *queueMgrDomainImpl) Register(opts *entity.QueueOptions) {
	if opts == nil {
		return
	}
	domain.safe.Lock()
	defer domain.safe.Unlock()
	domain.queueTree[opts.Name] = opts
}

func (domain *queueMgrDomainImpl) GetQueue(name string) *entity.QueueOptions {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if queue, ok := domain.queueTree[name]; ok {
		return queue
	}
	return nil
}

func (domain *queueMgrDomainImpl) AddPool(name string, creator func() interface{}) *queueMgrDomainImpl {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if _, ok := domain.queueTree[name]; ok {
		if _, ok = domain.processPool[name]; ok {
			return domain
		}
		domain.processPool[name] = sync.Pool{New: creator}
	}
	return domain
}

func (domain *queueMgrDomainImpl) GetProcessor(name string) (interface{}, bool) {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if pool, ok := domain.processPool[name]; ok {
		if p := pool.Get(); p != nil {
			return p, true
		}
	}
	return nil, false
}

func (domain *queueMgrDomainImpl) destroy() {
	defer runtime.SetFinalizer(domain, nil)
	domain.queueTree = nil
	domain.processPool = nil
}
