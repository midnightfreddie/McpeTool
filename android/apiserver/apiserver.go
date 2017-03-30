package apiserver

import (
	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
)

// func StartApiServer(path string) {
func StartApiServer() {
	path := `/storage/emulated/0/games/com.mojang/minecraftWorldsh4wKANYDAQA\=`
	world, err := world.OpenWorld(path)
	if err != nil {
		// go func() {
		// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 		// files, err := ioutil.ReadDir(path + "/db")
		// 		files, err := ioutil.ReadDir(`/storage/emulated/0/games/com.mojang/minecraftWorlds`)
		// 		if err != nil {
		// 			http.Error(w, err.Error(), 500)
		// 			return
		// 		}
		// 		io.WriteString(w, "<html><head><title>HelloWorld</title></head><body><h1 style=\"color: red;\">Hello world!</h1><ul>")
		// 		for _, file := range files {
		// 			// fmt.Println(file.Name())
		// 			io.WriteString(w, "<li>"+string(file.Name())+"</li>")
		// 		}
		// 		io.WriteString(w, err.Error())

		// 		io.WriteString(w, "</ul></body></html>")
		// 	})
		// 	http.ListenAndServe(":8080", nil)
		// }()
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
