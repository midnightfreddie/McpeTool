package main

import (
	"encoding/binary"
	"encoding/hex"

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
			Usage: "Returns level.dat in nbt-to-JSON format",
			Flags: []cli.Flag{
				pathFlag,
				outFlag,
				dumpFlag,
				yamlFlag,
				binaryFlag,
			},
			Action: func(c *cli.Context) error {
				var outData []byte
				var err error
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
					outData = []byte(hex.Dump(levelDat))
				case c.String("yaml") == "true":
					outData, err = nbt2json.Nbt2Yaml(levelDat, binary.LittleEndian, jsonComment+" | level.dat | Path "+worldPath)
					if err != nil {
						return cli.NewExitError(err, 1)
					}
				case c.String("binary") == "true":
					outData = levelDat
				default:
					outData, err = nbt2json.Nbt2Json(levelDat, binary.LittleEndian, jsonComment+" | level.dat | Path "+worldPath)
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
	},
}
