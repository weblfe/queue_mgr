package facede

import "github.com/weblfe/queue_mgr/entity"

// UserService 用户服务
type UserService interface {
	// Login 用户登录
	Login(params entity.UserLoginParams) (resp *entity.JsonResponse,err error)
	// CheckLogin 检查登录
	CheckLogin(params entity.UserStateCheckParams) (resp *entity.JsonResponse,err error)
}
