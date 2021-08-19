package blua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func dbModule(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("db", lt)
	L.RawSet(lt, lua.LString("test"), L.NewFunction(dbTest))
}

func dbTest(L *lua.LState) int {
	fmt.Println("temp code filler in dbTest")
	return 1
}
