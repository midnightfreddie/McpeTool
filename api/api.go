package api

import (
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/midnightfreddie/goleveldb/leveldb"
)

// Server is the http REST API server
func Server(db *leveldb.DB) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	return nil
}
