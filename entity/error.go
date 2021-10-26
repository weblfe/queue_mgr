package entity

import "errors"

var (
	// ErrorEmpty 空记录
	ErrorEmpty = errors.New("not exist")
	// ErrorPrivilegeLogin 请重新登录
	ErrorPrivilegeLogin = errors.New("please login")
	// ErrorDiType 类型不正确
	ErrorDiType = errors.New("container resolve type error")
	// ErrorRequired  推送参数异常
	ErrorRequired = errors.New("param required but given nil")
	// ErrorDecodeFailed  参数解码失败
	ErrorDecodeFailed  = errors.New("request param decode failed")
	// ErrorSupport 未知支持类型
	ErrorSupport = errors.New("unknown support type")
)

func IsLoginError(err error) bool {
	return err == ErrorPrivilegeLogin
}

func IsDiTypeError(err error) bool {
	return err == ErrorDiType
}

func IsRequiredError(err error) bool {
	return err == ErrorRequired
}

func IsEmptyError(err error) bool {
	return err == ErrorEmpty
}

func IsDecodeFailError(err error) bool {
		return err == ErrorDecodeFailed
}

func IsSupportError(err error) bool {
	return err == ErrorSupport
}
