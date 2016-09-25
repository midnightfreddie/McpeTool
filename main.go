package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/midnightfreddie/goleveldb/leveldb"
	"github.com/midnightfreddie/goleveldb/leveldb/opt"
	"github.com/quag/mcobj/nbt"
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
						switch key[12] {
						case 0x30, 0x31, 0x32, 0x76:
							fmt.Println(key)
						default:
							fmt.Println(string(key[:]))
						}
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
			Name:    "develop",
			Aliases: []string{"dev"},
			Usage:   "Random thing the dev is working on",
			Action: func(c *cli.Context) error {
				db, err := leveldb.OpenFile(c.Args().First(), nil)
				if err != nil {
					panic("error")
				}
				defer db.Close()

				// iter := db.NewIterator(nil, nil)
				// for iter.Next() {
				// 	key := iter.Key()
				// 	if len(key) == 9 && key[8] == 0x30 {
				// 		chunk := iter.Value()
				// 		for i := 0; i < 256; i++ {
				// 			cx := i / 16
				// 			y := 127
				// 			cz := i % 16
				// 			idx := 2048*cx + y + 128*cz
				// 			if (i%16)%2 == 0 {
				// 				chunk[idx] = 20
				// 			}
				// 			fmt.Printf("%d ", chunk[idx])
				// 		}
				// 		fmt.Printf("\n\n")
				// 	}
				// }
				// iter.Release()
				// err = iter.Error()
				key := [9]byte{1, 0, 0, 0, 1, 0, 0, 0, 0x30}
				value, err := db.Get(key[:], nil)
				if err != nil {
					panic(err.Error())
				}
				// make copy of value slice; not supposed to alter the value slice?
				chunk := make([]byte, len(value))
				// copy(chunk, value)
				for i := 1; i < 256; i++ {
					chunk[i*128] = 20
					chunk[i*128+1] = 9
				}
				cx := 0
				cz := 0
				for i := 0; i < 128; i++ {
					chunk[i+128*136] = 91
					chunk[2048*cx+128*cz+i] = 91
					switch (i / 15) % 4 {
					case 0:
						cx++
					case 1:
						cz++
					case 2:
						cx--
					case 3:
						cz--
					}
				}
				for i := 1; i < 40; i++ {
					key[0] = byte(i)
					err = db.Put(key[:], chunk, nil)
					if err != nil {
						panic(err.Error())
					}
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
