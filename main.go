package main

// Note: In my $GOPATH/src I have github.com/midnightfreddie/goleveldb/leveldb (addzlib branch) in place of github.com/syndtr/goleveldb/leveldb
//   This adds zlib decompression to the reader as compression type 2 which is needed to read MCPE ldb files
import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/quag/mcobj/nbt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Note: Open the .mcworld file as a zip--rename it to .mcworld.zip if needed--then copy the db folder
//  to the folder where you'll be running this program

func main() {
	// db, err := leveldb.OpenFile("db", nil)
	// Setting readOnly to true
	//   now thinking I can read directly from the zip file, perhaps
	o := &opt.Options{
		ReadOnly: true,
	}
	db, err := leveldb.OpenFile("db", o)
	if err != nil {
		panic("error")
	}
	defer db.Close()

	player, err := db.Get([]byte("~local_player"), nil)
	if err != nil {
		panic("error")
	}
	fmt.Println(hex.Dump(player[:]))
	nbtr := bytes.NewReader(player)
	mynbt := nbt.NewReader(nbtr)
	// out, _ := nbt.Parse(nbtr)
	// fmt.Println(json.Marshal(out))
	// out, _ := mynbt.ReadStruct()
	id, out, err := mynbt.ReadTag()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\n\n")
	fmt.Printf("\n%d%s\n", id, out)
	id, out, err = mynbt.ReadTag()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\n\n")
	fmt.Printf("\n%d%s\n", id, out)

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
	if err != nil {
		panic(err.Error())
	}
}

// http://minecraft.gamepedia.com/Pocket_Edition_level_format
