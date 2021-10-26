package repo

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/config"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type RedisRepository struct {
	redis.Client
}

var (
	redDb *RedisRepository
)

const (
	_redisPrefix = "redis."
	DefaultRedis = "default"
)

type RedisOptions struct {
	Addr         string `json:"addr"`
	Password     string `json:"password"`
	Db           int    `json:"db"`
	PoolSize     int    `json:"poolSize"`
	MinIdleConns int    `json:"minIdleConns"`
}

type RedisOptionMap map[string]RedisOptions

func NewRedisRepository(v ...*RedisOptions) *RedisRepository {
	if len(v) > 0 && v[0] != nil {
		return &RedisRepository{
			Client: *redis.NewClient(&redis.Options{
				Addr:         v[0].Addr,
				Password:     v[0].Password,
				DB:           v[0].Db,
				PoolSize:     v[0].PoolSize,
				MinIdleConns: v[0].MinIdleConns,
			}),
		}
	}
	var cfg = getRedisOptions("")
	return &RedisRepository{
		Client: *redis.NewClient(&redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.Db,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.MinIdleConns,
		}),
	}
}

// 获取默认redisRepo
func getRedisRepository() *RedisRepository {
	if redDb == nil {
		redDb = createRedisRepo("")
	}
	return redDb
}

// 创建redisRepo
func createRedisRepo(name string) *RedisRepository {
	var (
		key       string
		rd        = new(RedisRepository)
		container = GetContainerRepo()
	)
	if name == "" {
		key = "redis.default"
	}
	if !strings.HasPrefix(name, _redisPrefix) {
		key = GetRedisNameKey(name)
	}
	if err := container.Resolve(key, getRedisResolver(&rd)); err == nil {
		return rd
	}
	rd = NewRedisRepository(getRedisOptions(name))
	container.Cache(key, rd)
	return rd
}

// 获取解析器
func getRedisResolver(rd **RedisRepository) func(v interface{}) error {
	return func(v interface{}) error {
		if v == nil {
			return ErrorNil
		}
		if r, ok := v.(*RedisRepository); ok {
			*rd = r
			return nil
		}
		return ErrorType
	}
}

// 获取 redis option
func getRedisOptions(name string) *RedisOptions {
	var redisOpt = config.GetAppConfig().GetRedis(strings.TrimPrefix(name, _redisPrefix))
	return &RedisOptions{
		Addr:         redisOpt.GetAddr(),
		Password:     redisOpt.Auth,
		Db:           redisOpt.Db,
		PoolSize:     redisOpt.GetPoolSize(),
		MinIdleConns: redisOpt.GetMinIdleConns(),
	}
}

func GetRedisNameKey(name string) string {
	return _redisPrefix + name
}

func RedisOf() *RedisRepository {
	return getRedisRepository()
}

func RedisDb(name string) *RedisRepository {
	var (
		key    = name
		_redDb *RedisRepository
	)
	_redDb = createRedisRepo(key)
	return _redDb
}

func (rd *RedisRepository) GetString(key string) string {
	var r = rd.Get(key)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " GetString error:", err)
		return ""
	}
	return str
}

// GetMap 获取json数据,解析为hashMap
func (rd *RedisRepository) GetMap(key string) map[string]interface{} {
	var (
		data = rd.GetString(key)
		obj  = map[string]interface{}{}
	)
	if data == "" {
		return obj
	}
	if err := json.Unmarshal([]byte(data), &obj); err != nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " GetMap Error:", err.Error())
	}
	return obj
}

func (rd *RedisRepository) GetInt(key string) (int, bool) {
	var (
		data = rd.GetString(key)
	)
	if data == "" {
		return 0, false
	}
	if n, err := strconv.Atoi(data); err == nil {
		return n, true
	}
	return 0, false
}

// GetJson 获取json字符串数据,解析成为对应结构体
func (rd *RedisRepository) GetJson(key string, v interface{}) error {
	var (
		data = rd.GetString(key)
	)

	if data == "" || v == nil {
		return redis.Nil
	}
	if err := json.Unmarshal([]byte(data), v); err != nil {
		rd.getLogger().WithFields(map[string]interface{}{
			"rd.DB": rd.GetOptionDb(),
			"key":   key,
			"data":  data,
			"error": err.Error(),
		}).Infoln("redis.GetJson")
		return err
	}
	return nil
}

// GetJsonArr 获取json-array 字符串数据, 解析为字符串数组
func (rd *RedisRepository) GetJsonArr(key string) []string {
	var (
		arr = new([]string)
		r   = rd.Get(key)
	)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(err)
		return nil
	}
	if str == "" {
		return nil
	}
	err = json.Unmarshal([]byte(str), arr)
	if err != nil {
		rd.getLogger().Println("Redis.DB", rd.GetOptionDb(), key, "REDIS-Str: ", str)
		rd.getLogger().Println("Redis.DB", rd.GetOptionDb(), key, "REDIS-Error: ", err)
		return nil
	}
	return *arr
}

func (rd *RedisRepository) PushQueue(key string, v ...string) int64 {
	var (
		arr []interface{}
	)
	if len(v) <= 0 {
		return 0
	}
	for _, it := range v {
		arr = append(arr, it)
	}
	cmd := rd.LPush(key, arr...)
	n, err := cmd.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " PushQueue error:", err.Error(), v)
		return 0
	}
	return n
}

func (rd *RedisRepository) HGetString(key string, field string) string {
	var r = rd.HGet(key, field)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetString error:", err)
		return ""
	}
	return str
}

func (rd *RedisRepository) HGetInt(key string, field string) int {
	var r = rd.HGet(key, field)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetInt error:", err)
		return 0
	}
	n, errs := strconv.Atoi(str)
	if errs != nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetInt error:", err)
		return 0
	}
	return n
}

func (rd *RedisRepository) HGetBool(key string, field string) bool {
	var r = rd.HGet(key, field)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetBool error:", err)
		return false
	}
	b, errs := strconv.ParseBool(str)
	if errs != nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetBool error:", err)
		return false
	}
	return b
}

func (rd *RedisRepository) HGetFloat(key string, field string) float64 {
	var r = rd.HGet(key, field)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetFloat error:", err)
		return 0
	}
	b, errs := strconv.ParseFloat(str, 32)
	if errs != nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetFloat error:", err)
		return 0
	}
	return b
}

func (rd *RedisRepository) HGetJson(key string, field string) map[string]interface{} {
	var (
		r     = rd.HGet(key, field)
		_json = make(map[string]interface{})
	)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetJson error:", err)
		return nil
	}
	errs := json.Unmarshal([]byte(str), &_json)
	if errs != nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetJson error:", err)
		return nil
	}
	return _json
}

func (rd *RedisRepository) HGetStruct(key string, field string, v interface{}) error {
	var (
		r = rd.HGet(key, field)
	)
	if v == nil {
		return errors.New(key + "." + field + ",redis HGetStruct error: v is Nil Pointer")
	}
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, field, " HGetJson error:", err)
		return nil
	}
	errs := json.Unmarshal([]byte(str), v)
	if errs != nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " HGetJson error:", err)
		return errs
	}
	return nil
}

func (rd *RedisRepository) SetAny(key string, v interface{}, expire ...time.Duration) error {
	var status *redis.StatusCmd
	switch v.(type) {
	case string:
	case *string:
		v = *v.(*string)
	default:
		if v != "" {
			if _data, err := json.Marshal(v); err == nil {
				v = string(_data)
			}
		}
	}
	if len(expire) <= 0 {
		status = rd.Set(key, v, 0)
	} else {
		status = rd.Set(key, v, expire[0])
	}
	if _, err := status.Result(); err != nil {
		return err
	}
	return nil
}

func (rd *RedisRepository) SetArr(key string, field []string, duration ...time.Duration) bool {
	if len(duration) == 0 {
		duration = append(duration, 24*time.Hour)
	}
	var (
		data, err = json.Marshal(field)
	)
	if err != nil {
		rd.getLogger().Println(key, field, "Redis.DB", rd.GetOptionDb(), " Redis SetArr error:", err.Error())
		return false
	}
	return rd.SetEx(key, data, duration[0])
}

func (rd *RedisRepository) SetEx(key string, v interface{}, duration ...time.Duration) bool {
	if len(duration) == 0 {
		duration = append(duration, 24*time.Hour)
	}
	vf := reflect.ValueOf(v)
	if (vf.Kind() == reflect.Struct) || (vf.Kind() == reflect.Ptr && vf.Elem().Kind() == reflect.Struct) {
		_json, err := json.Marshal(v)
		if err == nil {
			v = string(_json)
		}
	}
	var r = rd.Set(key, v, duration[0])
	_, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " Redis SetEx Error: ", err.Error())
		return false
	}
	return true
}

func (rd *RedisRepository) HSet(key string, values map[string]interface{}) bool {
	var r = rd.HMSet(key, values)
	_, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " HSet error:", err, values)
		return false
	}
	return true
}

func (rd *RedisRepository) HCount(key string) int64 {
	var r = rd.HLen(key)
	n, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " HCount error:", err)
		return 0
	}
	return n
}

func (rd *RedisRepository) HGets(key string) map[string]string {
	var r = rd.HGetAll(key)
	m, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " HGets error: ", err)
		return nil
	}
	return m
}

func (rd *RedisRepository) HKeyExists(key string, field string) bool {
	var r = rd.HExists(key, field)
	b, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " HKeyExists error:", err)
		return false
	}
	return b
}

func (rd *RedisRepository) HKeyGet(key string, field string) string {
	var r = rd.HGet(key, field)
	str, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " HKeyGet Error:", err)
		return ""
	}
	return str
}

func (rd *RedisRepository) Locked(key string) bool {
	var r = rd.Exists(key)
	b, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " Redis Locked Error:", err)
		return false
	}
	return b > 0
}

func (rd *RedisRepository) Lock(key string, d ...time.Duration) bool {
	if len(d) <= 0 || d[0] <= 0 {
		d = append(d, time.Minute)
	}
	var r = rd.SetNX(key, time.Now().Unix()+int64(d[0]), d[0])
	b, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println("Redis.DB", rd.GetOptionDb(), "Redis LOCK error : <"+key+"> ", d[0], err)
		return false
	}
	return b
}

func (rd *RedisRepository) LockWithValue(key string, v interface{}, d ...time.Duration) bool {
	if len(d) <= 0 || d[0] <= 0 {
		d = append(d, time.Minute)
	}
	var r = rd.SetNX(key, v, d[0])
	b, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println("Redis.DB", rd.GetOptionDb(), "Redis LOCK error : <"+key+"> ", d[0], err)
		return false
	}
	return b
}

func (rd *RedisRepository) Empty(key string) bool {
	var r = rd.Get(key)
	b, err := r.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), "Redis Check Empty Error:", err)
		return false
	}
	return b == ""
}

func (rd *RedisRepository) PushQueueUnique(key string, v ...string) int64 {
	var (
		uuid string
		arr  []interface{}
	)
	if len(v) <= 0 {
		return 0
	}
	uuid = fmt.Sprintf("%x", md5.Sum([]byte(key)))
	for _, it := range v {
		k := fmt.Sprintf("unique_id:%s:%s", uuid, it)
		if rd.Has(k) {
			continue
		}
		arr = append(arr, it)
		rd.SetNX(k, it, 12*time.Hour)
	}
	if len(arr) <= 0 {
		return 0
	}
	cmd := rd.LPush(key, arr...)
	n, err := cmd.Result()
	if err != nil && err != redis.Nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " PushQueue error:", err.Error(), v)
		return 0
	}
	return n
}

func (rd *RedisRepository) Has(key string) bool {
	var (
		res    = rd.Exists(key)
		n, err = res.Result()
	)
	if err != nil {
		rd.getLogger().Println(key, "Redis.DB", rd.GetOptionDb(), " Redis <Has> Error:", err.Error())
		return false
	}
	return n > 0
}

// GetOptions 获取配置
func (rd *RedisRepository) GetOptions() redis.Options {
	var (
		opts = rd.Client.Options()
	)
	if opts == nil {
		rd.getLogger().Println(" Redis GetOptions Nil")
		return redis.Options{}
	}
	return *opts
}

// GetOptionDb 获取当前实例选择的Db
func (rd *RedisRepository) GetOptionDb() int {
	var (
		opts = rd.GetOptions()
	)
	return opts.DB
}

// 获取日志对象
func (rd *RedisRepository) getLogger() *logrus.Logger {
	return GetLogger("redis")
}
