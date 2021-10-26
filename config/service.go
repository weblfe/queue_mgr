package config

import (
	"fmt"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/utils"
	"os"
	"strconv"
	"strings"
)

// Service 服务配置对象
type Service struct {
	Name        string   `json:"name,omitempty"`
	Methods     []string `json:"methods,omitempty"`
	EntryPoints []string `json:"entrypoints,omitempty"`
	Widgets     []int    `json:"weights,omitempty"`
}

type ServiceKv map[string]*Service

const (
	defaultEntrypoint = "127.0.0.1"
)

func (s *ServiceKv) Keys() []string {
	var keys []string
	for k := range *s {
		keys = append(keys, k)
	}
	return keys
}

func (s *ServiceKv) MAdd(m map[string]interface{}) int {
	for k, v := range m {
		s.Add(k, v)
	}
	return 0
}

func (s *ServiceKv) Decode(content []byte) error {
	return utils.JsonDecode(content, s)
}

func (s *ServiceKv) Add(k string, v interface{}) {
	if _, ok := (*s)[k]; ok {
		return
	}
	switch v.(type) {
	case map[string]interface{}:
		(*s)[k] = CreateService(v)
		return
	}
	(*s)[k] = CreateService(k)
}

func (s *ServiceKv) ValueOf(s2 string, def ...interface{}) interface{} {
	if val, ok := (*s)[s2]; ok {
		return val
	}
	return def[0]
}

func (s *ServiceKv) Get(key string) (*Service, bool) {
	if val, ok := (*s)[key]; ok {
		return val, ok
	}
	if val, ok := (*s)[strings.ToLower(key)]; ok {
		return val, ok
	}
	return nil, false
}

func (s *ServiceKv) String() string {
	return utils.JsonEncode(s).String()
}

func CreateService(v ...interface{}) *Service {
	var key = ""
	if len(v) >= 0 && v[0] != nil {
		switch v[0].(type) {
		case map[string]interface{}:
			m := v[0].(map[string]interface{})
			arr := utils.MGetStrArr(m, "entrypoints")
			widgets := parserCreateWidgets(arr)
			return &Service{
				Name:        utils.MGet(m, "name"),
				Methods:     utils.MGetStrArr(m, "methods"),
				EntryPoints: arr,
				Widgets:     widgets,
			}
		}
	}
	var env = NewEnvKvGetter(key, "_")
	arr := env.GetArr("SERVICE_ENTRYPOINTS")
	widgets := parserCreateWidgets(arr)
	return &Service{
		Name:        env.Get("SERVICE_NAME"),
		Methods:     env.GetArr("SERVICE_METHODS"),
		EntryPoints: arr,
		Widgets:     widgets,
	}
}

// 解析创建权重
func parserCreateWidgets(arr []string) []int {
	var widgets []int
	for i, v := range arr {
		var (
			widget    = 100
			hasWidget = strings.Contains(v, "|")
		)
		if !hasWidget {
			arr[i] = strings.TrimSpace(v)
		} else {
			arrV := strings.SplitN(v, "|", 2)
			arr[i] = strings.TrimSpace(arrV[0])
			if n, err := strconv.Atoi(strings.TrimSpace(arrV[1])); err == nil && n >= 0 {
				widget = n
			}
		}
		widgets = append(widgets, widget)
	}
	return widgets
}

func (serv *Service) String() string {
	return utils.JsonEncode(serv).String()
}

func (serv *Service) GetEntryPoint() string {
	// var index = serv.Widgets
	if serv.EntryPoints == nil || len(serv.EntryPoints) <= 0 {
		var url = ""
		if serv.Name != "" {
			var key = strings.ToUpper(fmt.Sprintf("SERVER_%s_ENTRYPOINT", serv.Name))
			url = os.Getenv(key)
		}
		if url == "" {
			url = os.Getenv("SERVER_ENTRYPOINT")
		}
		if url == "" {
			return defaultEntrypoint
		}
		return url
	}
	return serv.EntryPoints[0]
}

// 创建redis 配置
func createServiceKv(v interface{}) *ServiceKv {
	switch v.(type) {
	case map[string]interface{}:
		kv := ServiceKv{}
		kv.MAdd(v.(map[string]interface{}))
		return &kv
	}
	return nil
}

// 注册
func registerServiceKvFactory(app *applicationConfiguration) {
	app.Register("services", func(v interface{}) facede.CfgKv {
		return createServiceKv(v)
	})
}
