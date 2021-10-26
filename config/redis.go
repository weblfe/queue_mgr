package config

import (
		"fmt"
		"github.com/weblfe/queue_mgr/facede"
		"github.com/weblfe/queue_mgr/utils"
		"strings"
)

type Redis struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	Auth    string `json:"password"`
	Prefix  string `json:"prefix"`
	Db      int    `json:"db"`
	Default bool   `json:"default"`
}

type RedisKv map[string]Redis

func CreateRedis(v ...interface{}) Redis {
	var key = ""
	if len(v) >= 0 && v[0] != nil {
		switch v[0].(type) {
		case string:
			str := v[0].(string)
			if str != "default" {
				key = str
			}
		case map[string]interface{}:
			m := v[0].(map[string]interface{})
			return Redis{
				Port:    utils.MGetInt(m, "port", 6379),
				Host:    utils.MGet(m, "host", "127.0.0.1"),
				Auth:    utils.MGet(m, "auth", ""),
				Prefix:  utils.MGet(m, "prefix", ""),
				Db:      utils.MGetInt(m, "db", 0),
				Default: utils.MGetBool(m, "default", false),
			}
		}
	}
	var env = NewEnvKvGetter(key, "_",true)
	return Redis{
		Port:    env.GetInt("REDIS_PORT", 6379),
		Host:    env.Get("REDIS_HOST", "127.0.0.1"),
		Auth:    env.Get("REDIS_AUTH", ""),
		Prefix:  env.Get("REDIS_PREFIX", ""),
		Db:      env.GetInt("REDIS_DB", 0),
		Default: env.GetBool("REDIS_DEFAULT", false),
	}
}

func (r *Redis) Get(key string, def ...interface{}) interface{} {
	switch strings.ToLower(key) {
	case "host":
		return r.Host
	case "port":
		return r.Port
	case "password":
		return r.Auth
	case "auth":
		return r.Auth
	case "db":
		return r.Db
	case "prefix":
		return r.Prefix
	case "default":
		return r.Default
	}
	return def[0]
}

func (r *Redis)GetAddr() string  {
		return fmt.Sprintf("%s:%d",r.Host,r.Port)
}

func (r *Redis)GetPoolSize() int  {
		return utils.GetEnvInt("REDIS_POOL_SIZE",100)
}

func (r *Redis)GetMinIdleConns() int  {
		return utils.GetEnvInt("REDIS_MIN_IDLE_CONNS",3)
}

func (data *RedisKv) Get(key string) (Redis, bool) {
	if v, ok := (*data)[key]; ok {
		return v, true
	}
	return CreateRedis(key), false
}

func (data *RedisKv) Size() int {
	return len(*data)
}

func (data *RedisKv) Keys() []string {
	var keys []string
	for k := range *data {
		keys = append(keys, k)
	}
	return keys
}

func (data *RedisKv) GetDefaultKey() string {
	for k, v := range *data {
		if v.Default {
			return k
		}
	}
	return "default"
}

func (data *RedisKv) Exists(key string) bool {
	if _, ok := (*data)[key]; ok {
		return true
	}
	return false
}

func (data *RedisKv) Add(key string, database Redis) *RedisKv {
	if _, ok := (*data)[key]; ok {
		return data
	}
	(*data)[key] = database
	return data
}

func (data *RedisKv) String() string {
	return utils.JsonEncode(data).String()
}

func (data *RedisKv) Decode(content []byte) error {
	return utils.JsonDecode(content, data)
}

func (data *RedisKv) ValueOf(key string, def ...interface{}) interface{} {
	if len(def) == 0 {
		def = append(def, nil)
	}
	if v, ok := (*data)[key]; ok {
		return v
	}
	return def[0]
}

func (data *RedisKv) MAdd(mArr map[string]interface{}) int {
	for k, v := range mArr {
		data.Add(k, CreateRedis(v))
	}
	return 0
}

// 创建redis 配置
func createRedisKv(v interface{}) *RedisKv {
	switch v.(type) {
	case map[string]interface{}:
		kv := RedisKv{}
		kv.MAdd(v.(map[string]interface{}))
		return &kv
	}
	return nil
}

// 注册
func registerRedisCfgFactory(app *applicationConfiguration) {
	app.Register("redis", func(v interface{}) facede.CfgKv {
		return createRedisKv(v)
	})
}
