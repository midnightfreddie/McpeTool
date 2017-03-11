package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

// PlayerResponse is the default JSON response object
type PlayerResponse struct {
	key        []byte
	keys       [][]byte
	data       []byte
	ApiVersion string   `json:"apiVersion"`
	Players    []Player `json:"players,omitempty"`
	StringKey  string   `json:"stringKey,omitempty"`
	HexKey     string   `json:"hexKey,omitempty"`
	Base64Data string   `json:"base64Data,omitempty"`
}

// NewPlayerResponse initializes and returns a Response object
func NewPlayerResponse() *PlayerResponse {
	return &PlayerResponse{ApiVersion: apiVersion}
}

// Player is the element type in the Response.Keys array
type Player struct {
	StringKey string `json:"stringKey,omitempty"`
	HexKey    string `json:"hexKey"`
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *PlayerResponse) Fill() {
	o.StringKey, o.HexKey = convertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	o.Players = make([]Player, len(o.keys))
	for i := range o.Players {
		o.Players[i].StringKey, o.Players[i].HexKey = convertKey(o.keys[i])
	}
}

func playerApi(world *world.World, path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var err error
		outData := NewPlayerResponse()
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
				localPlayerKey, _ := hex.DecodeString("7e6c6f63616c5f706c61796572")
				_, err = world.Get(localPlayerKey)
				if err != nil {
					http.Error(w, "~local_player not found", 404)
					return
				}
				outData.keys = make([][]byte, 1)
				outData.keys[0] = localPlayerKey
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
