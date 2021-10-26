package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/weblfe/queue_mgr/facede"
	"strings"
	"sync"
)

// 应用相关配置
type applicationConfiguration struct {
	constructor sync.Once
	lock        sync.RWMutex
	container   map[string]facede.CfgKv
	providers   map[string]func(interface{}) facede.CfgKv
}

var appConfig = newConfiguration()

func newConfiguration() *applicationConfiguration {
	var app = new(applicationConfiguration)
	app.lock = sync.RWMutex{}
	app.constructor = sync.Once{}
	app.container = make(map[string]facede.CfgKv)
	app.providers = make(map[string]func(interface{}) facede.CfgKv)
	return providers(app)
}

// 注册 配置构建器
func providers(app *applicationConfiguration) *applicationConfiguration {
	registerDbFactory(app)
	registerAppCfgFactory(app)
	registerRedisCfgFactory(app)
	registerServiceKvFactory(app)
	registerLocalStorageKvFactory(app)
	return app
}

// GetAppConfig 获取应用配置
func GetAppConfig() *applicationConfiguration {
	return appConfig
}

// LoadConfiguration 加载配置
func (l *applicationConfiguration) LoadConfiguration(app *viper.Viper) bool {
	l.constructor.Do(func() {
		l.init(app)
	})
	return true
}

// 初始化
func (l *applicationConfiguration) init(cfg *viper.Viper) {
	var (
		keyArr = cfg.AllKeys()
		cache  = map[string]bool{}
	)

	for _, k := range keyArr {
		if strings.Contains(k, ".") {
			arr := strings.Split(k, ".")
			k = arr[0]
		}
		if _, ok := cache[k]; ok {
			continue
		}
		v := cfg.Get(k)
		factory, ok := l.providers[k]
		if !ok {
			continue
		}
		kv := factory(v)
		if kv == nil {
			continue
		}
		l.Add(k, kv)
		cache[k] = true
	}
}

// Add 添加 配置组Kv
func (l *applicationConfiguration) Add(key string, v facede.CfgKv) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.container[key]; ok {
		return
	}
	l.container[key] = v
}

// GetKvObj 获取配置对象kv
func (l *applicationConfiguration) GetKvObj(key string) (facede.CfgKv, bool) {
	l.lock.RLocker().Lock()
	defer l.lock.RLocker().Unlock()
	if v, ok := l.container[key]; ok && v != nil {
		return v, true
	}
	return nil, false
}

// ValueOf 获取配置值
func (l *applicationConfiguration) ValueOf(key string, def ...interface{}) interface{} {
	def = append(def, nil)
	l.lock.RLocker().Lock()
	defer l.lock.RLocker().Unlock()
	if !strings.Contains(key, ".") {
		if v, ok := l.container[key]; ok && v != nil {
			return v
		}
		return def[0]
	}
	var arr = strings.Split(key, ".")
	if v, ok := l.container[arr[0]]; ok && v != nil {
		return v.ValueOf(strings.Join(arr[1:], ""), def[0])
	}
	return def[0]
}

// Register 注册
func (l *applicationConfiguration) Register(key string, factory func(interface{}) facede.CfgKv) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.providers[key]; ok {
		return
	}
	l.providers[key] = factory
}

// GetDatabase 获取数据连接配置
func (l *applicationConfiguration) GetDatabase(name ...string) Database {
	name = append(name, "")
	var db, ok = l.GetKvObj("database")
	if !ok {
		return CreateDatabase(name[0])
	}
	dbKv, ok := db.(*DatabaseKv)
	if !ok {
		return CreateDatabase(name[0])
	}
	d, _ := dbKv.Get(name[0])
	return d
}

// GetRedis 获取redis 配置
func (l *applicationConfiguration) GetRedis(name ...string) Redis {
	name = append(name, "")
	var db, ok = l.GetKvObj("redis")
	if !ok {
		return CreateRedis(name[0])
	}
	rdKv, ok := db.(*RedisKv)
	if !ok || rdKv == nil {
		return CreateRedis(name[0])
	}
	rd, _ := rdKv.Get(name[0])
	return rd
}

// GetService  服务配置
func (l *applicationConfiguration) GetService(name ...string) Service {
	name = append(name, "")
	var serv, ok = l.GetKvObj("services")
	if !ok {
		return *CreateService(name[0])
	}
	servKv, ok := serv.(*ServiceKv)
	if !ok {
		return *CreateService(name[0])
	}
	service, _ := servKv.Get(name[0])
	if service == nil {
		logrus.Info(name[0], " service not found")
		return *CreateService(name[0])
	}
	return *service
}

// GetServicesCnf  服务配置组
func (l *applicationConfiguration) GetServicesCnf(name ...string) ServiceKv {
	name = append(name, "")
	var serv, ok = l.GetKvObj("services")
	if !ok {
		return nil
	}
	servKv, ok := serv.(*ServiceKv)
	if !ok && servKv != nil {
		return nil
	}
	return *servKv
}

// GetLocalStorageCnf   本地存储配置组
func (l *applicationConfiguration) GetLocalStorageCnf(name ...string) LocalStorageKv {
	name = append(name, "")
	var serv, ok = l.GetKvObj("localStorage")
	if !ok {
		return nil
	}
	servKv, ok := serv.(LocalStorageKv)
	if !ok || servKv == nil {
		return nil
	}
	return servKv
}

// GetAppCnf 获取应用全局 配置
func (l *applicationConfiguration) GetAppCnf() App {

	var app, ok = l.GetKvObj("app")
	if !ok {
		return *CreateApp()
	}
	kv, ok := app.(*App)
	if !ok || kv == nil {
		return *CreateApp()
	}
	return *kv
}

func (l *applicationConfiguration) GetDbKv() DatabaseKv {
	var db, ok = l.GetKvObj("database")
	if !ok {
		return nil
	}
	dbKv, ok := db.(*DatabaseKv)
	if !ok || dbKv == nil {
		return nil
	}
	return *dbKv
}

func (l *applicationConfiguration) GetRedisKv() RedisKv {
	var db, ok = l.GetKvObj("redis")
	if !ok {
		return nil
	}
	rdKv, ok := db.(*RedisKv)
	if !ok || rdKv == nil {
		return nil
	}
	return *rdKv
}
