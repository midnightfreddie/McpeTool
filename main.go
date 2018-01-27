package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/urfave/cli"
)

func main() {
	var path string
	app := cli.NewApp()
	app.Name = "MCPE Tool"
	app.Version = "0.1.2b"
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
					return cli.NewExitError(err, 1);
				}
				defer world.Close()
				err = api.Server(&world)
				if err != nil {
					return cli.NewExitError(err, 1);
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
					return cli.NewExitError(err, 1);
				}
				defer world.Close()
				keys, err := world.GetKeys()
				if err != nil {
					return cli.NewExitError(err, 1);
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
			},
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(path)
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				defer world.Close()
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				value, err := world.Get(key)
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				if c.String("dump") == "true" {
					fmt.Println(hex.Dump(value))
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
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(path)
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				defer world.Close()
				base64Data, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				value, err := base64.StdEncoding.DecodeString(string(base64Data[:]))
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				err = world.Put(key, value)
				if err != nil {
					return cli.NewExitError(err, 1);
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
					return cli.NewExitError(err, 1);
				}
				defer world.Close()
				key, err := hex.DecodeString(c.Args().Get(0))
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				err = world.Delete(key)
				if err != nil {
					return cli.NewExitError(err, 1);
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
