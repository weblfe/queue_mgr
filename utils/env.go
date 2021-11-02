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
	matchRegexp = regexp.MustCompile(`^\$\{.+\}$`)
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

// GetEnvOk 是否获取成功
func GetEnvOk(key string) (string, bool) {
	var v = os.Getenv(key)
	if v == "" {
		return "", false
	}
	return v, true
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
	if depth[0] > 3 {
		return value
	}
	if value != "" && strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		var splits = strings.Split(value, "${")
		for _, v := range splits {
			if v == "" {
				continue
			}
			if count := strings.Count(v, "}"); count == 0 || count > 1 {
				continue
			}
			var key = strings.Split(v, "}")
			if len(key) < 1 {
				continue
			}
			var valueOr, ok = GetEnvOk(key[0])
			if !ok {
				continue
			}
			if valueOr != "" && strings.HasPrefix(valueOr, "${") && strings.HasSuffix(valueOr, "}") {
				valueOr = ParseEnvValue(valueOr, depth[0]+1)
			}
			if valueOr != "" {
				value = strings.ReplaceAll(value, fmt.Sprintf("${%s}", key[0]), valueOr)
			}
		}
	}
	return value
}
