package entity

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/utils"
	"strconv"
)

// UserLoginParams 用户登录参数
type UserLoginParams struct {
	UserAccount string `json:"user_login"` // 账号
	Password    string `json:"user_pass"`  // 密码
}

func (params *UserLoginParams) String() string {
	return utils.JsonEncode(params).String()
}

// UserStateCheckParams 用户登录参数
type UserStateCheckParams struct {
	Uid   int    `json:"uid"`   // 用户Id
	Token string `json:"token"` // 登录token
}

func (params *UserStateCheckParams) String() string {
	return utils.JsonEncode(params).String()
}

func (params *UserStateCheckParams) GetUid() string {
	return fmt.Sprintf("%d", params.Uid)
}

type UserLoginResponse struct {
	HttpCode int       `json:"ret"`
	Data     *UserData `json:"data,omitempty"`
	Msg      string    `json:"msg"`
}

type UserData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Info []*UserInfo `json:"info,omitempty"`
}

type UserInfo struct {
	Avatar        string `json:"avatar,omitempty"`
	AvatarThumb   string `json:"avatar_thumb"`
	Birthday      string `json:"birthday"`
	City          string `json:"city"`
	Coin          string `json:"coin,default=0"`
	Consumption   string `json:"consumption"`
	ID            string `json:"id"`
	IsReg         string `json:"isreg"`
	LastLoginTime string `json:"last_login_time"`
	Level         string `json:"level"`
	LevelAnchor   string `json:"level_anchor"`
	Location      string `json:"location"`
	LoginType     string `json:"login_type"`
	Province      string `json:"province"`
	Sex           string `json:"sex"`
	Signature     string `json:"signature"`
	Token         string `json:"token,omitempty"`
	UserNicename  string `json:"user_nicename"`
	Votestotal    string `json:"votestotal"`
}

func (u *UserInfo) GetUserId() int {
	if u.ID != "" {
		n, err := strconv.Atoi(u.ID)
		if err != nil {
			return n
		}
	}
	return 0
}

func (u *UserInfo) String() string {
	return utils.JsonEncode(u).String()
}

func (u *UserInfo) GetCoin() int {
	if u.Coin != "" {
		n, err := strconv.Atoi(u.Coin)
		if err != nil {
			return n
		}
	}
	return 0
}

func (u *UserInfo) GetToken() string {
	return u.Token
}

func (resp *UserLoginResponse) GetUser() *UserInfo {
	if resp.IsSuccess() {
		return nil
	}
	if resp.Data.Info == nil {
		return nil
	}
	user := resp.Data.Info[0]
	if user != nil {
		return user
	}
	return nil
}

func (resp *UserLoginResponse) String() string {
	return utils.JsonEncode(resp).String()
}

func (resp *UserLoginResponse) IsSuccess() bool {
	if resp.HttpCode != 200 {
		return false
	}
	if resp.Data != nil && resp.Data.Code != 0 {
		return false
	}
	return true
}

func CreateUserLoginResponse(data ...[]byte) *UserLoginResponse {
	var resp = &UserLoginResponse{}
	if len(data) <= 0 {
		if err := utils.JsonDecode(data[0], resp); err != nil {
			log.Infoln(err)
		}
	}
	return &UserLoginResponse{}
}
