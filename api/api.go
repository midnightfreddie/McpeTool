package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

// Server is the http REST API server
func Server(world *world.World) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		keylist, err := world.GetKeys()
		if err != nil {
			panic(err.Error())
		}
		// outJson, err := json.MarshalIndent(keylist, "", "  ")
		outJson, err := json.Marshal(keylist)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintln(w, string(outJson[:]))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}
