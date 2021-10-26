package facede

import "fmt"

// CfgKv 配置对象
type CfgKv interface {
	Keys() []string
	MAdd(map[string]interface{}) int
	Decode(content []byte) error
	ValueOf(string,...interface{}) interface{}
	fmt.Stringer
}
