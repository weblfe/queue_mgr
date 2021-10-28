package models

import "time"

type QueueFails struct {
	ID    uint   `xorm:" pk 'id'" json:"id"`
	AppID string `xorm:"'appid'" json:"appid"`
	// 队列状态 0: 消费失败, 1:消费异常
	Status uint `xorm:"'status'" json:"status"`
	// 已重试次数
	TryTimes uint `xorm:"'try_times'" json:"try_times"`
	// 异常信息
	Error string `xorm:"'error'" json:"error"`
	// 类型 amqp,mqtt,redis,native
	Type string `xorm:"'type'" json:"type"`
	// 队列名
	Queue    string `xorm:"'queue'" json:"queue"`
	Consumer string `xorm:"'consumer'" json:"consumer"`
	// 队列配置json
	Payloads  string    `xorm:"'payloads'" json:"payloads"`
	UpdatedAt time.Time `xorm:" updated 'updated_at'" json:"-"`
	CreatedAt time.Time `xorm:" created 'created_at'" json:"-"`
	baseModel
}

func (info *QueueFails) TableName() string {
	if info.table == "" {
		info.setTable("app_queue_fails")
	}
	return info.baseModel.TableName()
}
