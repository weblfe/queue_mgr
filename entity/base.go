package entity

import (
	"crypto/md5"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/weblfe/queue_mgr/utils"
	"reflect"
	"sort"
	"strconv"
	"time"
)

type (
	KvMap   map[string]interface{}
	KvArray []KvMap
	Arr     []string
)

func (kv KvMap) Len() int {
	return len(kv)
}

func (kv KvMap) String() string {
	return utils.JsonEncode(kv).String()
}

func (kv KvMap) Get(key string, def ...interface{}) interface{} {
	if len(def) <= 0 {
		def = append(def, nil)
	}
	if v, ok := kv[key]; ok {
		return v
	}
	return def[0]
}

func (kv KvMap) Exists(key string) bool {
	if _, ok := kv[key]; ok {
		return ok
	}
	return false
}

func (kv KvMap) GetStr(key string, def ...string) string {
	def = append(def, "")
	if v, ok := kv[key]; ok && v != "" {
		switch v.(type) {
		case string:
			return v.(string)
		case int:
			return fmt.Sprintf("%d", v.(int))
		case fmt.Stringer:
			return v.(fmt.Stringer).String()
		}
		return fmt.Sprintf("%v", v)
	}
	return def[0]
}

func (kv KvMap) GetInt(key string, def ...int) int {
	def = append(def, 0)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			if n, err := strconv.Atoi(str); err == nil {
				return n
			}
		case int:
			return v.(int)
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			if n, err := strconv.Atoi(str); err == nil {
				return n
			}
		}
	}
	return def[0]
}

func (kv KvMap) GetFloat(key string, def ...float64) float64 {
	def = append(def, 0)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			if n, err := strconv.ParseFloat(str, 64); err == nil {
				return n
			}
		case float64:
			return v.(float64)
		case float32:
			return float64(v.(float32))
		case int:
			return float64(v.(int))
		case int64:
			return float64(v.(int64))
		case int32:
			return float64(v.(int32))
		case uint:
			return float64(v.(uint))
		case uint64:
			return float64(v.(uint64))
		case uint32:
			return float64(v.(uint32))
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			if n, err := strconv.ParseFloat(str, 64); err == nil {
				return n
			}
		}
	}
	return def[0]
}

func (kv KvMap) GetBool(key string, def ...bool) bool {
	def = append(def, false)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			if n, err := strconv.ParseBool(str); err == nil {
				return n
			}
		case bool:
			return v.(bool)
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			if n, err := strconv.ParseBool(str); err == nil {
				return n
			}
		}
		str := fmt.Sprintf("%v", v)
		if n, err := strconv.ParseBool(str); err == nil {
			return n
		}
	}
	return def[0]
}

func (kv KvMap) GetDuration(key string, def ...time.Duration) time.Duration {
	def = append(def, 0)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			if n, err := time.ParseDuration(str); err == nil {
				return n
			}
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			if n, err := time.ParseDuration(str); err == nil {
				return n
			}
		case time.Duration:
			n := v.(time.Duration)
			return n
		}
	}
	return def[0]
}

func (kv KvMap) GetArr(key string, def ...[]string) []string {
	def = append(def, nil)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			return utils.Str2Arr(str)
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			return utils.Str2Arr(str)
		case []string:
			strArr := v.([]string)
			return strArr
		case []interface{}:
			return utils.Slice2StrArr(v.([]interface{}))
		}
	}
	return def[0]
}

func (kv KvMap) GetArrAny(key string, def ...[]interface{}) []interface{} {
	def = append(def, nil)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			var d = new([]interface{})
			if err := utils.JsonDecode([]byte(str), d); err != nil {
				return def[0]
			}
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			var d = new([]interface{})
			if err := utils.JsonDecode([]byte(str), d); err != nil {
				return def[0]
			}
		case []string:
			strArr := v.([]string)
			var d []interface{}
			for _, v := range strArr {
				d = append(d, v)
			}
			return d
		case []interface{}:
			return v.([]interface{})
		}
	}
	return def[0]
}

func (kv KvMap) GetKvMap(key string, def ...KvMap) KvMap {
	def = append(def, nil)
	if v, ok := kv[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			var d = new(map[string]interface{})
			if err := utils.JsonDecode([]byte(str), d); err != nil {
				return def[0]
			}
			return *d
		case fmt.Stringer:
			var d = new(map[string]interface{})
			str := v.(fmt.Stringer).String()
			if err := utils.JsonDecode([]byte(str), d); err != nil {
				return def[0]
			}
			return *d
		case map[string]string:
			strMap := v.(map[string]string)
			var d = make(map[string]interface{})
			for k, v := range strMap {
				d[k] = v
			}
			return d
		case map[string]interface{}:
			return v.(map[string]interface{})
		case []interface{}:
			arr := v.([]interface{})
			var d = make(map[string]interface{})
			for k, v := range arr {
				d[fmt.Sprintf("%d", k)] = v
			}
			return d
		}
	}
	return def[0]
}

func (kv KvMap) Add(key string, v interface{}) KvMap {
	kv[key] = v
	return kv
}

func (kv KvMap) Keys() []string {
	var keys []string
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (kv KvMap) HashCode() string {
	var hash = md5.New()
	hash.Write([]byte(kv.String()))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (kv KvMap) Foreach(iter func(k string, v interface{}) bool) {
	for k, v := range kv {
		if !iter(k, v) {
			return
		}
	}
}

func (kv KvMap) ConvertStrKv() map[string]string {
	var result = make(map[string]string)
	for k, v := range kv {
		result[k] = utils.NewStringer(v).Convert()
	}
	return result
}

func (kv KvMap) Copy() KvMap {
	var newKv = KvMap{}
	for k, v := range kv {
		newKv[k] = v
	}
	return newKv
}

func (kv KvMap) Merge(m KvMap) KvMap {
	if m == nil {
		return kv
	}
	if kv == nil {
		return m
	}
	for k, v := range m {
		kv[k] = v
	}
	return kv
}

func (kv KvMap)Filter(arr []string) KvMap {
	var kvObj = KvMap{}
	for _,k:= range arr {
		if v,ok:=kv[k];ok {
			kvObj[k] = v
		}
	}
	return kvObj
}

func NewKvArr() *KvArray {
	var arr = new(KvArray)
	return arr
}

func (arr *KvArray) Transform(items interface{}) KvArray {
	if items == nil {
		return *arr
	}
	var (
		v    = reflect.ValueOf(items)
		kind = v.Kind()
	)
	if kind == reflect.Interface || kind == reflect.Ptr {
		v = v.Elem()
		kind = v.Kind()
	}
	switch kind {
	case reflect.Array:
		var size = v.Len()
		for i := 0; i < size; i++ {
			var it = arr.parse(v.Index(i).Interface())
			if it == nil {
				continue
			}
			*arr = append(*arr, it)
		}
	case reflect.Slice:
		var size = v.Len()
		for i := 0; i < size; i++ {
			var it = arr.parse(v.Index(i).Interface())
			if it == nil {
				continue
			}
			*arr = append(*arr, it)
		}
	default:
		var it = arr.parse(v.Interface())
		if it != nil {
			*arr = append(*arr, it)
		}
	}
	return *arr
}

func (arr *KvArray) parse(v interface{}) KvMap {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case map[string]interface{}:
		return v.(map[string]interface{})
	case *map[string]interface{}:
		return *v.(*map[string]interface{})
	case KvMap:
		return v.(KvMap)
	case *KvMap:
		return *v.(*KvMap)
	}
	var (
		value = reflect.ValueOf(v)
		kind  = value.Kind()
	)
	switch kind {
	case reflect.Struct, reflect.Map:
		var (
			encoder = utils.JsonEncode(v)
			err     = encoder.Error()
			bytes   = encoder.Bytes()
		)
		if err != nil {
			return nil
		}
		var kv = &KvMap{}
		if err = utils.JsonDecode(bytes, kv); err != nil {
			return nil
		}
		return *kv
	default:
		return KvMap{
			"data": v,
		}
	}
}

func (arr *KvArray) Add(kv KvMap) KvArray {
	if arr == nil {
		return KvArray{kv}
	}
	*arr = append(*arr, kv)
	return *arr
}

func (arr Arr)Include(v string) bool {
	for _,value:=range arr {
		if value == v {
			return true
		}
	}
	return false
}

func (arr Arr)Index(v string) int {
	for i,value:=range arr {
		if value == v {
			return i
		}
	}
	return -1
}

func (arr *Arr)Len() int {
	return len(*arr)
}

func (arr *Arr) Sort()  {
	sort.Strings(*arr)
}


// argsToJsonBytes 读取query 参数转 json bytes
func argsToJsonBytes(args *fiber.Args) []byte {
	var data = make(map[string]interface{})
	args.VisitAll(func(key, v []byte) {
		data[string(key)] = string(v)
	})
	if len(data) > 0 {
		return utils.JsonEncode(data).Bytes()
	}
	return []byte(`{}`)
}
