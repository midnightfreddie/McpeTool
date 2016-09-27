package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/midnightfreddie/McpeTool/world"
)

type Key struct {
	KeyString string `json:"keyString,omitempty"`
	Base64Key string `json:"base64Key"`
	Key       []int  `json:"key"`
}

// SetKey is used to set the base64 and byte array versions of the key and ensure consistency
func (k *Key) SetKey(key []byte) {
	// json.Marshall will base64-encode byte arrays instead of making a JSON array, so making an array of ints to get desired behavior in JSON output
	k.Key = make([]int, len(key))
	allAscii := true
	for i := range key {
		k.Key[i] = int(key[i])
		if key[i] < 0x20 || key[i] > 0x7e {
			allAscii = false
		}
	}
	if allAscii {
		k.KeyString = string(key[:])
	}
	k.Base64Key = base64.StdEncoding.EncodeToString(key)
}

// KeyList is the structure used for JSON replies to key list requests
type KeyList struct {
	Keys []Key `json:"keys"`
}

// SetKeys is used to populate an array of Keys
func (k *KeyList) SetKeys(inKeyList [][]byte) {
	outKeyList := make([]Key, len(inKeyList))
	for i := 0; i < len(inKeyList); i++ {
		outKeyList[i].SetKey(inKeyList[i])
	}
	k.Keys = append(k.Keys, outKeyList...)
}

// Server is the http REST API server
func Server(world *world.World) error {
	http.HandleFunc("/api/v1/db/", func(w http.ResponseWriter, r *http.Request) {
		keylist, err := world.GetKeys()
		if err != nil {
			panic(err.Error())
		}
		outData := KeyList{}
		outData.SetKeys(keylist)
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
