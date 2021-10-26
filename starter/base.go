package starter

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type baseStarterConstructor struct {
	constructor sync.Once
	name        string
}

func (starter *baseStarterConstructor) init(initialize func(), name ...string) {
	starter.constructor.Do(initialize)
	if len(name) == 0 {
		name = append(name, starter.name)
	}
	log.Infoln(fmt.Sprintf("%s [starter.%s] StartUp Success", time.Now().Format(`2006-01-02 15:04:05`), name[0]))
}

func newStarterConstructor() baseStarterConstructor {
	return baseStarterConstructor{
		constructor: sync.Once{},
	}
}
