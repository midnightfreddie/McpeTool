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
	Base64Key string `json:"base64Key"`
	Key       []int  `json:"key"`
}

// SetKey is used to set the base64 and byte array versions of the key and ensure consistency
func (k *Key) SetKey(key []byte) {
	k.Key = make([]int, len(key))
	for i := range key {
		k.Key[i] = int(key[i])
	}
	k.Base64Key = base64.StdEncoding.EncodeToString(key)
}

// KeyList is the structure used for JSON replies to key list requests
type KeyList struct {
	KeyList []Key `json:"keyList"`
}

// SetKeys is used to populate an array of Keys
func (k *KeyList) SetKeys(inKeyList [][]byte) {
	outKeyList := make([]Key, len(inKeyList))
	for i := 0; i < len(inKeyList); i++ {
		outKeyList[i].SetKey(inKeyList[i])
	}
	k.KeyList = append(k.KeyList, outKeyList...)
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
