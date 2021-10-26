package starter

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
	"sync"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type databaseStarter struct {
	engines map[string]*xorm.Engine
	lock    sync.RWMutex
	baseStarterConstructor
}

const (
	defaultDriver  = "mysql"
	defaultEntity  = "default"
	defaultConnUrl = "root:123@/test?charset=utf8"
)

var (
	database = newDatabaseStarter()
)

func newDatabaseStarter() *databaseStarter {
	var starter = new(databaseStarter)
	starter.lock = sync.RWMutex{}
	starter.engines = make(map[string]*xorm.Engine)
	starter.baseStarterConstructor = newStarterConstructor()
	starter.name = "database"
	return starter
}

func GetDataBaseStarter() *databaseStarter {
	return database
}

func (starter *databaseStarter) StartUp() {
	starter.init(starter.boot)
}

func (starter *databaseStarter) boot() {
	if _, err := starter.GetDb(); err != nil {
		panic(err)
	}
}

func (starter *databaseStarter) GetDb(name ...string) (*xorm.Engine, error) {
	name = append(name, defaultEntity)
	starter.lock.Lock()
	defer starter.lock.Unlock()
	if v, ok := starter.engines[name[0]]; ok {
		return v, nil
	}
	var engine, err = xorm.NewEngine(starter.GetDriver(name...), starter.GetConnUrl(name...))

	if err == nil && engine != nil {
		starter.setting(name[0], engine)
		starter.engines[name[0]] = engine
	}
	return engine, err
}

func (starter *databaseStarter) setting(name string, engine *xorm.Engine) {
	starter.setPrefix(starter.getPrefix(name), engine)
}

func (starter *databaseStarter) setPrefix(prefix string, engine *xorm.Engine) {
	if prefix == "" {
		return
	}
	var (
		prefixMapper = names.NewPrefixMapper(names.SameMapper{}, prefix)
	)
	engine.SetTableMapper(prefixMapper)
}

func (starter *databaseStarter) setSnake(engine *xorm.Engine, prefix ...string) {
	if len(prefix) <= 0 {
		prefix = append(prefix, "")
	}
	var (
		prefixMapper = names.NewPrefixMapper(names.SnakeMapper{}, prefix[0])
	)
	engine.SetTableMapper(prefixMapper)
}

func (starter *databaseStarter) getPrefix(name string) string {
	if name == "" || name == "default" {
		if prefix := os.Getenv("DB_PREFIX"); prefix != "" {
			return prefix
		}
		return ""
	}
	if prefix := os.Getenv(fmt.Sprintf("DB_%s_PREFIX", name)); prefix != "" {
		return prefix
	}
	return ""
}

func (starter *databaseStarter) GetDriver(name ...string) string {
	driver := os.Getenv("DB_DRIVER")
	if driver != "" {
		return strings.ToLower(driver)
	}
	return defaultDriver
}

func (starter *databaseStarter) GetConnUrl(name ...string) string {
	return defaultConnUrl
}

func (starter *databaseStarter)env(key string,def...string) string  {
	return ""
}
