package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type (
	CaseMode int

	envTagDecoder struct {
		caseMode CaseMode
		prefix   string
		suffix   string
		tag      string
	}

	tagToken struct {
		Key     string
		Default string
		value   *reflect.Value
	}
)

const (
	envTag                  = "env"
	UnDefineCase   CaseMode = 0
	UpperCase      CaseMode = 1
	LowerCase      CaseMode = 2
	NormalCase     CaseMode = 3
	DateTimeLayout          = `2006-01-02 15:04:05`
)

func NewEnvDecoder(caseMode ...CaseMode) *envTagDecoder {
	caseMode = append(caseMode, UpperCase)
	var decoder = new(envTagDecoder)
	decoder.caseMode = caseMode[0]
	decoder.prefix = ""
	decoder.suffix = ""
	decoder.tag = envTag
	return decoder
}

func (decoder *envTagDecoder) SetPrefix(prefix string) *envTagDecoder {
	decoder.prefix = prefix
	return decoder
}

func (decoder *envTagDecoder) SetSuffix(suffix string) *envTagDecoder {
	decoder.suffix = suffix
	return decoder
}

func (decoder *envTagDecoder) GetPrefix() string {
	return decoder.prefix
}

func (decoder *envTagDecoder) GetSuffix() string {
	return decoder.suffix
}

func (decoder *envTagDecoder) Marshal(v interface{}) error {
	var tokens = decoder.parse(v)
	if len(tokens) <= 0 {
		return errors.New("tag parse failed")
	}
	return decoder.load(tokens)
}

func (decoder *envTagDecoder) parse(v interface{}, valueAddr ...reflect.Value) []*tagToken {
	if v == nil {
		return nil
	}
	var (
		tokens   []*tagToken
		valuePtr = reflect.ValueOf(v)
		kind     = valuePtr.Kind()
		canAddr  = valuePtr.CanAddr()
	)
	if kind == reflect.Ptr && !canAddr && valuePtr.Elem().Kind() == reflect.Struct {
		var (
			valueOf = valuePtr.Elem()
			value   = valueOf.Type()
			num     = value.NumField()
		)
		for i := 0; num > i; i++ {
			var values = decoder.values(value.Field(i), valueOf.Field(i))
			if len(values) <= 0 {
				continue
			}
			tokens = append(tokens, values...)
		}
		return tokens
	}
	if kind == reflect.Struct && !canAddr {
		var (
			value = valuePtr.Type()
			num   = valuePtr.NumField()
		)
		if len(valueAddr) > 0 {
			valuePtr = valueAddr[0]
		}
		for i := 0; num > i; i++ {
			var values = decoder.values(value.Field(i), valuePtr.Field(i))
			if len(values) <= 0 {
				continue
			}
			tokens = append(tokens, values...)
		}
	}
	return tokens
}

func (decoder *envTagDecoder) load(tokens []*tagToken) error {
	var err error
	defer func() {
		if v := recover(); v != nil {
			err = v.(error)
		}
	}()
	for _, v := range tokens {
		if v == nil || v.Key == "" || v.value == nil {
			continue
		}
		decoder.set(v.value, decoder.GetEnvOr(v.Key, v.Default))
	}
	return err
}

func (decoder *envTagDecoder) set(value *reflect.Value, data string) {
	if value == nil {
		return
	}
	if !value.CanSet() {
		return
	}
	var valueAny = value.Interface()
	// 时间类型处理
	switch valueAny.(type) {
	case time.Time:
		if d, err := time.Parse(time.RFC3339, data); err == nil {
			value.Set(reflect.ValueOf(d))
			return
		}
		if d, err := time.Parse(DateTimeLayout, data); err == nil {
			value.Set(reflect.ValueOf(d))
			return
		}
	case time.Duration:
		if d, err := time.ParseDuration(data); err == nil {
			value.Set(reflect.ValueOf(d))
			return
		}
	}
	// 基础类型映射解码
	var kind = value.Kind()
	switch kind {
	case reflect.Bool:
		if b, err := strconv.ParseBool(data); err == nil {
			value.Set(reflect.ValueOf(b))
		}
	case reflect.Int:
		if n, err := strconv.Atoi(data); err == nil {
			value.Set(reflect.ValueOf(n))
		}
	case reflect.Int8:
		if n, err := strconv.Atoi(data); err == nil {
			value.Set(reflect.ValueOf(int8(n)))
		}
	case reflect.Int16:
		if n, err := strconv.Atoi(data); err == nil {
			value.Set(reflect.ValueOf(int16(n)))
		}
	case reflect.Int32:
		if n, err := strconv.Atoi(data); err == nil {
			value.Set(reflect.ValueOf(int32(n)))
		}
	case reflect.Int64:
		if n, err := strconv.Atoi(data); err == nil {
			value.Set(reflect.ValueOf(int64(n)))
		}
	case reflect.Uint:
		if n, err := strconv.ParseUint(data, 10, 64); err == nil {
			value.Set(reflect.ValueOf(uint(n)))
		}
	case reflect.Uint8:
		if n, err := strconv.ParseUint(data, 10, 64); err == nil {
			value.Set(reflect.ValueOf(uint8(n)))
		}
	case reflect.Uint16:
		if n, err := strconv.ParseUint(data, 10, 64); err == nil {
			value.Set(reflect.ValueOf(uint16(n)))
		}
	case reflect.Uint32:
		if n, err := strconv.ParseUint(data, 10, 64); err == nil {
			value.Set(reflect.ValueOf(uint32(n)))
		}
	case reflect.Uint64:
		if n, err := strconv.ParseUint(data, 10, 64); err == nil {
			value.Set(reflect.ValueOf(n))
		}
	case reflect.Float32:
		if n, err := strconv.ParseFloat(data, 64); err == nil {
			value.Set(reflect.ValueOf(float32(n)))
		}
	case reflect.Float64:
		if n, err := strconv.ParseFloat(data, 64); err == nil {
			value.Set(reflect.ValueOf(n))
		}
	case reflect.Complex64:
		if n, err := strconv.ParseComplex(data, 64); err == nil {
			value.Set(reflect.ValueOf(complex64(n)))
		}
	case reflect.Complex128:
		if n, err := strconv.ParseComplex(data, 128); err == nil {
			value.Set(reflect.ValueOf(n))
		}
	case reflect.Map:
		var bytes = []byte(data)
		if err := decoder.bytesJsonDecoder(bytes, value.Addr().Interface(), true); err != nil {
			return
		}
	case reflect.Array:
		var bytes, ok = decoder.arrBytesDecoder(data)
		if !ok {
			return
		}
		if err := decoder.bytesJsonDecoder(bytes, value.Addr().Interface()); err != nil {
			return
		}
	case reflect.Slice:
		var bytes, ok = decoder.arrBytesDecoder(data)
		if !ok {
			return
		}
		if err := decoder.bytesJsonDecoder(bytes, value.Addr().Interface()); err != nil {
			return
		}
	case reflect.String:
		value.Set(reflect.ValueOf(data))
	case reflect.Struct:
		var bytes = []byte(data)
		if err := decoder.bytesJsonDecoder(bytes, value.Addr().Interface(), true); err != nil {
			return
		}
	}
	return
}

func (decoder *envTagDecoder) bytesJsonDecoder(bytes []byte, addr interface{}, check ...bool) error {
	check = append(check, false)
	if check[0] && !json.Valid(bytes) {
		return errors.New("bytes json valid failed")
	}
	return json.Unmarshal(bytes, addr)
}

func (decoder *envTagDecoder) arrBytesDecoder(data string) ([]byte, bool) {
	var bytes = []byte(data)
	if !json.Valid(bytes) {
		if strings.Contains(data, "{") || strings.Contains(data, "}") {
			return nil, false
		}
		if strings.Contains(data, "[") || strings.Contains(data, "]") {
			return nil, false
		}
		data = fmt.Sprintf("[%s]", strings.TrimSpace(data))
		bytes = []byte(data)
		if !json.Valid(bytes) {
			return nil, false
		}
	}
	return bytes, true
}

func (decoder *envTagDecoder) GetEnvOr(key string, def ...string) string {
	var k = decoder.make(key)
	if v := os.Getenv(k); v != "" {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

func (decoder *envTagDecoder) DebugKey(key string) string {
	return decoder.make(key)
}

func (decoder *envTagDecoder) make(key string) string {
	if key == "" {
		return ""
	}
	if strings.HasSuffix(key, "_") {
		key = strings.TrimSuffix(key, "_")
	}
	if strings.HasPrefix(key, "_") {
		key = strings.TrimPrefix(key, "_")
	}
	if decoder.prefix != "" {
		if strings.HasPrefix(decoder.suffix, "_") {
			key = decoder.prefix + key
		} else {
			key = fmt.Sprintf("%s_%s", decoder.prefix, key)
		}
	}
	if decoder.suffix != "" {
		if strings.HasPrefix(decoder.suffix, "_") {
			key = key + decoder.suffix
		} else {
			key = fmt.Sprintf("%s_%s", key, decoder.suffix)
		}
	}
	key = strings.TrimSpace(key)
	switch decoder.caseMode {
	case UpperCase:
		return strings.ToUpper(key)
	case LowerCase:
		return strings.ToLower(key)
	case NormalCase:
		return key
	}
	return key
}

func (decoder *envTagDecoder) values(field reflect.StructField, value reflect.Value) []*tagToken {
	var (
		key, defaultValue string
		tokenArr          []*tagToken
		token             = new(tagToken)
		kind              = field.Type.Kind()
		anonymous         = field.Anonymous
	)
	if !unicode.IsUpper([]rune(field.Name)[0]) {
		return tokenArr
	}
	// 匿名解构体
	if kind == reflect.Struct && value.CanAddr() && !anonymous {
		tokenArr = decoder.parse(value.Interface(), value)
		if len(tokenArr) > 0 {
			return tokenArr
		}
	}
	// 解析
	if kind == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct && !anonymous {
		tokenArr = decoder.parse(value.Interface())
		if len(tokenArr) > 0 {
			return tokenArr
		}
	}
	key, defaultValue = decoder.getTagInfo(field.Tag.Get(decoder.tag))
	if key == "" {
		return nil
	}
	token.Key = key
	token.Default = defaultValue
	token.value = &value
	tokenArr = append(tokenArr, token)
	return tokenArr
}

func (decoder *envTagDecoder) getTagInfo(tag string) (Key string, Default string) {
	if tag == "" {
		return "", ""
	}
	var pos = strings.Split(tag, ",")
	if len(pos) >= 2 {
		return strings.TrimSpace(pos[0]), strings.TrimSpace(pos[1])
	}
	return strings.TrimSpace(tag), ""
}
