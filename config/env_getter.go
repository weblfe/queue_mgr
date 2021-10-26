package config

import (
	"github.com/weblfe/queue_mgr/utils"
	"strings"
	"time"
)

type EnvGetter struct {
	key       string
	separator string
	index     int
	Upper     bool // 是否自动转大写 key
}

func NewEnvKvGetter(key, separator string, upper ...bool) *EnvGetter {
	upper = append(upper, false)
	return &EnvGetter{
		key, separator, 0, upper[0],
	}
}

func (env *EnvGetter) SetIndex(i int) *EnvGetter {
	env.index = i
	return env
}

func (env *EnvGetter) Get(key string, def ...string) string {
	return utils.GetEnvVal(env.sKey(key), def...)
}

func (env *EnvGetter) GetInt(key string, def ...int) int {
	return utils.GetEnvInt(env.sKey(key), def...)
}

func (env *EnvGetter) GetBool(key string, def ...bool) bool {
	return utils.GetEnvBool(env.sKey(key), def...)
}

func (env *EnvGetter) GetDuration(key string, def ...time.Duration) time.Duration {
	return utils.GetEnvDuration(env.sKey(key), def...)
}

func (env *EnvGetter) GetFloat(key string, def ...float64) float64 {
	return utils.GetEnvFloat(env.sKey(key), def...)
}

func (env *EnvGetter) GetTime(key string, def ...*time.Time) *time.Time {
	return utils.GetEnvTime(env.sKey(key), def...)
}

func (env *EnvGetter) GetMapper(key string, def ...map[string]interface{}) map[string]interface{} {
	return utils.GetEnvMapper(env.sKey(key), def...)
}

func (env *EnvGetter) sKey(key string) string {
	if env.separator == "" || env.key == "" {
		return env.format(key)
	}
	if env.separator != "" && env.key != "" {
		if !strings.Contains(env.key, env.separator) {
			return env.format(env.key + env.separator + key)
		}
		if strings.HasPrefix(key, env.separator) {
			return env.format(env.key + key)
		}
		arr := strings.Split(key, env.separator)
		val := env.key + env.separator
		return env.format(strings.Join(arr[env.index:], val+strings.Join(arr[env.index+1:], env.separator)))
	}
	return env.format(key)
}

func (env *EnvGetter) format(key string) string {
	if env.Upper && key != "" {
		return strings.ToUpper(key)
	}
	return key
}

func (env *EnvGetter) GetArr(key string, def ...[]string) []string {
	return utils.GetEnvArr(env.sKey(key), def...)
}
