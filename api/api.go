package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	return
}

// Server is the http REST API server
func Server(world *world.World) error {
	apiDbPath := "/api/v1/db/"
	http.HandleFunc(apiDbPath, func(w http.ResponseWriter, r *http.Request) {
		var err error
		outData := NewResponse()
		relPath := r.URL.Path[len(apiDbPath):]
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
		case "DELETE":
			if relPath == "" {
				http.Error(w, "Need to provide key to delete", 400)
				return
			}
			err = world.Delete(outData.key)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		case "PUT":
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading body: "+err.Error(), 400)
				return
			}
			inJson := Response{}
			err = json.Unmarshal(body, &inJson)
			if err != nil {
				http.Error(w, "Error parsing body: "+err.Error(), 400)
				return
			}
			data, err := base64.StdEncoding.DecodeString(inJson.Base64Data)
			if err != nil {
				http.Error(w, "Error decoding base64Data: "+err.Error(), 400)
				return
			}
			err = world.Put(outData.key, data)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			// http.Error(w, "Method "+r.Method+" is under development and not yet operational", 405)
			// return
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
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
	return nil
}
