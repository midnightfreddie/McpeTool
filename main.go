package main

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

func main() {
	var path string
	app := cli.NewApp()
	app.Name = "MCPE Tool"
	app.Version = "0.2.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Jim Nelson",
			Email: "jim@jimnelson.us",
		},
	}
	app.Copyright = "(c) 2018 Jim Nelson"
	app.Usage = "Reads and writes a Minecraft Pocket Edition world directory."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "path, p",
			Value:       ".",
			Usage:       "`FILEPATH` of world",
			EnvVar:      "MCPETOOL_WORLD",
			Destination: &path,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "api",
			Aliases: []string{"www"},
			Usage:   "Open world, start API at http://127.0.0.1:8080 . Control-c to exit.",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(path)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				err = api.Server(&world)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		{
			Name:    "keys",
			Aliases: []string{"k"},
			Usage:   "Lists all keys in the database.",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(path)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				keys, err := world.GetKeys()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				for i := 0; i < len(keys); i++ {
					fmt.Println(hex.EncodeToString(keys[i]))
				}
				return nil
			},
		},
		{
			Name:      "get",
			ArgsUsage: "<key>",
			Usage:     "Retruns a key's value in base64 format.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dump, d",
					Usage: "Display value as hexdump",
				},
				cli.BoolFlag{
					Name:  "json, j",
					Usage: "Display value as JSON. Only valid if value is NBT.",
				},
				cli.BoolFlag{
					Name:  "yaml, y",
					Usage: "Display value as YAML. Only valid if value is NBT.",
				},
			},
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(path)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				value, err := world.Get(key)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				stringKey, hexKey := api.ConvertKey(key)
				comment := "McpeTool v" + app.Version
				if stringKey != "" {
					comment += " | ASCII Key " + stringKey
				}
				comment += " | Hex Key " + hexKey + " | Path " + path
				if c.String("dump") == "true" {
					fmt.Println(hex.Dump(value))
				} else if c.String("json") == "true" {
					out, err := nbt2json.Nbt2Json(value, binary.LittleEndian, comment)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					fmt.Println(string(out[:]))
				} else if c.String("yaml") == "true" {
					out, err := nbt2json.Nbt2Yaml(value, binary.LittleEndian, comment)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					fmt.Println(string(out[:]))
				} else {
					fmt.Println(base64.StdEncoding.EncodeToString(value))
				}
				return nil
			},
		},
		{
			Name:      "put",
			ArgsUsage: "<key>",
			Usage:     "Put a key/value into the DB. The base64-encoded value read from stdin.",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json, j",
					Usage: "Use nbt2json JSON data as input",
				},
				cli.BoolFlag{
					Name:  "yaml, y",
					Usage: "Use YAML-ized nbt2json data as input",
				},
			},
			Action: func(c *cli.Context) error {
				var value []byte
				world, err := world.OpenWorld(path)
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				inputData, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				if c.String("json") == "true" {
					value, err = nbt2json.Json2Nbt(inputData[:], binary.LittleEndian)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				} else if c.String("yaml") == "true" {
					value, err = nbt2json.Yaml2Nbt(inputData, binary.LittleEndian)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				} else {
					value, err = base64.StdEncoding.DecodeString(string(inputData[:]))
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}
				err = world.Put(key, value)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		{
			Name:      "delete",
			ArgsUsage: "<key>",
			Usage:     "Deletes a key and its value.",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(path)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				err = world.Delete(key)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
