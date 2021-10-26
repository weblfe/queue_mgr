package utils

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
)

var (
	json           = jsoniter.ConfigCompatibleWithStandardLibrary
	EmptyJsonError = errors.New("empty json result")
)

type Json struct {
	data []byte
	err  error
}

func NewJson(data []byte, err error) Json {
	return Json{
		data, err,
	}
}

func (d Json) String() string {
	return string(d.data)
}

func (d Json) Empty() bool {
	if d.err != nil || len(d.data) == 0 {
		return true
	}
	return false
}

func (d Json) EmptyErr() error {
	return EmptyJsonError
}

func (d Json) IsEmptyErr(err error) bool {
	return err == EmptyJsonError
}

func (d Json) Bytes() []byte {
	return d.data
}

func (d Json) HasErr() bool {
	if d.err != nil {
		return true
	}
	return false
}

func (d Json) Error() error {
	return d.err
}

func (d Json) Decode(addr interface{}) error {
	return JsonDecode(d.Bytes(), addr)
}

// JsonEncode object序列化
func JsonEncode(data interface{}) Json {
	return NewJson(json.Marshal(data))
}

// JsonDecode json反序列化
func JsonDecode(data []byte, addr interface{}) error {
	return json.Unmarshal(data, addr)
}

func GetJsonDecoder() func([]byte, interface{}) error {
	return json.Unmarshal
}

func GetJsonEncoder() func(interface{}) ([]byte, error) {
	return json.Marshal
}
