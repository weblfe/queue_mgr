package entity

import "fmt"

type Code int

const (
	CodeVerify      Code = 2001 // 参数效验不通过
	CodeParamNil    Code = 2004 // 参数缺失
	CodeSystemError Code = 5001 // 系统异常
	CodeUndefined   Code = 4001 // 系统未定义异常
	CodeExits       Code = 2000 // 记录已存在
	CodeFail        Code = 1003 // 业务服务失败
	CodeSuccess     Code = 0    // 服务逻辑成功
	CodeParamDecode Code = -1   // 参数传输异常
)

func (code Code) Int() int {
	return int(code)
}

func (code Code) String() string {
	return fmt.Sprintf("%d", code)
}
