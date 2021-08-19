package blua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func dbModule(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("db", lt)
	L.RawSet(lt, lua.LString("get_keys"), L.NewFunction(dbGetKeys))
}

func dbGetKeys(L *lua.LState) int {
	fmt.Println("temp code filler in dbTest")
	keys, err := myWorld.GetKeys()
	if err != nil {
		panic(err)
	}
	// TODO: push byte array to Lua stack
	for _, k := range keys {
		fmt.Println(k)
	}
	return 1
}
