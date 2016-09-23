package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic("error")
	}
	defer db.Close()

	iter := db.NewIterator(nil, nil)

	for i := 1; i < 1000; iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Println(key)
		fmt.Println(value)
		i++
	}
	iter.Release()
	err = iter.Error()
}
