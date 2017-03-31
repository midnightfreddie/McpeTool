package apiserver

import (
	"io"
	"io/ioutil"
	"net/http"
)

// func StartApiServer(path string) {
func StartApiServer() {
	files, err := ioutil.ReadDir(`/storage/emulated/0/games/com.mojang/minecraftWorlds`)
	if err != nil {
		go func() {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "<html><head><title>ApiError</title></head><body><h1 style=\"color: red;\">API Error</h1><div>"+err.Error()+"</div></body></html>")
			})
			http.ListenAndServe(":8080", nil)
		}()
		return
	}
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// files, err := ioutil.ReadDir(path + "/db")
			io.WriteString(w, "<html><head><title>HelloWorld</title></head><body><h1 style=\"color: red;\">Hello world!</h1><ul>")
			for _, file := range files {
				// fmt.Println(file.Name())
				io.WriteString(w, "<li>"+string(file.Name())+"</li>")
			}
			// io.WriteString(w, err.Error())

			io.WriteString(w, "</ul></body></html>")
		})
		http.ListenAndServe(":8080", nil)
	}()
	// path := `/storage/emulated/0/games/com.mojang/minecraftWorldsh4wKANYDAQA\=`
	// world, err := world.OpenWorld(path)
	// if err != nil {
	// 	go func() {
	// 		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 			// files, err := ioutil.ReadDir(path + "/db")
	// 			files, err := ioutil.ReadDir(`/storage/emulated/0/games/com.mojang/minecraftWorlds`)
	// 			if err != nil {
	// 				http.Error(w, err.Error(), 500)
	// 				return
	// 			}
	// 			io.WriteString(w, "<html><head><title>HelloWorld</title></head><body><h1 style=\"color: red;\">Hello world!</h1><ul>")
	// 			for _, file := range files {
	// 				// fmt.Println(file.Name())
	// 				io.WriteString(w, "<li>"+string(file.Name())+"</li>")
	// 			}
	// 			io.WriteString(w, err.Error())

	// 			io.WriteString(w, "</ul></body></html>")
	// 		})
	// 		http.ListenAndServe(":8080", nil)
	// 	}()
	// 	return
	// }
	// defer world.Close()
	// go func() {
	// 	err = api.Server(&world)
	// 	if err != nil {
	// 		return
	// 	}
	// }()
}
