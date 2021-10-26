package entity

import "fmt"

type ServiceType int

const (
	UnSupport ServiceType = iota // 非开发,未支持类型
	Local
	RemoteApi
	RemoteJsonRpc
	RemoteGrpcRpc
	RemoteFastCgi
	RemoteWs
)

func (t ServiceType) Desc() string {
	switch t {
	case UnSupport:
		return "unSupport"
	case Local:
		return "local"
	case RemoteApi:
		return "http"
	case RemoteJsonRpc:
		return "jsonRpc"
	case RemoteGrpcRpc:
		return "grpc"
	case RemoteFastCgi:
		return "fastCgi"
	case RemoteWs:
		return "ws"
	}
	return ""
}

func (t ServiceType) Int() int {
	return int(t)
}

func (t ServiceType) String() string {
	return fmt.Sprintf("%d", t)
}

func (t ServiceType) KvStr() string {
	return fmt.Sprintf("%d=>%s", t, t.Desc())
}
