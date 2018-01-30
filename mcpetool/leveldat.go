package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"

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
			Usage: "Returns level.dat in JSON format",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "dump, d",
					Usage: "Output hexdump",
				},
				cli.BoolFlag{
					Name:  "yaml, y",
					Usage: "Output YAML",
				},
				cli.BoolFlag{
					Name:  "binary",
					Usage: "Output binary. Only use when redirecting output.",
				},
			},
			Action: func(c *cli.Context) error {
				myWorld, err := world.OpenWorld(worldPath)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer myWorld.Close()
				levelDat, err := myWorld.GetLevelDat()
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				switch {
				case c.String("dump") == "true":
					fmt.Println(hex.Dump(levelDat))
				case c.String("yaml") == "true":
					out, err := nbt2json.Nbt2Yaml(levelDat, binary.LittleEndian, jsonComment+" | level.dat | Path "+worldPath)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					fmt.Println(string(out[:]))
				case c.String("binary") == "true":
					err = binary.Write(os.Stdout, binary.LittleEndian, levelDat)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				default:
					out, err := nbt2json.Nbt2Json(levelDat, binary.LittleEndian, jsonComment+" | level.dat | Path "+worldPath)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
					fmt.Println(string(out[:]))
				}
				return nil
			},
		},
	},
}
