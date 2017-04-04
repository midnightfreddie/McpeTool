package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type World struct {
	Name        string `json:"name"`
	DirName     string `json:"dirName"`
	FqdnDirName string `json:"fqdnDirName"`
	Url         string `json:"url,omitempty"`
	DbUrl       string `json:"dbUrl,omitempty"`
	LevelUrl    string `json:"levelUrl,omitempty"`
	Error       string `json:"error,omitempty"`
}

func WorldInfo(worldPath, urlPrefix string) *World {
	splitPath := strings.Split(worldPath, "/")
	output := World{
		DirName:     splitPath[len(splitPath)-1],
		FqdnDirName: worldPath,
		Url:         urlPrefix + "/",
		DbUrl:       urlPrefix + "/db/",
		LevelUrl:    urlPrefix + "/level/",
	}
	// world, err := world.OpenWorld(worldPath)
	// if err != nil {
	// 	output.Error = "Opening world: " + err.Error()
	// 	return &output
	// }
	// defer world.Close()
	name, err := ioutil.ReadFile(output.FqdnDirName + `/levelname.txt`)
	if err != nil {
		output.Name = output.DirName
	} else {
		output.Name = string(name[:])
	}

	return &output
}

// WorldsResponse is the default JSON response object
type WorldsResponse struct {
	worldDirs  []string
	ApiVersion string  `json:"apiVersion"`
	Worlds     []World `json:"worlds,omitempty"`
	World      World   `json:"world,omitempty"`
}

// NewWorldsResponse initializes and returns a Response object
func NewWorldsResponse() *WorldsResponse {
	return &WorldsResponse{ApiVersion: apiVersion}
}

// Fill is used to convert the raw byte arrays to JSON-friendly data before returning to client
func (o *WorldsResponse) Fill(worldPath, urlPrefix string) {
	o.Worlds = make([]World, len(o.worldDirs))
	for i := range o.Worlds {
		// TODO: Error handling?
		urlEncoded, _ := url.Parse(urlPrefix)
		urlEncoded.Path += o.worldDirs[i]
		// o.Worlds[i] = urlEncoded.String()
		o.Worlds[i] = *WorldInfo(worldPath+"/"+o.worldDirs[i], urlEncoded.String())
	}
}

func worldsApi(worldsFilePath, path string) http.HandlerFunc {
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
		relPath := r.URL.Path[len(path):]
		if relPath != "" {
			urlencodedDir := strings.Split(relPath, "/")[0]
			worldDir, err := url.QueryUnescape(urlencodedDir)
			if err != nil {
				http.Error(w, "Error decoding url: "+err.Error(), 404)
				return
			}
			outData.World = *WorldInfo(worldsFilePath+"/"+worldDir, r.URL.Path[:len(path)]+urlencodedDir)
		}
		switch r.Method {
		case "GET":
			if relPath == "" {
				// worldsFilePath := `/storage/emulated/0/games/com.mojang/minecraftWorlds`
				dirs, err := ioutil.ReadDir(worldsFilePath)
				if err != nil {
					http.Error(w, "Error while reading minecraftWorlds folder: "+err.Error(), 500)
					return
				}
				outData.worldDirs = make([]string, len(dirs))
				for i, dir := range dirs {
					outData.worldDirs[i] = dir.Name()
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
		outData.Fill(worldsFilePath, r.URL.Path[:len(path)])
		outJson, err := json.MarshalIndent(outData, "", "  ")
		// outJson, err := json.Marshal(keylist)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Fprintln(w, string(outJson[:]))
	}
}
