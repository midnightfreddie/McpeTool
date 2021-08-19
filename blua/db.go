package blua

import (
	lua "github.com/yuin/gopher-lua"
)

func dbModule(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("db", lt)
	L.RawSet(lt, lua.LString("get_keys"), L.NewFunction(dbGetKeys))
	// Unsure if I need to define an empty key set
	// L.RawSet(lt, lua.LString("keys"), L.NewTable())
}

func dbGetKeys(L *lua.LState) int {
	keys, err := myWorld.GetKeys()
	if err != nil {
		panic(err)
	}
	lt := L.GetGlobal("db")
	if ltbl, ok := lt.(*lua.LTable); ok {
		kt := L.NewTable()
		L.RawSet(ltbl, lua.LString("keys"), kt)
		for _, k := range keys {
			kkt := L.NewTable()
			kt.Append(kkt)
			for _, kk := range k {
				kkt.Append(lua.LNumber(kk))
			}
		}
	}
	return 0
}
