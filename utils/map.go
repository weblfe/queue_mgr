package utils

import (
		"fmt"
		"strconv"
)

func MGet(m map[string]interface{}, key string, def ...string) string {
	def = append(def, "")
	if v, ok := m[key]; ok {
		if str, ok := v.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", v)
	}
	return def[0]
}

func MGetMap(m map[string]interface{}, key string, def ...map[string]interface{}) map[string]interface{} {
	def = append(def, nil)
	if v, ok := m[key]; ok {
		if kv, ok := v.(map[string]interface{}); ok {
			return kv
		}
	}
	return def[0]
}

func MGetBool(m map[string]interface{}, key string, def ...bool) bool {
	def = append(def, false)
	if v, ok := m[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			if b, err := strconv.ParseBool(str); err == nil {
				return b
			}
			return def[0]
		case bool:
			return v.(bool)
		case int:
			n := v.(int)
			if n >= 1 {
				return true
			}
			return false
		}
	}
	return def[0]
}

func MGetInt(m map[string]interface{}, key string, def ...int) int {
	def = append(def, 0)
	if v, ok := m[key]; ok {
		switch v.(type) {
		case string:
			str := v.(string)
			if b, err := strconv.Atoi(str); err == nil {
				return b
			}
			return def[0]
		case int32:
			return int(v.(int32))
		case int64:
			return int(v.(int64))
		case int:
			return v.(int)
		case uint:
			return int(v.(uint))
		case uint32:
			return int(v.(uint32))
		case uint64:
			return int(v.(uint64))
		}
	}
	return def[0]
}

func MGetFloat(m map[string]interface{}, key string, def ...float64) float64 {
		def = append(def, 0)
		if v, ok := m[key]; ok {
				switch v.(type) {
				case string:
						str := v.(string)
						if b, err := strconv.ParseFloat(str,64); err == nil {
								return b
						}
						return def[0]
				case int32:
						return float64(v.(int32))
				case int64:
						return float64(v.(int64))
				case int:
						return float64(v.(int))
				case uint:
						return float64(v.(uint))
				case uint32:
						return float64(v.(uint32))
				case uint64:
						return float64(v.(uint64))
				case float64:
						return v.(float64)
				case float32:
						return float64(v.(float32))
				}
		}
		return def[0]
}

func MGetAny(m map[string]interface{}, key string, def ...interface{}) interface{} {
	def = append(def, nil)
	if v, ok := m[key]; ok {
		return v
	}
	return def[0]
}

func MGetStrArr(m map[string]interface{}, key string, def ...[]string) []string {
	def = append(def, nil)
	if v, ok := m[key]; ok {
		switch v.(type) {
		case []string:
			return v.([]string)
		case string:
			str := v.(string)
			return Str2Arr(str)
		case fmt.Stringer:
			str := v.(fmt.Stringer).String()
			return Str2Arr(str)
		case []interface{}:
			str := Slice2StrArr(v.([]interface{}))
			return str
		}
	}
	return def[0]
}
