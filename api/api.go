package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/midnightfreddie/goleveldb/leveldb"
)

// Server is the http REST API server
func Server(db *leveldb.DB) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		iter := db.NewIterator(nil, nil)
		for iter.Next() {
			key := iter.Key()
			switch {
			case len(key) == 9:
				switch key[8] {
				case 0x30, 0x31, 0x32, 0x76:
					fmt.Fprintln(w, key)
				default:
					fmt.Fprintln(w, string(key[:]))
				}
			case len(key) == 13:
				switch key[12] {
				case 0x30, 0x31, 0x32, 0x76:
					fmt.Fprintln(w, key)
				default:
					fmt.Fprintln(w, string(key[:]))
				}
			default:
				fmt.Fprintln(w, string(key[:]))
			}
		}
		iter.Release()
		err := iter.Error()
		if err != nil {
			panic(err.Error())
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}
