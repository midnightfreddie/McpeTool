package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "MCPE Tool"
	app.Version = "0.0.0"
	app.Usage = "A utility to access Minecraft Pocket Edition .mcworld exported world files."

	app.Commands = []cli.Command{
		{
			Name:    "api",
			Aliases: []string{"www"},
			Usage:   "Open world and start http API. Hit control-c to exit.",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(c.Args().First())
				if err != nil {
					panic("error")
				}
				defer world.Close()
				err = api.Server(&world)
				if err != nil {
					panic("error")
				}
				return nil
			},
		},
		{
			Name:    "keys",
			Aliases: []string{"k"},
			Usage:   "Lists all keys in the database in base64 format. Be sure to include the path to the world folder, e.g. 'McpeTool keys path/to/world'",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(c.Args().First())
				if err != nil {
					return err
				}
				defer world.Close()
				keys, err := world.GetKeys()
				if err != nil {
					return err
				}
				for i := 0; i < len(keys); i++ {
					fmt.Println(base64.StdEncoding.EncodeToString(keys[i]))
				}
				return nil
			},
		},
		{
			Name:  "get",
			Usage: "Retruns the value of a key. Both key and value are in base64 format. e.g. 'McpeTool get path/to/world AAAAAAAAAAAw' for terrain chunk 0,0 or 'McpeTool get path/to/world fmxvY2FsX3BsYXllcg==' for ~local_player player data",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(c.Args().First())
				if err != nil {
					return err
				}
				defer world.Close()
				key, err := base64.StdEncoding.DecodeString(c.Args().Get(1))
				if err != nil {
					return err
				}
				value, err := world.Get(key)
				if err != nil {
					return err
				}
				fmt.Println(base64.StdEncoding.EncodeToString(value))
				return nil
			},
		},
		{
			Name:  "delete",
			Usage: "Deletes a key and its value. The key is in base64 format. e.g. 'McpeTool delete path/to/world AAAAAAAAAAAw' to delete terrain chunk 0,0 or 'McpeTool delete path/to/world fmxvY2FsX3BsYXllcg==' to delete ~local_player player data",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(c.Args().First())
				if err != nil {
					return err
				}
				defer world.Close()
				key, err := base64.StdEncoding.DecodeString(c.Args().Get(1))
				if err != nil {
					return err
				}
				err = world.Delete(key)
				if err != nil {
					return err
				}
				return nil
			},
		},
		{
			Name:    "develop",
			Aliases: []string{"dev"},
			Usage:   "Random thing the dev is working on",
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(c.Args().First())
				if err != nil {
					return err
				}
				defer world.Close()
				keys, err := world.GetKeys()
				if err != nil {
					return err
				}
				fmt.Printf("%v\n", keys)
				fmt.Println(world.FilePath())
				return nil
			},
		},
	}

	app.Run(os.Args)
}
