package libs

import lua "github.com/yuin/gopher-lua"

type (
	Loader func(L *lua.LState) int
)
