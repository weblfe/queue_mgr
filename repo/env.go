package repo

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	_varL       = "${"
	_varR       = "}"
	_timeLayout = "2006-01-02 15:04:05"
)

var (
	defaultEnvRepo = newEnvironmentRepository()
)

type (
	environmentRepository struct {
		lock  sync.RWMutex
		cache map[string]*envMethod
	}

	envMethod struct {
		prefix    string
		domainVar bool
		_varL     string
		_varR     string
	}
)

func newEnvironmentRepository() *environmentRepository {
	env := environmentRepository{}
	env.lock = sync.RWMutex{}
	env.cache = make(map[string]*envMethod)
	return &env
}

func (e *envMethod) Get(key string) string {
	if e.prefix != "" && !strings.HasPrefix(key, e.prefix) {
		key = e.prefix + key
	}
	vt := os.Getenv(key)
	if e.domainVar && vt != "" {
		vs := e.parseVal(vt)
		// 解析后更新,提供访问效率
		if vs != vt {
			e.Put(key, vs)
		}
		return vs
	}
	return vt
}

func (e *envMethod) Put(key string, v interface{}) bool {
	var str string
	switch v.(type) {
	case string:
		str = v.(string)
	case *string:
		str = *v.(*string)
	case fmt.Stringer:
		str = v.(fmt.Stringer).String()
	default:
		vs, err := json.Marshal(v)
		if err != nil {
			return false
		}
		str = string(vs)
	}
	if e.prefix != "" && !strings.HasPrefix(key, e.prefix) {
		key = e.prefix + key
	}
	if err := os.Setenv(key, str); err != nil {
		return false
	}
	return true
}

func (e *envMethod) SetPrefix(prefix string) *envMethod {
	e.prefix = prefix
	return e
}

func (e *envMethod) SetDomainVar(on bool) *envMethod {
	e.domainVar = on
	return e
}

func (e *envMethod) SetVarOption(l, r string) *envMethod {
	e._varL = l
	e._varR = r
	return e
}

func (e *envMethod) parseVal(v string) string {
	var l, r = e.getOption()
	if strings.Contains(v, l) && strings.Contains(v, r) {
		var (
			arr = strings.Split(v, l)
		)
		for _, vs := range arr {
			k := ""
			vt := ""
			if strings.Contains(vs, r) {
				k = l + vs
				vt = os.Getenv(strings.ReplaceAll(vs, r, ""))
			}
			// 是否可变变量
			vt = e.parseVal(vt)
			// 替换
			if k != "" {
				v = strings.ReplaceAll(v, k, vt)
			}
		}
	}
	return v
}

func (e *envMethod) getOption() (string, string) {
	if e._varL == "" {
		e._varL = _varL
	}
	if e._varR == "" {
		e._varR = _varR
	}
	return e._varL, e._varR
}

func (e *envMethod) GetInt(key string) int {
	var v = e.Get(key)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}

func (e *envMethod) GetBool(key string) bool {
	var v = e.Get(key)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

func (e *envMethod) GetFloat(key string) float32 {
	var v = e.Get(key)
	if v == "" {
		return 0
	}
	b, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0
	}
	return float32(b)
}

func (e *envMethod) GetDuration(key string) time.Duration {
	var v = e.Get(key)
	if v == "" {
		return 0
	}
	b, err := time.ParseDuration(v)
	if err != nil {
		return 0
	}
	return b
}

func (e *envMethod) GetDateTime(key string) *time.Time {
	var v = e.Get(key)
	if v == "" {
		return nil
	}
	t, err := time.Parse(_timeLayout, v)
	if err != nil {
		return nil
	}
	return &t
}

func (e *envMethod) GetArr(key string) []string {
	var v = e.Get(key)
	if v == "" {
		return nil
	}
	// json array
	if strings.Contains(v, "[") && strings.Contains(v, "]") {
		var strArr []string
		err := json.Unmarshal([]byte(v), &strArr)
		if err != nil {
			return nil
		}
		return strArr
	}
	if strings.Contains(v, ",") && !strings.Contains(v, "[") {
		return strings.Split(v, ",")
	}
	return nil
}

func (e *envMethod) GetIntArr(key string) []int {
	var v = e.Get(key)
	if v == "" {
		return nil
	}
	var intArr []int
	// json array
	if strings.Contains(v, "[") && strings.Contains(v, "]") {

		err := json.Unmarshal([]byte(v), &intArr)
		if err != nil {
			return nil
		}
		return intArr
	}
	// 逗号分割数组
	if strings.Contains(v, ",") && !strings.Contains(v, "[") {
		arr := strings.Split(v, ",")
		for _, v := range arr {
			n, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			intArr = append(intArr, n)
		}
		if len(intArr) != 0 {
			return intArr
		}
	}
	return nil
}

func (e *envMethod) GetMap(key string) map[string]interface{} {
	var v = e.Get(key)
	if v == "" {
		return nil
	}
	var m = make(map[string]interface{})
	// json array
	err := json.Unmarshal([]byte(v), &m)
	if err != nil {
		log.Println(map[string]interface{}{"ENV_GET_MAP_ERROR": err})
	}
	return m
}

func GetEnvRepository() *environmentRepository {
	return defaultEnvRepo
}

func (repo *environmentRepository) Env(prefix ...string) *envMethod {
	if len(prefix) <= 0 {
		prefix = append(prefix, "default")
	}
	repo.lock.Lock()
	defer repo.lock.Unlock()
	if v, ok := repo.cache[prefix[0]]; ok {
		return v
	}
	if prefix[0] == "default" {
		repo.cache[prefix[0]] = newPrefixEnvUtil(prefix...)
	}
	return repo.cache[prefix[0]]
}

func newPrefixEnvUtil(prefix ...string) *envMethod {
	var env = new(envMethod)
	if len(prefix) == 0 || prefix[0] == "" || prefix[0] == "default" {
		prefix[0] = ""
	}
	env.prefix = prefix[0]
	return env
}

func (e *envMethod) GetOf(key string, def ...string) string {
	var s = e.Get(key)
	if s == "" && len(def) != 0 && def[0] != "" {
		return def[0]
	}
	return s
}

func (e *envMethod) GetIntOf(key string, def ...int) int {
	var v = e.GetInt(key)
	if v == 0 && len(def) != 0 && def[0] != 0 {
		return def[0]
	}
	return v
}

func (e *envMethod) GetFloatOf(key string, def ...float32) float32 {
	var v = e.GetFloat(key)
	if v == 0 && len(def) != 0 && def[0] != 0 {
		return def[0]
	}
	return v
}

func (e *envMethod) GetArrOf(key string, def ...[]string) []string {
	var v = e.GetArr(key)
	if v == nil && len(def) != 0 && def[0] != nil {
		return def[0]
	}
	return v
}
