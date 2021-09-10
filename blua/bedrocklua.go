package blua

import (
	"os"

	"github.com/midnightfreddie/McpeTool/world"
	lua "github.com/yuin/gopher-lua"
)

var myWorld world.World

// Store the path of the world in $MCPETOOL_WORLD
var worldPath = os.Getenv("MCPETOOL_WORLD")

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
	myWorld, err = world.OpenWorld(worldPath)
	if err != nil {
		panic(err)
	}
	return 0
}
