package starter

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/weblfe/queue_mgr/config"
	"os"
)

type (
	settingStarter struct {
		cfg        *viper.Viper
		configDir  string
		configFile string
		baseStarterConstructor
	}
)

var (
	defaultSettingStarter = newSettingStarter()
)

const configFileName = "app.yml"

func GetSettingStarter() *settingStarter {
	return defaultSettingStarter
}

// 构建
func newSettingStarter() *settingStarter {
	var settings = new(settingStarter)
	settings.baseStarterConstructor = newStarterConstructor()
	settings.name = "settingStarter"
	return settings
}

// StartUp 启动 初始化 全局配置
func (starter *settingStarter) StartUp() {
	starter.init(starter.boot)
}

// GetViper 获取配置解析器
func (starter *settingStarter) GetViper() *viper.Viper {
	return starter.cfg
}

// 加载配置
func (starter *settingStarter) boot() {
	starter.initArgs()
	starter.initParser()
	// 应用 配置初始化
	config.GetAppConfig().LoadConfiguration(starter.cfg)
}

// 初始化配置读取参数
func (starter *settingStarter) initArgs() {
	if starter.configDir == "" {
		starter.configDir = GetArgumentsStarter().GetConfigDir()
	}
	if starter.configFile == "" {
		starter.configFile = GetArgumentsStarter().GetConfigFile()
	}
}

// 解析配置
func (starter *settingStarter) initParser() {
	if starter.cfg == nil {
		starter.cfg = viper.New()
	}
	if starter.cfg != nil && starter.configDir != "" && starter.exists(starter.configDir) {
		starter.cfg.AddConfigPath(starter.configDir)
		starter.cfg.SetConfigName(configFileName)
	}
	if starter.cfg != nil && starter.configFile != "" && starter.exists(starter.configFile) {
		starter.cfg.SetConfigFile(starter.configFile)
	}
	if err := starter.cfg.ReadInConfig(); err != nil {
		log.Error("ReadInConfig Error:" + err.Error())
	}
}

func (starter *settingStarter) exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}
