package models

import "time"

type QueryBindInfo struct {
	ID    uint   `xorm:" pk 'id'" json:"id"`
	AppID string `xorm:"'appid'" json:"appid"`
	// 绑定的消费队列
	Queue string `xorm:"'queue'" json:"queue"`
	// 消费器
	Consumer string `xorm:"'consumer'" json:"consumer"`
	// 状态  1:绑定,2:解绑
	Status uint `xorm:"'status'" json:"status"`
	// 消费器配置 json
	Properties string    `xorm:"'properties'" json:"properties"`
	UpdatedAt  time.Time `xorm:" updated 'updated_at'" json:"-"`
	CreatedAt  time.Time `xorm:" created 'created_at'" json:"-"`
	baseModel
}

func (info *QueryBindInfo) TableName() string {
	if info.table == "" {
		info.setTable("app_queue_binding")
	}
	return info.baseModel.TableName()
}


