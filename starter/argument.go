package starter

import (
	"flag"
	"os"
)

type (
	argumentsLoader struct {
		dir        string
		configFile string
		appPath    string
	}
)

var (
	defaultArgumentsStarter = newArgumentsStarter()
)

const (
	defaultCfgFilePath = "etc/"
	defaultCfgFile     = "etc/app.yml"
)

func newArgumentsStarter() *argumentsLoader {
	return new(argumentsLoader)
}

func (loader *argumentsLoader) Init() {
	loader.commandLine()
	loader.appPath = loader.GetAppPath()
}

// 注册命令参数
func (loader *argumentsLoader) commandLine() {
	flag.StringVar(&loader.dir, "d", defaultCfgFilePath, "set app configuration file dir")
	flag.StringVar(&loader.configFile, "c", defaultCfgFile, "set app configuration file")
	flag.Parse()
}

// GetAppPath 获取应用所在路径
func (loader *argumentsLoader) GetAppPath() string {
	if loader.appPath != "" {
		return loader.appPath
	}
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return "./"
}

// GetConfigDir 获取命令行指定读取应用配置目录
func (loader *argumentsLoader) GetConfigDir() string {
	return loader.dir
}

// GetConfigFile 获取命令指定应用配置文件
func (loader *argumentsLoader) GetConfigFile() string {
	return loader.configFile
}

func GetArgumentsStarter() *argumentsLoader {
	return defaultArgumentsStarter
}
