package config

import (
	"fmt"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/utils"
	"strings"
)

// Database 数据库配置
type Database struct {
	DbName     string `json:"name"`
	DbDriver   string `json:"driver"`
	DbUser     string `json:"user"`
	DbPassword string `json:"password"`
	DbPort     int    `json:"port"`
	DbHost     string `json:"host"`
	DbPrefix   string `json:"prefix"`
	DbOptions  string `json:"options"`
	DbDefault  bool   `json:"default,omitempty"`
}

type DatabaseKv map[string]Database

func CreateDatabase(v ...interface{}) Database {
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
			return Database{
				DbName:     utils.MGet(m, "name", "mysql"),
				DbDriver:   utils.MGet(m, "driver", "mysql"),
				DbUser:     utils.MGet(m, "user", "root"),
				DbPassword: utils.MGet(m, "password", "root"),
				DbPort:     utils.MGetInt(m, "port", 3306),
				DbHost:     utils.MGet(m, "host", "127.0.0.1"),
				DbPrefix:   utils.MGet(m, "prefix", ""),
				DbDefault:  utils.MGetBool(m, "default", false),
				DbOptions:  utils.MGet(m, "options", "charset=utf8mb4&parseTime=True&loc=Local"),
			}
		}
	}
	var env = NewEnvKvGetter(key, "_")
	return Database{
		DbName:     env.Get("DB_NAME", "mysql"),
		DbDriver:   env.Get("DB_DRIVER", "mysql"),
		DbUser:     env.Get("DB_USER", "root"),
		DbPassword: env.Get("DB_PASSWORD", "root"),
		DbPort:     env.GetInt("DB_PORT", 3306),
		DbHost:     env.Get("DB_HOST", "127.0.0.1"),
		DbPrefix:   env.Get("DB_PREFIX", ""),
		DbDefault:  env.GetBool("DB_DEFAULT", false),
		DbOptions:  env.Get("DB_OPTIONS", "charset=utf8mb4&parseTime=True&loc=Local"),
	}
}

// GetConnUrl conn 获取连接
func (db *Database) GetConnUrl() string {
	var url = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		db.DbUser, db.DbPassword, db.DbHost,
		db.DbPort, db.DbName,
	)
	if db.DbOptions != "" {
		if strings.Contains(url, "?") {
			return url + db.DbOptions
		} else {
			return url + "?" + db.DbOptions
		}
	}
	return url
}

func (db *Database) Get(key string, def ...interface{}) interface{} {
	switch strings.ToLower(key) {
	case "name":
		return db.DbName
	case "port":
		return db.DbPort
	case "host":
		return db.DbHost
	case "user":
		return db.DbUser
	case "password":
		return db.DbPassword
	case "options":
		return db.DbOptions
	case "prefix":
		return db.DbPrefix
	case "default":
		return db.DbDefault
	}
	return def[0]
}

func (data *DatabaseKv) Get(key string) (Database, bool) {
	if v, ok := (*data)[key]; ok {
		return v, true
	}
	return CreateDatabase(key), false
}

func (data *DatabaseKv) GetDefaultKey() string {
	for k, v := range *data {
		if v.DbDefault {
			return k
		}
	}
	return "default"
}

func (data *DatabaseKv) Exists(key string) bool {
	if _, ok := (*data)[key]; ok {
		return true
	}
	return false
}

func (data *DatabaseKv) Size() int {
	return len(*data)
}

func (data *DatabaseKv) Keys() []string {
	var keys []string
	for k := range *data {
		keys = append(keys, k)
	}
	return keys
}

func (data *DatabaseKv) Add(key string, database Database) *DatabaseKv {
	if _, ok := (*data)[key]; ok {
		return data
	}
	(*data)[key] = database
	return data
}

func (data *DatabaseKv) String() string {
	return utils.JsonEncode(data).String()
}

func (data *DatabaseKv) Decode(content []byte) error {
	return utils.JsonDecode(content, data)
}

func (data *DatabaseKv) ValueOf(key string, def ...interface{}) interface{} {
	if !strings.Contains(key, ".") {
		if d, ok := data.Get(key); ok {
			return d
		}
		return def[0]
	}
	strArr := strings.Split(key, ".")
	d, ok := data.Get(strArr[0])
	if !ok {
		return def[0]
	}
	return d.Get(strings.Join(strArr[1:], ""), def...)
}

func (data *DatabaseKv) MAdd(mArr map[string]interface{}) int {
	for k, v := range mArr {
		data.Add(k, CreateDatabase(v))
	}
	return 0
}

// 注册
func registerDbFactory(app *applicationConfiguration) {
	app.Register("database", createDatabaseKv)
}

// 创建database kv
func createDatabaseKv(v interface{}) facede.CfgKv {
	switch v.(type) {
	case map[string]interface{}:
		dataArr := DatabaseKv{}
		dataArr.MAdd(v.(map[string]interface{}))
		return &dataArr
	}
	return nil
}
