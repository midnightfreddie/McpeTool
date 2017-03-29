package androidtest

import (
	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
)

func StartApiServer(path string) {
	world, err := world.OpenWorld(path)
	if err != nil {
		return
	}
	defer world.Close()
	go func() {
		err = api.Server(&world)
		if err != nil {
			return
		}
	}()
}
