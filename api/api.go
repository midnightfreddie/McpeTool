package api

import (
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

var apiVersion = "1.0"

// Response is the default JSON response object
type Response struct {
	key        []byte
	keys       [][]byte
	data       []byte
	ApiVersion string `json:"apiVersion"`
	Keys       []Key  `json:"keys,omitempty"`
	StringKey  string `json:"stringKey,omitempty"`
	HexKey     string `json:"hexKey,omitempty"`
	Base64Data string `json:"base64Data,omitempty"`
}

// NewResponse initializes and returns a Response object
func NewResponse() *Response {
	return &Response{ApiVersion: apiVersion}
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *Response) Fill() {
	o.StringKey, o.HexKey = ConvertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	o.Keys = make([]Key, len(o.keys))
	for i := range o.Keys {
		o.Keys[i].StringKey, o.Keys[i].HexKey = ConvertKey(o.keys[i])
	}
}

// Key is the element type in the Response.Keys array
type Key struct {
	StringKey string `json:"stringKey,omitempty"`
	HexKey    string `json:"hexKey"`
}

// ConvertKey takes a byte array and returns a string if all characters are printable (else "")  hex-string-encoded versions of key
func ConvertKey(k []byte) (stringKey, hexKey string) {
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
