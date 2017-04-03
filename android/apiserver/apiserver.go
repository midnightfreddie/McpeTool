package apiserver

import "github.com/midnightfreddie/McpeTool/api"

func StartApiServer() {
	// path := `/storage/emulated/0/games/com.mojang/minecraftWorlds/h4wKANYDAQA=`
	go func() {
		// world, err := world.OpenWorld(path)
		// if err != nil {
		// 	panic(err)
		// }
		// defer world.Close()
		// err = api.Server(&world)
		// if err != nil {
		// 	panic(err)
		// }
		api.WorldsServer()
	}()
}
