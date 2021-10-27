package plugins

import (
	"testing"
)

func TestNewLuaPlugin(t *testing.T) {
	var plugin = NewLuaPlugin()
	plugin.SetLoader(CreateLoader).Boot()
	if state := plugin.GetLState(); state == nil {
		t.Error("error lua state")
	}
}
