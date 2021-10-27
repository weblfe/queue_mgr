package plugins

import (
	"errors"
	"github.com/weblfe/queue_mgr/plugins/libs"
	"github.com/yuin/gopher-lua"
	"io"
	"runtime"
		"sort"
		"sync"
	"time"
)

type (
	PluginBootLoader func() (*PluginOptions, error)

	LuaState struct {
		lua.LState
	}

	luaPluginImpl struct {
		lvm         *LuaState
		constructor *sync.Once
		bootAt      time.Time
		loader      PluginBootLoader
		options     *PluginOptions
		extLibs     []*LuaRegistryFunction
	}

	PluginOptions struct {
		Extends []*LuaRegistryFunction
		lua.Options
	}

	LuaRegistryFunction struct {
		LName     string
		LFunction lua.LGFunction
	}
)

func NewLuaPlugin(options ...PluginOptions) *luaPluginImpl {
	var plugin = new(luaPluginImpl)
	if len(options) > 0 {
		plugin.options = &options[0]
	} else {
		plugin.options = NewDefaultOptions()
	}
	return plugin.init()
}

func NewDefaultOptions() *PluginOptions {
	var opt = new(PluginOptions)
	opt.SkipOpenLibs = false

	return opt
}

func NewLuaState(options ...lua.Options) *LuaState {
	var (
		state    = new(LuaState)
		luaState = lua.NewState(options...)
	)
	state.LState = *luaState
	return state
}

func (options *PluginOptions) GetLuaOptions() lua.Options {
	var opt = options.Options
	return opt
}

func (plugin *luaPluginImpl) init() *luaPluginImpl {
	if plugin == nil {
		return nil
	}
	if plugin.lvm == nil {
		if plugin.options == nil {
			plugin.lvm = NewLuaState()
		} else {
			plugin.lvm = NewLuaState(plugin.options.GetLuaOptions())
		}
	}
	if plugin.constructor == nil {
		plugin.constructor = &sync.Once{}
	}
	runtime.SetFinalizer(plugin, (*luaPluginImpl).destroy)
	return plugin
}

func (plugin *luaPluginImpl) SetLoader(loader PluginBootLoader) *luaPluginImpl {
	if loader == nil || !plugin.bootAt.IsZero() {
		return plugin
	}
	plugin.loader = loader
	return plugin
}

func (plugin *luaPluginImpl) Boot() {
	if plugin == nil {
		return
	}
	plugin.constructor.Do(func() {
		plugin.initLoader().loads()
		plugin.bootAt = time.Now()
	})
}

func (plugin *luaPluginImpl) initLoader() *luaPluginImpl {
	if plugin.loader == nil {
		return plugin
	}
	if opts, err := plugin.loader(); err == nil && opts != nil {
		plugin.register(opts.Extends)
	}
	return plugin
}

func (plugin *luaPluginImpl) register(libs []*LuaRegistryFunction) {
	if len(libs) <= 0 || !plugin.bootAt.IsZero() {
		return
	}
	plugin.extLibs = append(plugin.extLibs, libs...)
}

func (plugin *luaPluginImpl) loads() {
	var vm = plugin.GetVM()
	if plugin.extLibs == nil || vm == nil {
		return
	}
	var (
		libRegistries  []LuaRegistryFunction
		cache = make(map[string]bool)
	)
	for _, lib := range plugin.extLibs {
		if lib == nil || lib.LName == "" || lib.LFunction == nil {
			continue
		}
		if _, ok := cache[lib.LName]; ok {
			continue
		}
			libRegistries = append(libRegistries, *lib)
	}

	if len(libRegistries) > 0 {
		plugin.extend(vm, libRegistries)
	}
}

func (plugin *luaPluginImpl) extend(state *LuaState, extends []LuaRegistryFunction) {
	if len(extends) <= 0 || state == nil {
		return
	}
	for _, lib := range extends {
		state.Push(state.NewFunction(lib.LFunction))
		state.Push(lua.LString(lib.LName))
		state.Call(1, 0)
	}
}

func (plugin *luaPluginImpl) GetVM() *LuaState {
	return plugin.lvm
}

func (plugin *luaPluginImpl) GetLState() *lua.LState {
	var state = &plugin.lvm.LState
	return state
}

func (plugin *luaPluginImpl) Eval(data []byte) error {
	return plugin.GetVM().DoString(string(data))
}

func (plugin *luaPluginImpl) EvalExpr(luaExpr string) error {
	return plugin.GetVM().DoString(luaExpr)
}

func (plugin *luaPluginImpl) LoadFile(file string) (*lua.LFunction, error) {
	return plugin.GetVM().LoadFile(file)
}

func (plugin *luaPluginImpl) Libs() []string {
	var (
		libArr []string
		global = plugin.GetLState().G
	)
	if global == nil || global.Global == nil {
		return libArr
	}
	global.Global.ForEach(func(value lua.LValue, value2 lua.LValue) {
		if value.Type() == lua.LTNil || value2.Type() == lua.LTNil {
			return
		}
		libArr = append(libArr, value.String())
	})
	if len(libArr) >0 {
			sort.Strings(libArr)
	}
	return libArr
}

func (plugin *luaPluginImpl) LoadByIo(reader io.ReadCloser, name string) (*lua.LFunction, error) {
	if reader == nil {
		return nil, errors.New("reader nil")
	}
	defer func() {
		_ = reader.Close()
	}()
	return plugin.GetVM().Load(reader, name)
}

func (plugin *luaPluginImpl) destroy() {
	runtime.SetFinalizer(plugin, nil)
	if plugin.lvm != nil {
		plugin.lvm.Close()
		plugin.lvm = nil
	}
	plugin.extLibs = nil
	plugin.constructor = nil
}

func CreateExtendsLoader() (*PluginOptions, error) {
	var opts = &PluginOptions{
		Extends: []*LuaRegistryFunction{
			{
				LName:     libs.LoggerModule,
				LFunction: libs.NewLuaLoggerTables(),
			},
		},
	}
	return opts, nil
}
