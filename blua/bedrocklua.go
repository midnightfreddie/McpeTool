package blua

import (
	"github.com/midnightfreddie/McpeTool/world"
	lua "github.com/yuin/gopher-lua"
)

var myWorld world.World

// Blua injects functions into a gopher-lua state
func Blua(L *lua.LState) error {
	L.SetGlobal("open_world", L.NewFunction(openWorld))
	dbModule(L)
	return nil
}

func openWorld(L *lua.LState) int {
	var err error
	// if myWorld is initialized
	if myWorld != (world.World{}) {
		myWorld.Close()
	}
	// TODO: Make path a lua parameter
	path := `C:\Users\jim\AppData\Local\Packages\Microsoft.MinecraftUWP_8wekyb3d8bbwe\LocalState\games\com.mojang\minecraftWorlds\fEgcYZ-hBwA=`
	myWorld, err = world.OpenWorld(path)
	if err != nil {
		panic(err)
	}
	return 0
}
