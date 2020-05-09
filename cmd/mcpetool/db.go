package main

import (
	"encoding/hex"
	"fmt"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli/v2"
)

var dbCommand = cli.Command{
	Name:  "db",
	Usage: "List, get, put, or delete leveldb keys",
	Subcommands: []*cli.Command{
		{
			Name:    "list",
			Aliases: []string{"keys", "k"},
			Flags: []cli.Flag{
				&pathFlag,
			},
			Usage: "Lists all keys in the database.",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(worldPath)
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
				&pathFlag,
				&outFlag,
				&jsonFlag,
				&yamlFlag,
				&dumpFlag,
				&binaryFlag,
			},
			Action: func(c *cli.Context) error {
				var outData []byte
				var err error
				nbt2json.UseBedrockEncoding()
				nbt2json.UseLongAsString()
				world, err := world.OpenWorld(worldPath)
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
				comment := jsonComment
				if stringKey != "" {
					comment += " | ASCII Key " + stringKey
				}
				comment += " | Hex Key " + hexKey + " | Path " + worldPath
				if c.String("dump") == "true" {
					outData = []byte(hex.Dump(value))
				} else if c.String("json") == "true" {
					outData, err = nbt2json.Nbt2Json(value, comment)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				} else if c.String("yaml") == "true" {
					outData, err = nbt2json.Nbt2Yaml(value, comment)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					// } else if c.String("base64") == "true" {
					// 	outData = []byte(base64.StdEncoding.EncodeToString(value))
				} else if c.String("binary") == "true" {
					outData = value
				} else {
					outData, err = nbt2json.Nbt2Json(value, comment)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				}
				err = writeOutput(outFile, outData)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
		{
			Name:      "put",
			ArgsUsage: "<key>",
			Usage:     "Put a key/value into the DB. Overwrites the key if already present. Input is base64-formatted by default.",
			Flags: []cli.Flag{
				&pathFlag,
				&inFlag,
				&jsonFlag,
				&yamlFlag,
				&binaryFlag,
			},
			Action: func(c *cli.Context) error {
				var value []byte
				var err error
				nbt2json.UseBedrockEncoding()
				nbt2json.UseLongAsString()
				world, err := world.OpenWorld(worldPath)
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				inputData, err := readInput(inFile)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				if c.String("json") == "true" {
					value, err = nbt2json.Json2Nbt(inputData)
				} else if c.String("yaml") == "true" {
					value, err = nbt2json.Yaml2Nbt(inputData)
					// } else if c.String("base64") == "true" {
					// 	value, err = base64.StdEncoding.DecodeString(string(inputData[:]))
				} else if c.String("binary") == "true" {
					value = inputData
				} else {
					value, err = nbt2json.Json2Nbt(inputData[:])
				}
				if err != nil {
					return cli.NewExitError(err, 1)
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
			Flags: []cli.Flag{
				&pathFlag,
			},
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(worldPath)
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
	},
}
