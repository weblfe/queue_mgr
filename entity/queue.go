package entity

import "github.com/weblfe/queue_mgr/utils"

type (
	// QueueOptions 队列选项
	QueueOptions struct {
		Name         string `json:"name,omitempty" yaml:"name"`         // 队列名
		Type         string `json:"type,omitempty" yaml:"type"`         // 队列类型
		Handler      string `json:"handler" yaml:"handler"`             // 处理器
		ConnUrl      string `json:"conn_url" yaml:"conn_url"`           // amqp,mqtt 链接url
		Queue        string `json:"queue" yaml:"queue"`                 // 队列名
		Topic        string `json:"topic" yaml:"topic"`                 // 主题名
		ErrorHandler string `json:"error_handler" yaml:"error_handler"` // 异常处理器
	}
)

func NewQueueOption() *QueueOptions {
	var options = new(QueueOptions)
	return options
}

func DecoderOptions(bytes []byte) *QueueOptions {
	if len(bytes) <= 0 {
		return nil
	}
	var opts = NewQueueOption()
	if err := utils.JsonDecode(bytes, opts); err != nil {
		return nil
	}
	return opts
}

func (options *QueueOptions) GetTopic() string {
	return options.Topic
}

func (options *QueueOptions) GetQueue() string {
	return options.Queue
}

func (options *QueueOptions) GetType() string {
	return options.Type
}
