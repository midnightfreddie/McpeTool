package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/midnightfreddie/goleveldb/leveldb"
)

// Server is the http REST API server
func Server(db *leveldb.DB) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		keylist := [][]byte{}
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			key := iter.Key()
			tmp := make([]byte, len(key))
			copy(tmp, key)
			keylist = append(keylist, tmp)
			// outJson, err := json.Marshal(key)
			// json.NewEncoder(w).Encode(key)
			// if err != nil {
			// 	panic(err.Error())
			// }
			// fmt.Println(string(outJson[:]))
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			panic(err.Error())
		}
		// fmt.Fprintf(w, "%v\n\n", keylist)
		outJson, err := json.MarshalIndent(keylist, "", "  ")
		// outJson, err := json.Marshal(keylist)
		if err != nil {
			panic(err.Error())
		}
		fmt.Fprintln(w, string(outJson[:]))
		// json.NewEncoder(w).Encode(keylist)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}
