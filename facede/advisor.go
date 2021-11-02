package facede

// Advisor 通知
type Advisor interface {
	SetOpts(key string, value interface{})
	Notification(data []byte, level ...string) error
}
