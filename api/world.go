package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

func worldApi(world *world.World, path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var err error
		outData := NewResponse()
		relPath := r.URL.Path[len(path):]
		if relPath != "" {
			outData.key, err = hex.DecodeString(relPath)
			if err != nil {
				http.Error(w, err.Error()+"\n"+relPath+": URL key must be a byte array coded in hex digits, two digits per byte", 400)
				return
			}
		}
		switch r.Method {
		case "GET":
			if relPath == "" {
				outData.keys, err = world.GetKeys()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			} else {
				outData.data, err = world.Get(outData.key)
				if err != nil {
					if err.Error() == "leveldb: not found" {
						http.Error(w, "key not found", 404)
						return
					}
					http.Error(w, err.Error(), 500)
					return
				}
			}
		case "HEAD":
			return
		default:
			http.Error(w, "Method "+r.Method+" not supported", 405)
			return
		}
		outData.Fill()
		outJson, err := json.MarshalIndent(outData, "", "  ")
		// outJson, err := json.Marshal(keylist)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintln(w, string(outJson[:]))
	})
}
