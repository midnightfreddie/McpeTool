package main

import (
	"os"
	"time"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/urfave/cli"
)

const appVersion = "0.2.2-alpha"
const jsonComment = "MCPE Tool v" + appVersion

var path string

func main() {
	app := cli.NewApp()
	app.Name = "mcpetool"
	app.Version = appVersion
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
		levelDatCommand,
		dbCommand,
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
	}

	app.Run(os.Args)
}
