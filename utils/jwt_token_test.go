package utils

import (
	"fmt"
	"github.com/subosito/gotenv"
	"os"
	"testing"
	"time"
)

func TestJwtTokenDecode(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var (
		secret     = os.Getenv("APP_SECRET")
		token     = os.Getenv("TEST_APP_TOKEN")
		data, err = JwtTokenDecode(token,secret)
	)

	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%v\n",data.String())
	if data != nil {
		if err = data.Verify(); err != nil {
			t.Error(err)
		}
	}
}

func TestJwtTokenEncode(t *testing.T) {
	if err := gotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	var (
		data = AuthData{
			Uid:      "123",
			Role:     1,
			Scope:    os.Getenv("APP_SCOPE"),
			ExpireAt: time.Now().Add(time.Hour).Unix(),
		}
		secret     = os.Getenv("APP_SECRET")
		token, err = JwtTokenEncode(data, secret)
	)

	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("token 生成失败")
	}
	fmt.Println("token: ", token)
}
