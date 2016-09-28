package api

import (
	"encoding/base64"
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
	Base64Key  string `json:"base64Key,omitempty"`
	Key        []int  `json:"key,omitempty"`
	Base64Data string `json:"base64Data,omitempty"`
}

func NewResponse() *Response {
	return &Response{ApiVersion: apiVersion}
}

// Fill is used to conver the raw data to JSON-friendly data before returning to client
func (o *Response) Fill() {
	o.KeyString, o.Base64Key, o.Key = convertKey(o.key)
	o.Keys = make([]Key, len(o.keys))
	for i := range o.Keys {
		o.Keys[i].KeyString, o.Keys[i].Base64Key, o.Keys[i].Key = convertKey(o.keys[i])
	}
}

type Key struct {
	key       []byte
	KeyString string `json:"keyString,omitempty"`
	Base64Key string `json:"base64Key"`
	Key       []int  `json:"key"`
}

func convertKey(k []byte) (keyString, base64Key string, intArray []int) {
	// json.Marshall will base64-encode byte arrays instead of making a JSON array, so making an array of ints to get desired behavior in JSON output
	intArray = make([]int, len(k))
	allAscii := true
	for i := range k {
		intArray[i] = int(k[i])
		if k[i] < 0x20 || k[i] > 0x7e {
			allAscii = false
		}
	}
	if allAscii {
		keyString = string(k[:])
	}
	base64Key = base64.StdEncoding.EncodeToString(k)
	return
}

// Fill is used to set the base64 and int array versions of the key
func (k *Key) Fill() {
	k.KeyString, k.Base64Key, k.Key = convertKey(k.key)
}

// // KeyList is the structure used for JSON replies to key list requests
// type KeyList struct {
// 	Keys []Key `json:"keys"`
// }

// // SetKeys is used to populate an array of Keys
// func (k *KeyList) SetKeys(inKeyList [][]byte) {
// 	outKeyList := make([]Key, len(inKeyList))
// 	for i := 0; i < len(inKeyList); i++ {
// 		outKeyList[i].SetKey(inKeyList[i])
// 	}
// 	k.Keys = append(k.Keys, outKeyList...)
// }

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
		// outData := KeyList{}
		// outData.SetKeys(keylist)

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
