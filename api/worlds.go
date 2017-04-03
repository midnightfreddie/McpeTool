package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// WorldsResponse is the default JSON response object
type WorldsResponse struct {
	key        []byte
	keys       [][]byte
	data       []byte
	temp       []string
	ApiVersion string   `json:"apiVersion"`
	Keys       []Key    `json:"keys,omitempty"`
	StringKey  string   `json:"stringKey,omitempty"`
	HexKey     string   `json:"hexKey,omitempty"`
	Base64Data string   `json:"base64Data,omitempty"`
	Temp       []string `json:temp`
}

// NewWorldsResponse initializes and returns a Response object
func NewWorldsResponse() *WorldsResponse {
	return &WorldsResponse{ApiVersion: apiVersion}
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *WorldsResponse) Fill(urlPrefix string) {
	o.StringKey, o.HexKey, _ = convertKey(o.key)
	o.Base64Data = base64.StdEncoding.EncodeToString(o.data)
	o.Keys = make([]Key, len(o.keys))
	for i := range o.Keys {
		o.Keys[i].StringKey, o.Keys[i].HexKey, _ = convertKey(o.keys[i])
	}
	o.Temp = make([]string, len(o.temp))
	for i := range o.Temp {
		// TODO: Error handling?
		urlEncoded, _ := url.Parse(urlPrefix)
		urlEncoded.Path += o.temp[i]
		o.Temp[i] = urlEncoded.String()
	}
}

func worldsApi(worldsFilePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set Origin headers for CORS
		// yoinked from http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang Matt Bucci's answer
		// Could/should go in a Handle not HandleFunc, but I'm not yet quite sure how to do that with the default mux
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		outData := NewWorldsResponse()
		// FIXME: hackish and not robust
		path := "/api/v1/worlds/"
		relPath := r.URL.Path[len(path):]
		// if relPath != "" {
		// 	outData.key, err = hex.DecodeString(relPath)
		// 	if err != nil {
		// 		http.Error(w, err.Error()+"\n"+relPath+": URL key must be a byte array coded in hex digits, two digits per byte", 400)
		// 		return
		// 	}
		// }
		switch r.Method {
		case "GET":
			if relPath == "" {
				// worldsFilePath := `/storage/emulated/0/games/com.mojang/minecraftWorlds`
				dirs, err := ioutil.ReadDir(worldsFilePath)
				if err != nil {
					http.Error(w, "Error while reading minecraftWorlds folder: "+err.Error(), 500)
					return
				}
				outData.temp = make([]string, len(dirs))
				for i, dir := range dirs {
					outData.temp[i] = dir.Name()
				}
			}
			//  else {
			// 	outData.data, err = world.Get(outData.key)
			// 	if err != nil {
			// 		if err.Error() == "leveldb: not found" {
			// 			http.Error(w, "key not found", 404)
			// 			return
			// 		}
			// 		http.Error(w, err.Error(), 500)
			// 		return
			// 	}
			// }
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
	}
}
