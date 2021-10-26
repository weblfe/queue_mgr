package facede

import "github.com/weblfe/drivers/rabbitmq"

type QueueEntry interface {
	Stop()
	QueueDeclare(queue string, options ...func(params interface{})) error
	Push(data interface{}, queue ...string) error
	Pop(callback func(broker rabbitmq.MessageWrapper), queue ...string) error
}
