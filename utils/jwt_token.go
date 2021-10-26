package utils

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type (
	AuthData struct {
		Uid      string
		Role     int
		Scope    string
		ExpireAt int64
		Extra    interface{}
	}

	DataSetter interface {
		Set(key string, v interface{})
	}

	DataGetter interface {
		Get(key string) interface{}
	}
)

var (
	AuthDataKeys            = []string{Uid, Role, Scope, ExpireAt, ExpireAtAlias}
	ErrorTokenExpired       = errors.New("token expired")
	ErrorTokenMissUid       = errors.New("token missing uid")
	ErrorTokenMissScope     = errors.New("token missing scope")
	ErrorTokenMissSecret    = errors.New("token missing secret")
	ErrorTokenInvalid       = errors.New("token  invalid")
	ErrorTokenMethodInvalid = errors.New("token  method error")
	ErrorMissSecret         = errors.New("missing secret")
)

const (
	DefaultSafeDuration = 7200
	SecretKey           = "APP_SECRET"
	Uid                 = "uid"
	Role                = "role"
	Scope               = "scope"
	ExpireAt            = "exp"
	ExpireAtAlias       = "expired"
)

// JwtTokenEncode 生成token
func JwtTokenEncode(data AuthData, secret ...string) (string, error) {
	if data.ExpireAt <= 0 {
		data.ExpireAt = time.Now().Add(time.Second * DefaultSafeDuration).Unix()
	}
	var token = jwt.NewWithClaims(jwt.SigningMethodHS256, createMapClaims(&data, AuthDataKeys))
	if len(secret) <= 0 || secret[0] == "" {
		secret[0] = os.Getenv(SecretKey)
	}
	var key = secret[0]
	if key == "" {
		return "", ErrorMissSecret
	}
	return token.SignedString([]byte(key))
}

// JwtTokenDecode 解析token
func JwtTokenDecode(token string, secret ...string) (*AuthData, error) {
	secret = append(secret, "")
	if len(secret) <= 0 || secret[0] == "" {
		secret[0] = os.Getenv(SecretKey)
	}
	var (
		key       = secret[0]
		data      = &AuthData{}
		mapClaims = jwt.MapClaims{}
	)
	info, err := jwt.ParseWithClaims(token, &mapClaims, createKeyFunc(key))
	if err != nil {
		return nil, err
	}
	if !info.Valid {
		return nil, ErrorTokenInvalid
	}
	for _, k := range AuthDataKeys {
		if v, ok := mapClaims[k]; ok {
			data.Set(k, v)
		}
	}
	return data, nil
}

// Set 更新更新键
func (data *AuthData) Set(key string, v interface{}) {
	switch key {
	case "uid":
		switch v.(type) {
		case string:
			data.Uid = v.(string)
		case fmt.Stringer:
			data.Uid = v.(fmt.Stringer).String()
		case fmt.GoStringer:
			data.Uid = v.(fmt.GoStringer).GoString()
		case []byte:
			data.Uid = string(v.([]byte))
		default:
			data.Uid = fmt.Sprintf("%v", v)
		}
	case "role":
		switch v.(type) {
		case int:
			data.Role = v.(int)
		case int64:
			data.Role = int(v.(int64))
		case float64:
			data.Role = int(v.(float64))
		case float32:
			data.Role = int(v.(float32))
		case string:
			var str = v.(string)
			data.Role, _ = strconv.Atoi(str)
		}
	case "scope":
		switch v.(type) {
		case string:
			data.Scope = v.(string)
		case fmt.Stringer:
			data.Scope = v.(fmt.Stringer).String()
		case fmt.GoStringer:
			data.Scope = v.(fmt.GoStringer).GoString()
		case []byte:
			data.Scope = string(v.([]byte))
		default:
			data.Scope = fmt.Sprintf("%v", v)
		}
	case "exp","expired","expiredAt":
		switch v.(type) {
		case int:
			data.ExpireAt = int64(v.(int))
		case int64:
			data.ExpireAt = v.(int64)
		case float64:
			data.ExpireAt = int64(v.(float64))
		case float32:
			data.ExpireAt = int64(v.(float32))
		case string:
			var str = v.(string)
			data.ExpireAt, _ = strconv.ParseInt(str, 0, 64)
		}
	}
}

func (data *AuthData) Get(key string) interface{} {
	switch key {
	case "uid":
		return data.Uid
	case "role":
		return data.Role
	case "scope":
		return data.Scope
	case "exp":
		return data.ExpireAt
	case "extra":
		return data.Extra
	}
	return nil
}

func (data *AuthData) Verify() error {
	var (
		now   = time.Now()
		timer = time.Unix(data.ExpireAt, 0)
	)
	if timer.Before(now) {
		return ErrorTokenExpired
	}
	if data.Uid == "" {
		return ErrorTokenMissUid
	}
	if data.Scope == "" {
		return ErrorTokenMissScope
	}
	return nil
}

func (data *AuthData) Encode() string {
	var builder = url.Values{
		"uid":     {data.Uid},
		"scope":   {data.Scope},
		"extra":   {JsonEncode(data.Extra).String()},
		"exp":     {time.Unix(data.ExpireAt, 0).Format(`2006-01-02 15:04:05`)},
		"expired": {fmt.Sprintf("%d", data.ExpireAt)},
		"role":    {fmt.Sprintf("%d", data.Role)},
	}
	return builder.Encode()
}

func (data *AuthData)String() string {
	return JsonEncode(map[string]interface{}{
		"uid":     data.Uid,
		"scope":   data.Scope,
		"extra":   JsonEncode(data.Extra).String(),
		"exp":     time.Unix(data.ExpireAt, 0).Format(`2006-01-02 15:04:05`),
		"expired": fmt.Sprintf("%d", data.ExpireAt),
		"role":    fmt.Sprintf("%d", data.Role),
	}).String()
}

func (data *AuthData) IsExpiredError(err error) bool {
	return err == ErrorTokenExpired
}

func (data *AuthData) IsMissUidError(err error) bool {
	return err == ErrorTokenMissUid
}

func (data *AuthData) IsMissScopeError(err error) bool {
	return err == ErrorTokenMissScope
}

func (data *AuthData) CheckScope(scope ...string) bool {
	for _, v := range scope {
		if v == data.Scope || strings.ToLower(v) == data.Scope {
			return true
		}
	}
	return false
}

// 创建 秘钥解析函数
func createKeyFunc(secret string) func(*jwt.Token) (interface{}, error) {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrorTokenMethodInvalid
		}
		if secret == "" {
			return nil, ErrorTokenMissSecret
		}
		return []byte(secret), nil
	}
}

// 创建 map Claims
func createMapClaims(data DataGetter, keys ...[]string) jwt.MapClaims {
	var claims = jwt.MapClaims{}
	keys = append(keys, AuthDataKeys)
	if len(keys) <= 0 || len(keys[0]) <= 0 {
		return claims
	}
	for _, k := range keys[0] {
		claims[k] = data.Get(k)
	}
	return claims
}
