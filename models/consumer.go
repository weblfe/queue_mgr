package models

import "time"

type ConsumerInfo struct {
	ID    uint   `xorm:" pk 'id'" json:"id"`
	AppId string `xorm:"'appid'" json:"appid"`
	// 队列状态 0:未启动消费,1:消费中,2:idle(空闲),3:暂停消费
	Status uint `xorm:"'status'" json:"status"`
	// 类型 fastcgi,api,grpc,shell,ws,mysql
	Type string `xorm:"'type'" json:"type"`
	// 消费器名
	Name string `xorm:"'name'" json:"name"`
	// 消费器配置json
	Properties string `xorm:"'properties'" json:"properties"`
	// 备注说明
	Comment   string    `xorm:"'comment'" json:"comment"`
	UpdatedAt time.Time `xorm:" updated 'updated_at'" json:"-"`
	CreatedAt time.Time `xorm:" created 'created_at'" json:"-"`
	baseModel
}

func (info *ConsumerInfo) TableName() string {
	if info.table == "" {
		info.setTable("app_queue_consumers")
	}
	return info.baseModel.TableName()
}
