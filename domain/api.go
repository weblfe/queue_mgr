package domain

import (
	"runtime"
	"sync"
)

type apiDomainImpl struct {
	safe sync.RWMutex
	pool map[string]sync.Pool
}

func NewApiDomain() *apiDomainImpl {
	var domain = new(apiDomainImpl)
	domain.init()
	return domain
}

func (domain *apiDomainImpl) init() {
	domain.safe = sync.RWMutex{}
	domain.pool = make(map[string]sync.Pool)
	runtime.SetFinalizer(domain, (*apiDomainImpl).destroy)
}

func (domain *apiDomainImpl) Get(name string) (interface{},bool) {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if v, ok := domain.pool[name]; ok {
		if q:=v.Get();q!=nil {
			return q,true
		}
	}
	return nil,false
}

func (domain *apiDomainImpl) Add(name string, newer func() interface{}) {
	domain.safe.Lock()
	defer domain.safe.Unlock()
	if _, ok := domain.pool[name]; ok {
		return
	}
	domain.pool[name] = sync.Pool{New: newer}
}

func (domain *apiDomainImpl) destroy() {
	defer runtime.SetFinalizer(domain, nil)
	domain.pool = nil
}
