package starter

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type loggerStarter struct {
	cache map[string]*logrus.Entry
	lock  sync.RWMutex
	baseStarterConstructor
}

var defaultLoggerStarter = newLoggerStarter()

func newLoggerStarter() *loggerStarter {
	var logger = new(loggerStarter)
	logger.lock = sync.RWMutex{}
	logger.cache = make(map[string]*logrus.Entry)
	logger.baseStarterConstructor = newStarterConstructor()
	logger.name = "loggerStarter"
	return logger
}

func (starter *loggerStarter) StartUp() {
	starter.init(starter.boot)
}

func (starter *loggerStarter) boot() {

}

func GetLoggerStarter() *loggerStarter {
	return defaultLoggerStarter
}