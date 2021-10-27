package domain

import (
	"github.com/gofiber/fiber/v2"
	"github.com/weblfe/queue_mgr/entity"
	"github.com/weblfe/queue_mgr/facede"
	"runtime"
	"sync"
)

type (

	FastCgiCreator func(ctx *fiber.Ctx) (facede.Handler, error)

	fastCgiMgrDomainImpl struct {
		constructor sync.Once
		safe        sync.RWMutex
		fastCgiPool map[entity.FastCgiType]sync.Pool
		factories   map[entity.FastCgiType]FastCgiCreator
	}

)

var (
	fastCgiDefault = NewFastCgiDomain()
)

const (
	FastCGIQueueType = "fastCGI"
	HeaderQueueType  = "Queue-Handler-Type"
)

func NewFastCgiDomain() *fastCgiMgrDomainImpl {
	var domain = new(fastCgiMgrDomainImpl)
	return domain.init()
}

func GetFastCGIDomain() *fastCgiMgrDomainImpl {
	return fastCgiDefault
}

func (domain *fastCgiMgrDomainImpl) init() *fastCgiMgrDomainImpl {
	domain.constructor = sync.Once{}
	domain.fastCgiPool = make(map[entity.FastCgiType]sync.Pool)
	domain.factories = make(map[entity.FastCgiType]FastCgiCreator)
	runtime.SetFinalizer(domain, (*fastCgiMgrDomainImpl).destroy)
	return domain
}

func (domain *fastCgiMgrDomainImpl) Resolve(ctx *fiber.Ctx) (facede.Handler, error) {
	var tyName = ctx.Get(HeaderQueueType, FastCGIQueueType)
	return domain.get(entity.FastCgiType(tyName), ctx)
}

func (domain *fastCgiMgrDomainImpl) get(name entity.FastCgiType, ctx *fiber.Ctx) (facede.Handler, error) {
	// 池化处理获取
	if pool := domain.GetPool(name); pool != nil {
		if creator := pool.Get(); creator != nil {
			switch creator.(type) {
			case FastCgiCreator:
				return creator.(FastCgiCreator)(ctx)
			case facede.Handler:
				return creator.(facede.Handler), nil
			}
		}
	}
	// 特殊处理, 非池化
	if creator := domain.GetFactory(name); creator != nil {
		if creator == nil {
			return nil, entity.ErrorEmpty
		}
		return creator(ctx)
	}
	return nil, entity.ErrorSupport
}

func (domain *fastCgiMgrDomainImpl) GetPool(name entity.FastCgiType) *sync.Pool {
	if pool, ok := domain.fastCgiPool[name]; ok {
		return &pool
	}
	return nil
}

func (domain *fastCgiMgrDomainImpl) GetFactory(name entity.FastCgiType) FastCgiCreator {
	if factory, ok := domain.factories[name]; ok && factory != nil {
		return factory
	}
	return nil
}

func (domain *fastCgiMgrDomainImpl) AddPool(name entity.FastCgiType, creator FastCgiCreator) bool {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if _, ok := domain.factories[name]; ok {
		if _, ok2 := domain.fastCgiPool[name]; ok2 {
			return true
		}
		domain.fastCgiPool[name] = sync.Pool{
			New: domain.createFactory(name),
		}
		return true
	}
	domain.factories[name] = creator
	domain.fastCgiPool[name] = sync.Pool{New: domain.createFactory(name)}
	return true
}

func (domain *fastCgiMgrDomainImpl) RegisterFactory(name entity.FastCgiType, creator FastCgiCreator) bool {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if fn, ok := domain.factories[name]; ok && fn != nil {
		return true
	}
	domain.factories[name] = creator
	return true
}

func (domain *fastCgiMgrDomainImpl) createFactory(name entity.FastCgiType) func() interface{} {
	var creator, ok = domain.factories[name]
	if !ok {
		return nil
	}
	return func() interface{} {
		return creator
	}
}

func (domain *fastCgiMgrDomainImpl) destroy() {
	defer runtime.SetFinalizer(domain, nil)
	domain.factories = nil
	domain.fastCgiPool = nil
}
