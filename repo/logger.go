package repo

import (
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type (
	loggerRepository struct {
		cache   sync.Map
		writers []io.WriteCloser
	}

	loggerParams struct {
		Driver string
		File   string
		Level  string
		Writer io.Writer
		Hooks  []string
	}
)

var (
	helper       = newLoggerRepo()
	defaultHooks = []string{"rotate", "error", "webHookNotify"}
)

// 构建新logger 服务
func newLoggerRepo() *loggerRepository {
	var loggerRepo = &loggerRepository{
		cache: sync.Map{},
	}
	runtime.SetFinalizer(loggerRepo, (*loggerRepository).destroy)
	return loggerRepo
}

// GetLogger 获取日志对象
func GetLogger(name ...string) *log.Logger {
	name = append(name, "default")
	return helper.GetLogger(name[0])
}

func (repo *loggerRepository) GetLogger(name string) *log.Logger {
	var app = strings.ToUpper(name)
	logger, ok := repo.cache.Load(app)
	if !ok {
		logApp := repo.createLogger(app)
		if logApp != nil {
			repo.cache.Store(app, logApp)
			return logApp
		}
		return repo.getLogger()
	}
	switch logger.(type) {
	case *log.Logger:
		return logger.(*log.Logger)
	}
	return repo.getLogger()
}

func (repo *loggerRepository) createLogger(name string) *log.Logger {
	var params = repo.getLoggerParams(name)
	switch params.Driver {
	case "file":
		params.Writer = repo.getFileWriter(params.File)
	case "stdout":
		params.Writer = os.Stdout
	case "stderr":
		params.Writer = os.Stderr
	case "multi":
		params.Writer = repo.getMultiWriter(params.File)
	default:
		if params.Writer == nil {
			params.Writer = os.Stdout
		}
	}
	logger := log.New()
	logger.Out = params.Writer
	logger.Level = params.getLevel()
	return logger
}

func (repo *loggerRepository) getMultiWriter(params string) io.Writer {
	if params == "" {
		return os.Stdout
	}
	var (
		ioArr []io.Writer
		arr   = strings.Split(params, ",")
	)
	for _, it := range arr {
		ioArr = append(ioArr, repo.getFileWriter(it))
	}
	return io.MultiWriter(ioArr...)
}

func (repo *loggerRepository) getFileWriter(file string) io.Writer {
	if file == "" {
		log.Infoln("empty file parameter")
		return os.Stdout
	}
	switch file {
	case "stderr":
		return os.Stderr
	case "stdout":
		return os.Stdout
	}
	fileAbs, _err := filepath.Abs(file)
	if _err != nil {
		log.Infoln(_err)
		return os.Stdout
	}
	file = fileAbs
	_, err := os.Stat(file)
	if err == nil {
		fd, _err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if _err != nil {
			log.Infoln(_err)
			return os.Stdout
		}
		repo.writers = append(repo.writers, fd)
		return fd
	}
	if os.IsNotExist(err) {
		dir := filepath.Dir(file)
		_ = os.MkdirAll(dir, os.ModePerm)
		fd, err2 := os.OpenFile(file, os.O_CREATE|os.O_RDWR, os.ModePerm)
		if err2 != nil {
			log.Infoln(err)
			return os.Stdout
		}
		repo.writers = append(repo.writers, fd)
		return fd
	}
	log.Infoln(err)
	return os.Stdout
}

func (repo *loggerRepository) getLoggerParams(name string) loggerParams {
	if !strings.HasSuffix(name, "_") {
		name = name + "_"
	}
	var (
		hooks     []string
		logFile   = strings.ToUpper(name + "LOGGER_FILE")
		logDriver = strings.ToUpper(name + "LOGGER_DRIVER")
		logLevel  = strings.ToUpper(name + "LOGGER_LEVEL")
		logHooks  = strings.ToUpper(name + "LOGGER_HOOKS")
		file      = os.Getenv(logFile)
		driver    = os.Getenv(logDriver)
		level     = os.Getenv(logLevel)
		hook      = os.Getenv(logHooks)
	)
	if driver == "" {
		driver = "stdout"
	}
	if file == "" {
		file = "stdout"
	}
	if level == "" {
		level = "info"
	}
	// 文件驱动
	if file != "stdout" && file != "stderr" && driver == "stdout" {
		driver = "file"
	}
	hook = strings.TrimSpace(hook)
	if hook == "" {
		hooks = repo.getDefaultHooks()
	} else {
		hooks = strings.Split(hook, ",")
	}

	return loggerParams{
		Driver: driver,
		File:   file,
		Level:  strings.ToLower(level),
		Hooks:  hooks,
	}
}

func (repo *loggerRepository) getDefaultHooks() []string {
	var hooks = os.Getenv("LOGGER_HOOKS")
	if hooks == "" {
		return []string{}
	}
	return strings.Split(strings.TrimSpace(hooks), ",")
}

func (repo *loggerRepository) getLogger() *log.Logger {
	var logger = log.New()
	logger.Out = os.Stdout
	logger.Formatter = &log.JSONFormatter{}
	return logger
}

func (params loggerParams) getLevel() log.Level {
	switch params.Level {
	case "fatal":
		return log.FatalLevel
	case "error":
		return log.ErrorLevel
	case "warn":
		return log.WarnLevel
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "trace":
		return log.TraceLevel
	}
	return log.InfoLevel
}

func (repo *loggerRepository) destroy() {
	runtime.SetFinalizer(repo, nil)
	if repo.writers == nil {
		return
	}
	for _, closer := range repo.writers {
		if closer != nil {
			_ = closer.Close()
		}
	}
	repo.writers = nil
}
