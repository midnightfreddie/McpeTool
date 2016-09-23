package main

import (
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, _ := leveldb.OpenFile("db/test.ldb", nil)
	defer db.Close()
}
