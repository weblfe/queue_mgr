package middlewares

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/facede"
	"github.com/weblfe/queue_mgr/repo"
	"github.com/weblfe/queue_mgr/utils"
	"os"
)

type (
	Jwt struct {
		secret    string
		scope     string
		appID     string
		headerKey string
		appIDKey  string
		skips     []string
		debug     uint8
		logger    *logrus.Logger
		storage   facede.SecretStorage
	}

	Option struct {
		Secret    string
		Scope     string
		HeaderKey string
	}
)

func NewJwt() *Jwt {
	var jwt = new(Jwt)
	jwt.scope = utils.GetEnvVal("JWT_SCOPE", "")
	jwt.headerKey = utils.GetEnvVal("JWT_HEADER_KEY", "Authorization")
	return jwt
}

func NewJwtWare(options ...*Option) fiber.Handler {
	var ware = NewJwt()
	if len(options) <= 0 || options[0] == nil {
		return ware.Handler
	}
	var opt = options[0]
	ware.SetScope(opt.Scope).SetSecret(opt.Secret).SetHeaderKey(opt.HeaderKey)
	return ware.Handler
}

func (ware *Jwt) SetSecret(secret string) *Jwt {
	if ware.secret == "" && secret != "" {
		ware.secret = secret
	}
	return ware
}

func (ware *Jwt) SetAppIDKey(key string) *Jwt {
	if ware.appIDKey == "" && key != "" {
		ware.appIDKey = key
	}
	return ware
}

func (ware *Jwt) SetScope(scope string) *Jwt {
	if ware.scope == "" && scope != "" {
		ware.scope = scope
	}
	return ware
}

func (ware *Jwt) SetHeaderKey(key string) *Jwt {
	if ware.headerKey == "" && key != "" {
		ware.headerKey = key
	}
	return ware
}

func (ware *Jwt) GetKey() string {
	if ware.headerKey == "" {
		return "Authorization"
	}
	return ware.headerKey
}

func (ware *Jwt) decode(c *fiber.Ctx) error {
	var (
		key         = ware.GetKey()
		accessToken = c.Get(key)
	)
	if accessToken == "" {
		return errors.New("miss access token")
	}
	var appID = ware.getAppID(c)
	if appID == "" {
		return errors.New("miss appID")
	}
	var secret = ware.GetAppSecret(appID)
	if secret == "" {
		return errors.New("miss app secret")
	}
	var data, err = utils.JwtTokenDecode(accessToken, secret)
	if err != nil {
		return err
	}
	if err = data.Verify(); err != nil {
		return err
	}
	if ware.scope != "" {
		if data.Scope != ware.getScope() && data.Scope != "*" {
			return errors.New("scope error")
		}
	}
	c.Request().Header.Add("X-UID", data.Uid)
	c.Request().Header.Add("X-Role", fmt.Sprintf("%d", data.Role))
	return nil
}

func (ware *Jwt) GetAppSecret(appID string) string {
	// 本应用AppID
	if id := ware.getID(); id != "" && appID == id {
		return ware.getSecret()
	}
	var storage = ware.getSecretStorage()
	if storage == nil {
		return ""
	}
	return storage.GetSecretByAppID(appID)
}

func (ware *Jwt) getSecret() string {
	if ware.secret == "" {
		ware.secret = os.Getenv("APP_SECRET")
	}
	return ware.secret
}

func (ware *Jwt) getID() string {
	if ware.appID == "" {
		ware.appID = os.Getenv("APP_ID")
	}
	return ware.appID
}

func (ware *Jwt) getScope() string {
	if ware.scope == "" {
		ware.scope = os.Getenv("APP_SCOPE")
	}
	return ware.scope
}

func (ware *Jwt) getSecretStorage() facede.SecretStorage {
	if ware.storage == nil {
		//ware.storage = models.NewAppInfo()
	}
	return ware.storage
}

func (ware *Jwt) getAppID(c *fiber.Ctx) string {
	var key = ware.getAppIDKey()
	if appID := c.Get(key); appID != "" {
		return appID
	}
	var id = c.Params("appID")
	if id == "" {
		return ware.getID()
	}
	return id
}

func (ware *Jwt) getAppIDKey() string {
	if ware.appIDKey == "" {
		return "AppID"
	}
	return ware.appIDKey
}

func (ware *Jwt) skip(c *fiber.Ctx) bool {
	var path = c.Path()
	if len(ware.skips) > 0 {
		for _, v := range ware.skips {
			if v == path {
				return true
			}
		}
	}
	if ware.GetDebug() {
		return true
	}
	return false
}

func (ware *Jwt) GetDebug() bool {
	if ware.debug == 0 {
		if utils.GetEnvBool("JWT_FILTER_OFF") {
			ware.debug = 1
		} else {
			ware.debug = 2
		}
	}
	if ware.debug == 1 {
		return true
	}
	return false
}

func (ware *Jwt) Handler(c *fiber.Ctx) error {
	if c.Method() == fiber.MethodOptions || ware.skip(c) {
		return c.Next()
	}
	if err := ware.decode(c); err == nil {
		return c.Next()
	}
	if c.Get(fiber.HeaderAccept) == fiber.MIMEApplicationJSON {
		var data = fiber.Map{
			"httpCode": fiber.StatusUnauthorized,
			"data": fiber.Map{
				"code": fiber.StatusUnauthorized,
				"msg":  `please try login`,
			},
		}
		ware.getLogger().WithField("request",c.Request()).Errorln("jwt decode failed")
		return c.Status(fiber.StatusUnauthorized).JSON(data)
	}
	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	return c.Status(fiber.StatusUnauthorized).Send([]byte(`please try login`))
}

func (ware *Jwt) getLogger() *logrus.Logger {
	if ware.logger == nil {
		ware.logger = repo.GetLogger("middleware")
	}
	return ware.logger
}
