package domain

import (
	"github.com/gofiber/fiber/v2"
	"github.com/weblfe/queue_mgr/entity"
	"sync"
	"time"
)

type PHPFastCgiDomainImpl struct {
	params entity.KvMap
}

func NewPHPFastCgiDomain() *PHPFastCgiDomainImpl {
	var domain = new(PHPFastCgiDomainImpl)
	domain.init()
	return domain
}

func (domain *PHPFastCgiDomainImpl) init() {
	domain.params = entity.KvMap{}
}

func (domain *PHPFastCgiDomainImpl) Parse() error {
	return nil
}

func (domain *PHPFastCgiDomainImpl) reset() *PHPFastCgiDomainImpl {
	domain.params = entity.KvMap{}
	return domain
}

func (domain *PHPFastCgiDomainImpl) Register(pool *sync.Pool) {
	if pool == nil {
		return
	}
	pool.Put(domain.reset())
}

func (domain *PHPFastCgiDomainImpl) Cancel() {
	panic("implement me")
}

func (domain *PHPFastCgiDomainImpl) Type() string {
	panic("implement me")
}

func (domain *PHPFastCgiDomainImpl) Handle(ctx *fiber.Ctx) error {
	panic("implement me")
}

func (domain *PHPFastCgiDomainImpl) SetTimeout(duration time.Duration) {
	domain.params.Add("timeout",duration)
}
