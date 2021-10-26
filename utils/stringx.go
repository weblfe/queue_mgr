package utils

import (
	"fmt"
	"strings"
)

func Str2Arr(content string) []string {
	if json.Valid([]byte(content)) {
		var arr = new([]string)
		err := JsonDecode([]byte(content), arr)
		if err == nil {
			return *arr
		}
		fmt.Println("GetEnvMapper Error:", err)
	}
	if strings.Contains(content, ",") {
		return strings.Split(content, ",")
	}
	if strings.Contains(content, " ") {
		return strings.Split(content, " ")
	}
	if strings.Contains(content, "\n") {
		return strings.Split(content, "\n")
	}
	return []string{content}
}

func Slice2StrArr(array []interface{}) []string {
	var arr []string
	for _, v := range array {
		switch v.(type) {
		case string:
			arr = append(arr, v.(string))
		case *string:
			arr = append(arr, *v.(*string))
		case fmt.Stringer:
			arr = append(arr, v.(fmt.Stringer).String())
		default:
			arr = append(arr, fmt.Sprintf("%v", v))
		}
	}
	return arr
}

type stringer struct {
	value interface{}
	err   error
}

func NewStringer(v interface{}) *stringer {
	return &stringer{
		value: v,
	}
}

func (v *stringer) Error() error {
	return v.err
}

func (v *stringer) Convert() string {
	var data = v.value
	switch data.(type) {
	case string:
		return data.(string)
	case []byte:
		return string(data.([]byte))
	case fmt.Stringer:
		return data.(fmt.Stringer).String()
	case fmt.GoStringer:
		return data.(fmt.GoStringer).GoString()
	}
	return fmt.Sprintf("%v", data)
}

func (v *stringer)String() string  {
	return v.Convert()
}

func (v *stringer)GoString() string  {
	return v.Convert()
}
