package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

// WorldResponse is the default JSON response object
type WorldResponse struct {
	key        []byte
	keys       [][]byte
	data       []byte
	ApiVersion string `json:"apiVersion"`
	Keys       []Key  `json:"keys,omitempty"`
	StringKey  string `json:"stringKey,omitempty"`
	HexKey     string `json:"hexKey,omitempty"`
	Base64Data string `json:"base64Data,omitempty"`
}

// NewWorldResponse initializes and returns a Response object
func NewWorldResponse() *WorldResponse {
	return &WorldResponse{ApiVersion: apiVersion}
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *WorldResponse) Fill() {
	o.StringKey, o.HexKey, _ = convertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	o.Keys = make([]Key, len(o.keys))
	for i := range o.Keys {
		o.Keys[i].StringKey, o.Keys[i].HexKey, _ = convertKey(o.keys[i])
	}
}

func worldApi(world *world.World, path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var err error
		outData := NewWorldResponse()
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
