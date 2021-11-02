package models

import (
	"github.com/weblfe/queue_mgr/entity"
	"time"
	"xorm.io/builder"
)

type QueueInfo struct {
	ID    uint   `xorm:" pk 'id'" json:"id"`
	AppID string `xorm:"'appid'" json:"appid"`
	// 队列状态 0:未启动消费,1:消费中,2:idle(空闲),3:暂停消费
	Status uint `xorm:"'status'" json:"status"`
	// 类型 amqp,mqtt,redis,native
	Type string `xorm:"'type'" json:"type"`
	// 队列名
	Name string `xorm:"'name'" json:"name"`
	// 消费者数量
	ConsumerMaxNum uint `xorm:"'consumer_max_num'" json:"consumer_max_num"`
	// 队列配置json
	Properties string `xorm:"'properties'" json:"properties"`
	// 备注说明
	Comment   string    `xorm:"'comment'" json:"comment"`
	UpdatedAt time.Time `xorm:" updated 'updated_at'" json:"-"`
	CreatedAt time.Time `xorm:" created 'created_at'" json:"-"`
	baseModel
}

func (info *QueueInfo) TableName() string {
	if info.table == "" {
		info.setTable("app_queues")
	}
	return info.baseModel.TableName()
}

func (info *QueueInfo) Create(params entity.QueueParams) error {
	return nil
}

func (info *QueueInfo) GetByCond(params builder.Cond) (*QueueInfo, error) {
	return nil, nil
}

func (info *QueueInfo)GetBinding() *QueryBindInfo  {
		return nil
}