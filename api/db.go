package api

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"encoding/binary"

	"github.com/midnightfreddie/McpeTool/world"
	"github.com/midnightfreddie/nbt2json"
)

// DbResponse is the default JSON response object
type DbResponse struct {
	key          []byte
	keys         [][]byte
	data         []byte
	ApiVersion   string          `json:"apiVersion"`
	Keys         []Key           `json:"keys,omitempty"`
	StringKey    string          `json:"stringKey,omitempty"`
	HexKey       string          `json:"hexKey,omitempty"`
	Base64Key    string          `json:"base64Key,omitempty"`
	Base64Data   string          `json:"base64Data,omitempty"`
	Nbt2JsonData json.RawMessage `json:"nbt2jsonData,omitempty"`
	// HexDumpData  string          `json:"hexDumpData,omitempty"`
}

// NewDbResponse initializes and returns a Response object
func NewDbResponse() *DbResponse {
	return &DbResponse{ApiVersion: apiVersion}
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *DbResponse) Fill(urlPrefix string) {
	o.StringKey, o.HexKey, o.Base64Key = convertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	// Not checking error...if it works, field is populated. If not, field is nil. That works.
	o.Nbt2JsonData, _ = nbt2json.Nbt2Json(o.data, binary.LittleEndian)
	// o.HexDumpData = hex.Dump(o.data)
	o.Keys = make([]Key, len(o.keys))
	for i := range o.Keys {
		// o.Keys[i] = world.KeyInfo(o.keys[i])
		o.Keys[i].StringKey, o.Keys[i].HexKey, o.Keys[i].Base64Key = convertKey(o.keys[i])
		if urlPrefix != "" && o.Keys[i].HexKey != "" {
			o.Keys[i].Url = urlPrefix + o.Keys[i].HexKey
		}
	}
}

// Key is the element type in the Response.Keys array
type Key struct {
	StringKey string `json:"stringKey,omitempty"`
	HexKey    string `json:"hexKey"`
	Base64Key string `json:"base64Key,omitempty"`
	Url       string `json:"url,omitempty"`
}

func dbApi(world *world.World, path string) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		var err error

		// Set Origin headers for CORS
		// yoinked from http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang Matt Bucci's answer
		// Could/should go in a Handle not HandleFunc, but I'm not yet quite sure how to do that with the default mux
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		outData := NewDbResponse()
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
			inJson := DbResponse{}
			err = json.Unmarshal(body, &inJson)
			if err != nil {
				http.Error(w, "Error parsing body: "+err.Error(), 400)
				return
			}
			var data []byte
			if len(inJson.Nbt2JsonData) > 0 {
				data, err = nbt2json.Json2Nbt(inJson.Nbt2JsonData, binary.LittleEndian)
				if err != nil {
					http.Error(w, "Error decoding nbt2jsonData: "+err.Error(), 400)
					return
				}
			} else {
				data, err = base64.StdEncoding.DecodeString(inJson.Base64Data)
				if err != nil {
					http.Error(w, "Error decoding base64Data: "+err.Error(), 400)
					return
				}
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
		// TODO: URL prefix should be a variable and configurable. Or perhaps pulled from server.
		outData.Fill("http://127.0.0.1:8080" + r.URL.Path[:len(path)])
		outJson, err := json.MarshalIndent(outData, "", "  ")
		// outJson, err := json.Marshal(keylist)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintln(w, string(outJson[:]))
	})
}
