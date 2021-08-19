package blua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func dbModule(L *lua.LState) {
	fmt.Println("dbModule")
	lt := L.NewTable()
	lt.Append(lua.LString("test string"))
	L.SetGlobal("db", lt)
	L.RawSet(lt, lua.LString("get_keys"), L.NewFunction(dbGetKeys))
	L.RawSet(lt, lua.LString("keys"), L.NewTable())
}

func dbGetKeys(L *lua.LState) int {
	fmt.Println("temp code filler in dbTest")
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
			// fmt.Println(k)
			for _, kk := range k {
				kkt.Append(lua.LNumber(kk))
			}
		}
	}
	return 0
}
