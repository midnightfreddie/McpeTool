package blua

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Blua injects functions into a gopher-lua state
func Blua(L *lua.LState) error {
	fmt.Println("blua doesn't do anything yet")
	return nil
}
