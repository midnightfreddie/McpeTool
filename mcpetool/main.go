package main

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"time"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/urfave/cli"
)

const appVersion = "0.2.2-alpha"
const jsonComment = "MCPE Tool v" + appVersion

var worldPath string

// Write to file or to stdout if outFile is "-"
func writeOutput(outFile string, outData []byte) error {
	if outFile == "-" {
		err := binary.Write(os.Stdout, binary.LittleEndian, outData)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
	} else {
		err := ioutil.WriteFile(outFile, outData, 0644)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
	}
	return nil
}

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
			Destination: &worldPath,
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
				world, err := world.OpenWorld(worldPath)
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
