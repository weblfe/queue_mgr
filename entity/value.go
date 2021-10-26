package entity

import (
	"encoding/json"
	"fmt"
	"github.com/weblfe/queue_mgr/utils"
	"reflect"
	"strconv"
	"time"
)

type valueConverter struct {
	value interface{}
}

func NewValue(v interface{}) *valueConverter {
	var conv = new(valueConverter)
	return conv.setValue(v)
}

func (value *valueConverter) setValue(v interface{}) *valueConverter {
	if v == nil {
		return value
	}
	value.value = v
	return value
}

func (value *valueConverter) Typeof() string {
	var ref = reflect.ValueOf(value.value)
	return ref.Type().Kind().String()
}

func (value *valueConverter) Empty() bool {
	if value == nil || value.value == nil || value.value == "" || value.value == 0 {
		return true
	}
	return false
}

func (value *valueConverter) IsNull() bool {
	return value == nil || value.value == nil
}

func (value *valueConverter) IsZero() bool {
	var ref = reflect.ValueOf(value.value)
	return ref.IsZero()
}

func (value *valueConverter) Interface() interface{} {
	return value.value
}

func (value *valueConverter) String() string {
	if value.value == nil {
		return ""
	}
	var (
		v    = reflect.ValueOf(value.value)
		data = v.Interface()
	)
	switch data.(type) {
	case string:
		return data.(string)
	case fmt.Stringer:
		return data.(fmt.Stringer).String()
	case fmt.GoStringer:
		return data.(fmt.GoStringer).GoString()
	case json.Marshaler:
		var (
			marshal    = data.(json.Marshaler)
			bytes, err = marshal.MarshalJSON()
		)
		if err != nil {
			return ""
		}
		return string(bytes)
	}
	return fmt.Sprintf("%v", data)
}

func (value *valueConverter) IsNumber() bool {
	var (
		v    = reflect.ValueOf(value.value)
		kind = v.Kind()
	)
	if kind == reflect.Interface || kind == reflect.Ptr {
		kind = v.Elem().Kind()
	}
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func (value *valueConverter) Int() int {
	var (
		v    = value.valueOf()
		kind = v.Kind()
	)
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int(v.Uint())
	case reflect.Float32, reflect.Float64:
		return int(v.Float())
	case reflect.String:
		var (
			str    = v.String()
			n, err = strconv.Atoi(str)
		)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

func (value *valueConverter) Float() float64 {
	var (
		v    = value.valueOf()
		kind = v.Kind()
	)
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.String:
		var (
			str    = v.String()
			n, err = strconv.ParseFloat(str, 64)
		)
		if err != nil {
			return 0
		}
		return n
	default:
		return 0
	}
}

func (value *valueConverter) Duration() time.Duration {
	if value.value == nil {
		return 0
	}
	var (
		v    = reflect.ValueOf(value.value)
		data = v.Interface()
	)
	switch data.(type) {
	case time.Duration:
		return data.(time.Duration)
	case *time.Duration:
		return *data.(*time.Duration)
	case string:
		var (
			str    = data.(string)
			d, err = time.ParseDuration(str)
		)
		if err != nil {
			return 0
		}
		return d
	case fmt.Stringer:
		var (
			str    = data.(fmt.Stringer).String()
			d, err = time.ParseDuration(str)
		)
		if err != nil {
			return 0
		}
		return d
	case fmt.GoStringer:
		var (
			str    = data.(fmt.GoStringer).GoString()
			d, err = time.ParseDuration(str)
		)
		if err != nil {
			return 0
		}
		return d
	case json.Marshaler:
		var (
			marshal    = data.(json.Marshaler)
			bytes, err = marshal.MarshalJSON()
		)
		if err != nil {
			return 0
		}
		var (
			str = string(bytes)
		)
		if d, err := time.ParseDuration(str); err == nil {
			return d
		}
	default:
		return 0
	}
	return 0
}

func (value *valueConverter) Time() (time.Time, bool) {
	if value.value == nil {
		return time.Time{}, false
	}
	var (
		v    = reflect.ValueOf(value.value)
		data = v.Interface()
	)
	switch data.(type) {
	case time.Time:
		return data.(time.Time), true
	case *time.Time:
		return *data.(*time.Time), true
	case string:
		var str = data.(string)
		if t := value.convertTime(str); t != nil {
			return *t, true
		}
	case fmt.Stringer:
		var str = data.(fmt.Stringer).String()
		if t := value.convertTime(str); t != nil {
			return *t, true
		}
	case fmt.GoStringer:
		var str = data.(fmt.GoStringer).GoString()
		if t := value.convertTime(str); t != nil {
			return *t, true
		}
	case json.Marshaler:
		var (
			marshal    = data.(json.Marshaler)
			bytes, err = marshal.MarshalJSON()
		)
		if err != nil {
			return time.Time{}, false
		}
		if t := value.convertTime(string(bytes)); t != nil {
			return *t, true
		}
	default:
		return time.Time{}, false
	}
	return time.Time{}, false
}

func (value *valueConverter) convertTime(str string) *time.Time {
	if str == "" {
		return nil
	}
	if t, err := time.Parse(utils.DateTimeLayout, str); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RFC3339Nano, str); err == nil {
		return &t
	}
	if t, err := time.Parse(time.UnixDate, str); err == nil {
		return &t
	}
	if t, err := time.Parse(time.RubyDate, str); err == nil {
		return &t
	}
	return nil
}

func (value *valueConverter) valueOf() reflect.Value {
	var (
		v    = reflect.ValueOf(value.value)
		kind = v.Kind()
	)
	if kind == reflect.Interface || kind == reflect.Ptr {
		return v.Elem()
	}
	return v
}
