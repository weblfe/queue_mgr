package service

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/weblfe/queue_mgr/entity"
)

// 用户服务
type userServiceImpl struct {
	baseServiceImpl
}

// GetUserServiceImpl 获取获取
func GetUserServiceImpl() *userServiceImpl {
	var service = new(userServiceImpl)
	service.initService(service.init)
	return service
}

func (u *userServiceImpl) init() {
	u.setType(entity.RemoteApi)
	u.setServiceID("userService")
	u.initAddr()
}

func (u *userServiceImpl) Login(params entity.UserLoginParams) (resp *entity.JsonResponse, err error) {
	var (
		args   = fiber.AcquireArgs()
		client = u.GetClient("Login.UserLogin")
	)
	args.Set("service", "Login.UserLogin")
	args.Set("user_login", params.UserAccount)
	args.Set("user_pass", params.Password)
	if err := u.PostForm(client, args); err != nil {
		return nil, err
	}
	code, body, errs := client.Bytes()
	if len(errs) <= 0 {
		errs = append(errs, NoError)
	}
	if code != 200 {
		return entity.CreateResponse(), errors.New(errs[0].Error())
	}
	return entity.CreateResponse(body), nil
}

func (u *userServiceImpl) CheckLogin(params entity.UserStateCheckParams) (resp *entity.JsonResponse, err error) {
	var (
		args   = fiber.AcquireArgs()
		client = u.GetClient("User.Iftoken")
	)
	args.Set("uid", params.GetUid())
	args.Set("token", params.Token)
	args.Set("service", "User.Iftoken")
	if err = u.PostForm(client, args); err != nil {
		return nil, err
	}
	code, body, errs := client.Bytes()
	if len(errs) <= 0 {
		errs = append(errs, NoError)
	}
	if code != 200 {
		return entity.CreateResponse(), errors.New(errs[0].Error())
	}
	return entity.CreateResponse(body), nil
}
