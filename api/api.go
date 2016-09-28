package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

var apiVersion = "0.0"

// Response is the default JSON response object
type Response struct {
	key        []byte
	keys       [][]byte
	data       []byte
	ApiVersion string `json:"apiVersion"`
	Context    string `json:"context,omitempty"`
	Keys       []Key  `json:"keys,omitempty"`
	KeyString  string `json:"keyString,omitempty"`
	HexKey     string `json:"hexKey,omitempty"`
	Base64Data string `json:"base64Data,omitempty"`
}

// NewResponse initializes and returns a Response object
func NewResponse() *Response {
	return &Response{ApiVersion: apiVersion}
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *Response) Fill() {
	o.KeyString, o.HexKey = convertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	o.Keys = make([]Key, len(o.keys))
	for i := range o.Keys {
		o.Keys[i].KeyString, o.Keys[i].HexKey = convertKey(o.keys[i])
	}
}

// Key is the element type in the Response.Keys array
type Key struct {
	KeyString string `json:"keyString,omitempty"`
	HexKey    string `json:"hexKey"`
}

// convertKey takes a byte array and returns a string if all characters are printable (else "")  hex-string-encoded versions of key
func convertKey(k []byte) (keyString, hexKey string) {
	allAscii := true
	for i := range k {
		if k[i] < 0x20 || k[i] > 0x7e {
			allAscii = false
		}
	}
	if allAscii {
		keyString = string(k[:])
	}
	hexKey = hex.EncodeToString(k)
	return
}

// Server is the http REST API server
func Server(world *world.World) error {
	http.HandleFunc("/api/v1/db/", func(w http.ResponseWriter, r *http.Request) {
		var err error
		outData := NewResponse()
		outData.keys, err = world.GetKeys()
		if err != nil {
			panic(err.Error())
		}
		outData.Fill()
		outJson, err := json.MarshalIndent(outData, "", "  ")
		// outJson, err := json.Marshal(keylist)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintln(w, string(outJson[:]))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}
