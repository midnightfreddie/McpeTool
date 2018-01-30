package main

import (
	"encoding/binary"
	"fmt"

	"github.com/midnightfreddie/McpeTool/world"
	"github.com/midnightfreddie/nbt2json"
	"github.com/urfave/cli"
)

var levelDatCommand = cli.Command{
	Name:  "leveldat",
	Usage: "Get or put level.dat data",
	Subcommands: []cli.Command{
		{
			Name:  "get",
			Usage: "Returns level.dat in nbt2json YAML format",
			Action: func(c *cli.Context) error {
				myWorld, err := world.OpenWorld(path)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer myWorld.Close()
				levelDat, err := myWorld.GetLevelDat()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				out, err := nbt2json.Nbt2Yaml(levelDat, binary.LittleEndian, jsonComment+" | level.dat | Path "+path)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				fmt.Println(string(out[:]))

				return nil
			},
		},
	},
}
