package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"
	"encoding/binary"

	"github.com/midnightfreddie/McpeTool/world"
	"github.com/midnightfreddie/nbt2json"
)

// Level is the default JSON response object
type Level struct {
	ApiVersion string       `json:"apiVersion"`
	FilePath   string       `json:"filePath,omitempty"`
	LevelDat   LevelDatInfo `json:"levelDat,omitempty"`
}

// LevelDatInfo represents the level.dat file
type LevelDatInfo struct {
	Version int32           `json:"version"`
	NBT     json.RawMessage `json:"nbt,omitempty"`
}

// NewLevel initializes and returns a Response object
func NewLevel() *Level {
	return &Level{ApiVersion: apiVersion}
}

func levelApi(world *world.World, path string) {
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

		outData := NewLevel()
		relPath := r.URL.Path[len(path):]
		if relPath != "" {
		}
		switch r.Method {
		case "GET":
			if relPath == "" {
				myLevelDat, err := world.GetLevelDat()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				buf := bytes.NewReader(myLevelDat)
				err = binary.Read(buf, binary.LittleEndian, &outData.LevelDat.Version)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
				outData.LevelDat.NBT, err = nbt2json.Nbt2Json(myLevelDat[8:], binary.LittleEndian)
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			} else {
				http.Error(w, "no URLs under "+path, 404)
				return
			}
		case "HEAD":
			return
		default:
			http.Error(w, "Method "+r.Method+" not supported", 405)
			return
		}
		outJson, err := json.MarshalIndent(outData, "", "  ")
		// outJson, err := json.Marshal(keylist)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintln(w, string(outJson[:]))
	})
}
