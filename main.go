package main

// Note: In my $GOPATH/src I have github.com/midnightfreddie/goleveldb/leveldb (addzlib branch) in place of github.com/syndtr/goleveldb/leveldb
//   This adds zlib decompression to the reader as compression type 2 which is needed to read MCPE ldb files
import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

// Note: Open the .mcworld file as a zip--rename it to .mcworld.zip if needed--then copy the db folder
//  to the folder where you'll be running this program

func main() {
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic("error")
	}
	defer db.Close()

	player, err := db.Get([]byte("~local_player"), nil)
	if err != nil {
		panic("error")
	}
	fmt.Println(string(player[:]))

	// iterate and print the first 10 key/value pairs
	iter := db.NewIterator(nil, nil)
	for i := 1; i < 1; iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(key)
		fmt.Println(value)
		i++
	}
	iter.Release()
	err = iter.Error()
	fmt.Println(err)
}

// http://minecraft.gamepedia.com/Pocket_Edition_level_format
