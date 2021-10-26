package starter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
	"github.com/weblfe/queue_mgr/config"
	"github.com/weblfe/queue_mgr/router"
	"github.com/weblfe/queue_mgr/utils"
	"os"
	"time"
)

type (
	appStarter struct {
		app *fiber.App
		baseStarterConstructor
	}
)

var (
	defaultAppStarter = newAppStarter()
)

const (
	serverHeader   = "d2VibGludXhnYW1l"
	defaultAppName = "appCdnServ v0.1.0"
)

func newAppStarter() *appStarter {
	var appStarter = new(appStarter)
	appStarter.baseStarterConstructor = newStarterConstructor()
	appStarter.name = "appStarter"
	return appStarter
}

// 载入 路由
func (starter *appStarter) router() {
	var app = starter.App()
	// http
	router.Http(app)
}

func (starter *appStarter) StartUp() {
	starter.init(starter.boot)
}

func (starter *appStarter) boot() {
	time.Local, _ = time.LoadLocation(os.Getenv("TZ"))
	// 初始化 应用http engine
	starter.app = fiber.New(starter.configure())
	// 载入路由
	starter.router()
}

// 构造 http应用启动配置
func (starter *appStarter) configure() fiber.Config {
	return fiber.Config{
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  serverHeader,
		Prefork:       utils.GetEnvBool("APP_FORK_ENABLE", false), // fork
		AppName:       utils.GetEnvVal("APP_NAME", defaultAppName),
		JSONEncoder:   utils.GetJsonEncoder(),
		JSONDecoder:   utils.GetJsonDecoder(),
	}
}

func (starter *appStarter) wsConfigure() websocket.Config {
	return websocket.Config{
		HandshakeTimeout: utils.GetEnvDuration("app_ws_handshake_timeout", time.Minute),
	}
}

func (starter *appStarter) App() *fiber.App {
	return starter.app
}

func (starter *appStarter) Run() {
	var app = starter.getAppCnf()
	if err := starter.App().Listen(app.GetAddr()); err != nil {
		log.Infoln(err)
	}
	log.Infoln("stopping")
}

func (starter *appStarter) getAppCnf() config.App {
	return config.GetAppConfig().GetAppCnf()
}

func GetAppStarter() *appStarter {
	return defaultAppStarter
}
