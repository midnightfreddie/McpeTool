package main

// Note: In my $GOPATH/src I have github.com/midnightfreddie/goleveldb/leveldb (addzlib branch) in place of github.com/syndtr/goleveldb/leveldb
//   This adds zlib decompression to the reader as compression type 2 which is needed to read MCPE ldb files
import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/quag/mcobj/nbt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/urfave/cli"
)

func proofOfConcept() {
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
	id, out, err := mynbt.ReadTag()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("\n\n")
	fmt.Printf("\n%d%s\n", id, out)

	// iterate and print the first 10 key/value pairs
	iter := db.NewIterator(nil, nil)
	for i := 1; i < 10; iter.Next() {
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

func main() {
	app := cli.NewApp()
	app.Name = "MCPE Tool"
	app.Version = "0.0.0"
	app.Usage = "A utility to access Minecraft Portable Edition .mcworld exported world files."

	app.Commands = []cli.Command{
		{
			Name:    "keys",
			Aliases: []string{"k"},
			Usage:   "Lists all keys in the database. Be sure to include the path to the db, e.g. 'McpeTool keys db'",
			Action: func(c *cli.Context) error {
				o := &opt.Options{
					ReadOnly: true,
				}
				db, err := leveldb.OpenFile(c.Args().First(), o)
				if err != nil {
					panic("error")
				}
				defer db.Close()

				iter := db.NewIterator(nil, nil)
				for iter.Next() {
					key := iter.Key()
					switch {
					case len(key) == 9:
						switch key[8] {
						case 0x30, 0x31, 0x32, 0x76:
							fmt.Println(key)
						default:
							fmt.Println(string(key[:]))
						}
					case len(key) == 13:
						fmt.Println(key)
					default:
						fmt.Println(string(key[:]))
					}
				}
				iter.Release()
				err = iter.Error()
				if err != nil {
					panic(err.Error())
				}
				return nil
			},
		},
		{
			Name:    "proofofconcept",
			Aliases: []string{"poc"},
			Usage:   "Run the original POC code which assumes a folder \"db\" is present with the *.ldb and other level files",
			Action: func(c *cli.Context) error {
				proofOfConcept()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
