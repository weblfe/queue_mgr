package starter

import (
	log "github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	"os"
	"path/filepath"
	"sync"
)

type envStarter struct {
	envFile string
	lock    sync.RWMutex
	baseStarterConstructor
}

var defaultEnvStarter = newEnvStarter()

func GetEnvStarter() *envStarter {
	return defaultEnvStarter
}

func newEnvStarter() *envStarter {
	var _starter = envStarter{}
	_starter.envFile = ""
	_starter.lock = sync.RWMutex{}
	_starter.baseStarterConstructor = newStarterConstructor()
	_starter.name = "envStarter"
	return &_starter
}

func (starter *envStarter) Init() {
	starter.lock.Lock()
	defer starter.lock.Unlock()
	if starter.envFile == "" {
		starter.envFile = starter.getFile()
	}
	return
}

// 获取文件路径
func (starter *envStarter) getFile() string {
	var appDir = GetArgumentsStarter().GetAppPath()
	if appDir == "" {
		dir, err := os.Getwd()
		if err == nil {
			panic(err)
		}
		appDir = dir
	}
	var (
		path      = filepath.Join(appDir, ".env")
		file, err = filepath.Abs(path)
	)
	if err != nil {
		return ""
	}
	return file
}

// StartUp 启动
func (starter *envStarter) StartUp() {
	starter.init(starter.boot)
}

// 加载环境
func (starter *envStarter) boot() {
	if starter.envFile == "" {
		return
	}
	if _, err := os.Stat(starter.envFile); err != nil {
		if os.IsNotExist(err) {
			return
		}
	}
	err := gotenv.Load(starter.envFile)
	if err != nil {
		panic(err)
	}
	log.Infoln("environment loaded " + starter.envFile + " ok ")
}
