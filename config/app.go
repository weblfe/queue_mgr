package config

import (
	"fmt"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/utils"
	"strings"
)

type App struct {
	Port    int    `json:"port,default=8080"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

func CreateApp(v ...interface{}) *App {
	if len(v) >= 0 && v[0] != nil {
		switch v[0].(type) {
		case map[string]interface{}:
			m := v[0].(map[string]interface{})
			return &App{
				Port:    utils.MGetInt(m, "port", 0),
				Name:    utils.MGet(m, "name", "robot"),
				Version: utils.MGet(m, "version", "0.1.0"),
			}
		}
	}
	return &App{
		Port:    utils.GetEnvInt("APP_PORT", 80),
		Name:    utils.GetEnvVal("APP_NAME", "robot"),
		Version: utils.GetEnvVal("APP_VERSION", "0.1.0"),
	}
}

func (app *App) String() string {
	return utils.JsonEncode(app).String()
}

func (app *App) Decode(content []byte) error {
	return utils.JsonDecode(content, app)
}

func (app *App) ValueOf(key string, def ...interface{}) interface{} {
	if len(def) == 0 {
		def = append(def, nil)
	}
	switch strings.ToLower(key) {
	case "name":
		return app.Name
	case "port":
		return app.Port
	}
	return def[0]
}

func (app *App) Keys() []string {
	return []string{
		"name", "port", "version",
	}
}

func (app *App) MAdd(mArr map[string]interface{}) int {
	app.Port = utils.MGetInt(mArr, "port", app.Port)
	app.Name = utils.MGet(mArr, "name", app.Name)
	app.Version = utils.MGet(mArr, "version", app.Version)
	return 0
}

func (app *App) GetAddr() string {
	if app.Port <=0 {
			app.Port = utils.GetEnvInt("APP_PORT",8081)
	}
	return fmt.Sprintf(":%d", app.Port)
}

// 注册
func registerAppCfgFactory(app *applicationConfiguration) {
	app.Register("app", func(v interface{}) facede.CfgKv {
		return CreateApp(v)
	})
}
