package starter

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/weblfe/queue_mgr/repo"
	"os"
	"sync"
	"time"
)

type localStorageStarter struct {
	dbPath string
	lock   sync.RWMutex
	baseStarterConstructor
}

var (
	localStorage = newLocalStorageStarter()
)

func newLocalStorageStarter() *localStorageStarter {
	var starter = new(localStorageStarter)
	starter.dbPath = ""
	starter.lock = sync.RWMutex{}
	starter.baseStarterConstructor = newStarterConstructor()
	starter.name = "localStorage"
	return starter
}

func GetLocalStorageStarter() *localStorageStarter {
	return localStorage
}

func (starter *localStorageStarter) StartUp() {
	starter.init(starter.boot)
}

func (starter *localStorageStarter) boot() {
	var (
		dbs     = starter.getDbs()
		factory = repo.GetLocalStorageRepo()
	)
	if len(dbs) <= 0 {
		return
	}
	for _, v := range dbs {
		var db, err = starter.openDb(v)
		if err != nil {
			log.Infoln(fmt.Sprintf("%s, localStorage Db:%s, open error:%s", time.Now().Format(`2006-01-02 15:04:05`), v, err.Error()))
		}
		if db != nil {
			factory.RegisterDb(v, db)
		}
	}
}

func (starter *localStorageStarter) GetConnUrl(name ...string) string {
	name = append(name, "default")
	var key = name[0]
	if starter.dbPath != "" && key == "default" {
		return starter.dbPath
	}

	return defaultConnUrl
}

func (starter *localStorageStarter) getDbs() []string {
	var keys []string

	return keys
}

func (starter *localStorageStarter) openDb(name string) (*leveldb.DB, error) {
	var (
		file = starter.GetConnUrl(name)
	)
	if file == "" {
		return nil, errors.New("db file not configured")
	}
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}
	return leveldb.OpenFile(file, starter.GetOptions(name))
}

func (starter *localStorageStarter) GetOptions(name ...string) *opt.Options {
	return &opt.Options{}
}
