package libs

import (
	"github.com/sirupsen/logrus"
	"github.com/yuin/gopher-lua"
)

type (
	LuaFunctionTable struct {
		logger *logrus.Logger
	}
)

var (
	LoggerModule = "logger"
	LoggerFuncs  = map[string]lua.LGFunction{
		"logInfo":      logInfo,
		"logInfoLn":    logInfoLn,
		"logError":     logError,
		"logErrorLn":   logErrorLn,
		"logDebug":     logDebug,
		"logDebugLn":   logDebugLn,
		"createLogger": createLogger,
	}
)

func NewLuaLoggerTables() lua.LGFunction {
	return func(L *lua.LState) int {
		var mod lua.LValue
		if len(LoggerFuncs) <= 0 {
			return 0
		}
		mod = L.RegisterModule(LoggerModule, LoggerFuncs)
		L.Push(mod)
		return 1
	}
}

func logInfoLn(L *lua.LState) int {
	var argc = L.GetTop()
	if argc <= 0 {
		return 0
	}
	var args []interface{}
	for i := 0; i < argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		args = append(args, v.String())
	}
	if len(args) <= 0 {
		return 0
	}
	logrus.Infoln(args...)
	return 1
}

func logInfo(L *lua.LState) int {
	var argc = L.GetTop()
	if argc <= 0 {
		return 0
	}
	var args []interface{}
	for i := 0; i < argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		args = append(args, v.String())
	}
	if len(args) <= 0 {
		return 0
	}
	logrus.Info(args...)
	return 1
}

func logDebug(L *lua.LState) int {
	var argc = L.GetTop()
	if argc <= 0 {
		return 0
	}
	var args []interface{}
	for i := 0; i < argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		args = append(args, v.String())
	}
	if len(args) <= 0 {
		return 0
	}
	logrus.Debug(args...)
	return 1
}

func logDebugLn(L *lua.LState) int {
	var argc = L.GetTop()
	if argc <= 0 {
		return 0
	}
	var args []interface{}
	for i := 0; i < argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		args = append(args, v.String())
	}
	if len(args) <= 0 {
		return 0
	}
	logrus.Debugln(args...)
	return 1
}

func logErrorLn(L *lua.LState) int {
	var argc = L.GetTop()
	if argc <= 0 {
		return 0
	}
	var args []interface{}
	for i := 0; i < argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		args = append(args, v.String())
	}
	if len(args) <= 0 {
		return 0
	}
	logrus.Errorln(args...)
	return 1
}

func logError(L *lua.LState) int {
	var argc = L.GetTop()
	if argc <= 0 {
		return 0
	}
	var args []interface{}
	for i := 0; i < argc; i++ {
		v := L.CheckAny(i)
		if v == lua.LNil {
			continue
		}
		args = append(args, v.String())
	}
	if len(args) <= 0 {
		return 0
	}
	logrus.Error(args...)
	return 1
}

// function createLogger(file string,level int,mode int) logger
func createLogger(L *lua.LState) int {
	//var argc = L.GetTop()
	return 1
}
