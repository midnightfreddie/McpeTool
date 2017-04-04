package api

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

var apiVersion = "1.0"

// type apiConfig struct {
// 	worldsPath    string
// 	h             func(config apiConfig, w http.ResponseWriter, r *http.Request)
// 	urlPathPrefix string
// 	relPath       string
// 	serverurl     string
// }

// TODO: this moved to world/keys.go; remove its use from this package and delete
// ACTUALLY: No. This package shouldn't rely on world/keys.go for its formatting
// convertKey takes a byte array and returns a string if all characters are printable (else "")  hex-string-encoded versions of key
func convertKey(k []byte) (stringKey, hexKey, base64Key string) {
	allAscii := true
	for i := range k {
		if k[i] < 0x20 || k[i] > 0x7e {
			allAscii = false
		}
	}
	if allAscii {
		stringKey = string(k[:])
	}
	hexKey = hex.EncodeToString(k)
	base64Key = base64.StdEncoding.EncodeToString(k)
	return
}

// Server is the http REST API server
func Server(world *world.World) error {

	// http handler functions defined in other files in this package
	// dbApi(world, "/api/v1/db/")
	worldApi(world, "/api/v1/world/")
	playerApi(world, "/api/v1/player/")
	// levelApi(world, "/api/v1/level/")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
	return nil
}

// WorldsServer is the REST API server that lists and opens worlds on demand
func WorldsServer(worldsPath string) error {
	http.HandleFunc("/api/v1/worlds/", worldsApi(worldsPath, "/api/v1/worlds/"))
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
	return nil
}
