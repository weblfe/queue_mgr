package service

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/drivers/rabbitmq"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/repo"
	"sync"
)

type (
	// 队列处理器
	queueProcessorServiceImpl struct {
		maxNum        int
		queue         string
		logger        *logrus.Logger
		locker        sync.RWMutex
		processorList []facede.QueueEntry
	}
)

var (
	queueProcessor = NewQueueProcessorService()
)

func GetQueueProcessor() *queueProcessorServiceImpl {
	return queueProcessor
}

func NewQueueProcessorService() *queueProcessorServiceImpl {
	var service = queueProcessorServiceImpl{
		maxNum: 3,
		locker: sync.RWMutex{},
	}
	return &service
}

func (queue *queueProcessorServiceImpl) SetProcessorCap(num int) *queueProcessorServiceImpl {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	queue.maxNum = num
	return queue
}

func (queue *queueProcessorServiceImpl) GetCap() int {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	return queue.maxNum
}

func (queue *queueProcessorServiceImpl) GetSize() int {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	return len(queue.processorList)
}

func (queue *queueProcessorServiceImpl) GetQueue() string {
	if queue.queue == "" {
		queue.queue = "default"
	}
	return queue.queue
}

// SetQueue 在 添加队列实例前设置队列名
func (queue *queueProcessorServiceImpl) SetQueue(name string) *queueProcessorServiceImpl {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	if queue.queue == "" {
		queue.queue = name
	}
	return queue
}

func (queue *queueProcessorServiceImpl) Close() {
	queue.locker.Lock()
	defer queue.locker.Unlock()
	for _, entry := range queue.processorList {
		entry.Stop()
	}
	queue.processorList = nil
}

func (queue *queueProcessorServiceImpl) Add(entry facede.QueueEntry, callback func(broker rabbitmq.MessageWrapper)) error {
	var capacity = queue.GetCap()
	if queue.GetSize() >= capacity {
		return errors.New(fmt.Sprintf("queueProcessorServiceImpl cap<%d> limited", capacity))
	}
	queue.locker.Lock()
	defer queue.locker.Unlock()
	queue.processorList = append(queue.processorList, entry)
	queueName := queue.GetQueue()
	return repo.GetPoolRepo().Add(func() {
		if err := entry.Pop(callback, queueName); err != nil {
			queue.getLogger().WithField("error", err).Errorln(err)
		} else {
			queue.getLogger().Infoln(queueName, ".queue->stop")
		}
	})
}

func (queue *queueProcessorServiceImpl) getLogger() *logrus.Logger {
	if queue.logger == nil {
		queue.logger = logrus.New()
	}
	return queue.logger
}

func (queue *queueProcessorServiceImpl) SetLogger(log *logrus.Logger) *queueProcessorServiceImpl {
	if queue.logger == nil {
		queue.logger = log
	}
	return queue
}
