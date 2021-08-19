package blua

import (
	lua "github.com/yuin/gopher-lua"
)

func worldModule(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("world", lt)
	L.RawSet(lt, lua.LString("load"), L.NewFunction(openWorld))

}
