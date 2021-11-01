package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	envVarReg = regexp.MustCompile(`(\$\{.+\})`)
)

// GetEnvVal 获取字符串
func GetEnvVal(key string, def ...string) string {
	return GetEnvOr(key, def...)
}

// GetEnvOr 获取环境变量
func GetEnvOr(key string, def ...string) string {
	def = append(def, "")
	var v = os.Getenv(key)
	if v == "" {
		return def[0]
	}
	return v
}

// GetEnvBool 获取bool
func GetEnvBool(key string, def ...bool) bool {
	def = append(def, false)
	var b = GetEnvVal(key)
	if b == "" {
		return def[0]
	}
	bVal, err := strconv.ParseBool(b)
	if err != nil {
		fmt.Println("GetEnvBool Error:", err)
		return def[0]
	}
	return bVal
}

// GetEnvInt 获取整形
func GetEnvInt(key string, def ...int) int {
	def = append(def, 0)
	var b = GetEnvVal(key)
	if b == "" {
		return def[0]
	}
	intVal, err := strconv.Atoi(b)
	if err != nil {
		fmt.Println("GetEnvInt Error:", err)
		return def[0]
	}
	return intVal
}

// GetEnvFloat 获取浮点数
func GetEnvFloat(key string, def ...float64) float64 {
	def = append(def, 0)
	var b = GetEnvVal(key)
	if b == "" {
		return def[0]
	}
	floatVal, err := strconv.ParseFloat(b, 64)
	if err != nil {
		fmt.Println("GetEnvFloat Error:", err)
		return def[0]
	}
	return floatVal
}

// GetEnvDuration 获取时长
func GetEnvDuration(key string, def ...time.Duration) time.Duration {
	def = append(def, 0)
	var b = GetEnvVal(key)
	if b == "" {
		return def[0]
	}
	d, err := time.ParseDuration(b)
	if err != nil {
		fmt.Println("GetEnvDuration Error:", err)
		return def[0]
	}
	return d
}

// GetEnvTime 获取时间
func GetEnvTime(key string, def ...*time.Time) *time.Time {
	def = append(def, nil)
	var (
		err error
		d   time.Time
		b   = GetEnvVal(key)
	)
	if b == "" {
		return def[0]
	}
	if strings.Contains(b, ":") {
		d, err = time.Parse(`2006-01-02 15:04:05`, b)
		if err != nil {
			fmt.Println("GetEnvTime Error:", err)
			return def[0]
		}
	} else {
		d, err = time.Parse(`2006-01-02`, b)
		if err != nil {
			fmt.Println("GetEnvTime Error:", err)
			return def[0]
		}
	}
	return &d
}

// GetEnvMapper 获取hashMap
func GetEnvMapper(key string, def ...map[string]interface{}) map[string]interface{} {
	def = append(def, nil)
	var b = GetEnvVal(key)
	if b == "" {
		return def[0]
	}
	if !json.Valid([]byte(b)) {
		return def[0]
	}
	var mapper = make(map[string]interface{})
	if err := JsonDecode([]byte(b), &mapper); err != nil {
		fmt.Println("GetEnvMapper Error:", err)
		return def[0]
	}
	return mapper
}

// GetEnvArr 获取字符串数据
func GetEnvArr(key string, def ...[]string) []string {
	def = append(def, nil)
	var b = GetEnvVal(key)
	if b == "" {
		return def[0]
	}
	return Str2Arr(b)
}

// ParseEnvValue 解析
func ParseEnvValue(value string, depth ...int) string {
	depth = append(depth, 0)
	if value != "" && envVarReg.Match([]byte(value)) {
		var subArr = envVarReg.FindAllSubmatch([]byte(value), -1)
		if len(subArr) <= 0 {
			return value
		}
		for _, v := range subArr {
			var (
				name = strings.TrimPrefix(strings.TrimSuffix(string(v[0]), `}`), `${`)
				valueOr  = GetEnvOr(name)
			)
			if valueOr == "" {
				continue
			}
			if envVarReg.Match([]byte(valueOr)) {
				if depth[0] <= 3 {
					valueOr = ParseEnvValue(valueOr, depth[0]+1)
				}
			}
			value = strings.ReplaceAll(value, string(v[0]), valueOr)
		}
	}
	return value
}
