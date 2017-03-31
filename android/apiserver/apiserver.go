package apiserver

import (
	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
)

// func StartApiServer(path string) {
func StartApiServer() {
	// go func() {
	// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 		// files, err := ioutil.ReadDir(path + "/db")
	// 		// files, err := ioutil.ReadDir(`/storage/emulated/0/games/com.mojang/minecraftWorlds`)
	// 		files, err := ioutil.ReadDir(`/storage/emulated/0/games/com.mojang/minecraftWorlds/3`)
	// 		// files, err := ioutil.ReadDir(`/storage/emulated/0/games/com.mojang/minecraftWorlds/h4wKANYDAQA\=`)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		io.WriteString(w, "<html><head><title>HelloWorld</title></head><body><h1 style=\"color: red;\">Hello world!</h1><ul>")
	// 		for _, file := range files {
	// 			// fmt.Println(file.Name())
	// 			io.WriteString(w, "<li>"+string(file.Name())+"</li>")
	// 		}
	// 		io.WriteString(w, "</ul></body></html>")
	// 	})
	// 	http.ListenAndServe(":8080", nil)
	// }()
	// path := `/storage/emulated/0/games/com.mojang/minecraftWorlds/h4wKANYDAQA\=`
	path := `/storage/emulated/0/games/com.mojang/minecraftWorlds/mod`

	go func() {
		world, err := world.OpenWorld(path)
		if err != nil {
			panic(err)
		}
		defer world.Close()
		err = api.Server(&world)
		if err != nil {
			panic(err)
		}
	}()
}
