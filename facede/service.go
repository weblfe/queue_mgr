package facede

import "github.com/weblfe/queue_mgr/entity"

type Service interface {
	// Type 服务服务类型
	Type() entity.ServiceType
	// ServiceID 服务ID |服务名
	ServiceID() string
	// Addr 服务地址
	Addr() string
	// Desc 服务描述
	Desc() string
}
