package repo

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/weblfe/drivers/rabbitmq"
	"sync"
)

type RabbitmqUtils struct {
	client          *rabbitmq.Client
	params          rabbitmq.PubSubParams
	container       sync.Map
	locker          sync.RWMutex
	consumerOptions []func(params interface{})
	ctrl            chan bool
}

func RabbitmqOf(namespace ...string) RabbitmqUtils {
	namespace = append(namespace, "")
	return RabbitmqUtils{
		params: rabbitmq.PubSubParams{
			Ctx:     context.Background(),
			ConnUrl: "",
			Cfg:     nil,
			Entry:   namespace[0],
		},
		container:       sync.Map{},
		locker:          sync.RWMutex{},
		ctrl:            make(chan bool, 2),
		consumerOptions: []func(params interface{}){},
	}
}

// CreateGroupRabbitmq 批量创建自动ack 消费队列
func CreateGroupRabbitmq(num int, namespace ...string) []*RabbitmqUtils {
	var queueArr []*RabbitmqUtils
	for i := 0; i < num; i++ {
		entry := RabbitmqOf(namespace...)
		queueArr = append(queueArr, &entry)
	}
	return queueArr
}

// CreateGroupAckRabbitmq  手动Ack队列批量创建
func CreateGroupAckRabbitmq(num int, namespace ...string) []*RabbitmqUtils {
	var queueArr []*RabbitmqUtils
	for i := 0; i < num; i++ {
		entry := RabbitmqOf(namespace...)
		entry.AddConsumerParamsOptions(WithAckConsumeOption)
		queueArr = append(queueArr, &entry)
	}
	return queueArr
}

func (utils *RabbitmqUtils) GetBroker() *rabbitmq.Broker {
	return rabbitmq.CreateBroker(utils.params.GetBrokerCfg())
}

func (utils *RabbitmqUtils) getClient() *rabbitmq.Client {
	if utils.client == nil {
		utils.client = rabbitmq.NewPublisher(utils.params)
	}
	return utils.client
}

// Push 推送消息
func (utils *RabbitmqUtils) Push(data interface{}, queues ...string) error {
	queues = append(queues, "default")
	queue := queues[0]
	if err := utils.QueueDeclare(queue); err != nil {
		return err
	}
	msg := rabbitmq.MessageParamsOf("", queue, data)
	return utils.getClient().Send(*msg)
}

// QueueDeclare 队列定义
func (utils *RabbitmqUtils) QueueDeclare(queue string, options ...func(params interface{})) error {
	q, ok := utils.container.Load(queue + ".declare")
	if ok && q != nil {
		return nil
	}
	if len(options) <= 0 {
		options = append(options, WithGoodQueueOptions)
	}
	var client = utils.getClient()
	if client == nil {
		return errors.New("queue client connection failed")
	}
	channel := client.GetChannel()
	params := utils.getQueueArgs(queue)
	// 设置队列参数
	for _, setting := range options {
		setting(&params)
	}
	queueIns, err := channel.QueueDeclare(
		queue,
		params.Durable,
		params.AutoDelete,
		params.Exclusive,
		params.NoWait,
		params.Args,
	)
	if err == nil {
		utils.container.Store(queue+".declare", queueIns)
		// log.Infoln("queueParams",fmt.Sprintf("%v",params))
	}
	return err
}

// 队列参数
func (utils *RabbitmqUtils) getQueueArgs(queue string) rabbitmq.QueueParams {
	return rabbitmq.QueueParams{
		Durable:    true,
		Key:        queue,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     true,
	}
}

// Pop 消费 [一个个消费]
func (utils *RabbitmqUtils) Pop(callback func(broker rabbitmq.MessageWrapper), queues ...string) error {
	queues = append(queues, "default")
	queue := queues[0]
	// 定义消费队列
	if err := utils.QueueDeclare(queue); err != nil {
		return err
	}
	// 限制值消费数量
	if err := utils.getClient().Qos(); err != nil {
		return err
	}
	utils.registerHandler(queue, callback)
	return utils.wait(queue)
}

// 阻塞
func (utils *RabbitmqUtils) wait(queue string) error {
	var (
		channel, err = utils.getClient().Receive(utils.getConsumerParams(queue))
	)
	if err != nil {
		return err
	}
	conn:=utils.getClient().GetBroker().GetConnection()
	// 消息异常
	defer func() {
		if _err := recover(); _err != nil {
			log.Infoln("queue.error", _err)
			if conn!=nil {
				conn.ConnectionState()
			}
		}
	}()
	for {
		select {
		case msg := <-channel:
			utils.dispatch(queue, msg)
		case c := <-utils.ctrl:
			if c {
				log.Infoln("stop-ctrl")
				return nil
			}
		}
	}

}

// AddConsumerParamsOptions 添加 消费参数设置处理器
func (utils *RabbitmqUtils) AddConsumerParamsOptions(options ...func(params interface{})) {
	utils.locker.Lock()
	defer utils.locker.Unlock()
	if len(options) <= 0 {
		return
	}
	utils.consumerOptions = append(utils.consumerOptions, options...)
}

// 消费参数
func (utils *RabbitmqUtils) getConsumerParams(queue string) rabbitmq.ConsumerParams {
	params, ok := utils.container.Load(queue + ".consumerParams")
	if ok {
		switch params.(type) {
		case rabbitmq.ConsumerParams:
			return params.(rabbitmq.ConsumerParams)
		}
	}
	var paramsNew = rabbitmq.ConsumerParams{
		Queue:        queue,
		Key:          "",
		Exchange:     "",
		ExchangeType: "",
		AutoAck:      false,
		Exclusive:    false,
		NoLocal:      false,
		NoWait:       false,
	}
	if len(utils.consumerOptions) < 0 {
		utils.locker.Lock()
		defer utils.locker.Unlock()
		for _, setting := range utils.consumerOptions {
			setting(&paramsNew)
		}
	}
	// _json,_:=json.Marshal(paramsNew)
	// log.Infoln("json.params:", string(_json))
	utils.container.Store(queue+".consumerParams",paramsNew)
	return paramsNew
}

// 注册消费处理
func (utils *RabbitmqUtils) registerHandler(queue string, handler func(broker rabbitmq.MessageWrapper)) {
	var (
		entry, ok = utils.container.Load(queue)
		handlers  []func(broker rabbitmq.MessageWrapper)
	)
	if !ok {
		entry = handlers
	} else {
		switch entry.(type) {
		case []func(broker rabbitmq.MessageWrapper):
			handlers = entry.([]func(broker rabbitmq.MessageWrapper))
		default:
			panic(errors.New("callbacks types error" + fmt.Sprintf("%T", entry)))
			return
		}
	}
	utils.locker.Lock()
	defer utils.locker.Unlock()
	handlers = append(handlers, handler)
	utils.container.Store(queue, handlers)
}

func (utils *RabbitmqUtils) Stop() {
	utils.ctrl <- true
}

func (utils *RabbitmqUtils) dispatch(queue string, delivery amqp.Delivery) {
	utils.locker.Lock()
	defer utils.locker.Unlock()
	var (
		entry, ok = utils.container.Load(queue)
		handlers  []func(broker rabbitmq.MessageWrapper)
	)
	if !ok {
		log.Infoln(string(delivery.Body))
		return
	}
	switch entry.(type) {
	case []func(broker rabbitmq.MessageWrapper):
		handlers = entry.([]func(broker rabbitmq.MessageWrapper))
	default:
		panic(errors.New("callbacks types error " + fmt.Sprintf("%T", entry)))
		return
	}
	for _, handler := range handlers {
		if handler == nil {
			continue
		}
		handler(rabbitmq.NewMessageWrapper(delivery))
	}
}

// WithGoodQueueOptions 推荐配置
func WithGoodQueueOptions(params interface{}) {
	if params == nil {
		return
	}
	switch params.(type) {
	case *rabbitmq.QueueParams:
		queueParams := params.(*rabbitmq.QueueParams)
		queueParams.Durable = true
		queueParams.AutoDelete = false
		queueParams.Exclusive = false
		queueParams.NoWait = true
	}
	return
}

// WithAckConsumeOption 手动ack
func WithAckConsumeOption(params interface{}) {
	if params == nil {
		return
	}
	switch params.(type) {
	case *rabbitmq.ConsumerParams:
		queueParams := params.(*rabbitmq.ConsumerParams)
		queueParams.AutoAck = false
	}
	return
}
