package plugins

import (
	"fmt"
	"testing"
)

func TestNewLuaPlugin(t *testing.T) {
	var plugin = NewLuaPlugin()
	plugin.SetLoader(CreateExtendsLoader).Boot()
	if state := plugin.GetLState(); state == nil {
		t.Error("error lua state")
	}
}

func TestLuaPluginImpl_EvalExpr(t *testing.T) {
	var plugin = NewLuaPlugin()
	plugin.SetLoader(CreateExtendsLoader).Boot()
	if state := plugin.GetLState(); state == nil {
		t.Error("error lua state")
	}
	var luaScript = `
log=require("logger")
os=require("os")
a=1+1
table={}
print("hello go-lua!")
print("a:",a)
log.logInfoLn("a:",a)
log.logInfoLn(log)
log.logInfoLn("env home:",os.getenv("HOME"))
`
	if err := plugin.EvalExpr(luaScript); err != nil {
		t.Error("lua error", err.Error())
	}
}

func TestLuaPluginImpl_Libs(t *testing.T) {
	var plugin = NewLuaPlugin()
	plugin.SetLoader(CreateExtendsLoader).Boot()
	if state := plugin.GetLState(); state == nil {
		t.Error("error lua state")
	}
	libs := plugin.Libs()
	if len(libs) <= 0 {
		t.Error("list lua libs failed")
	} else {
		for i, v := range libs {
			fmt.Println(fmt.Sprintf("lua[%d].%s", i, v))
		}
	}
}
