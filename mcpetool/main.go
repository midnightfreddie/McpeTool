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

const appVersion = "0.3.0-beta"
const jsonComment = "MCPE Tool v" + appVersion

var worldPath, inFile, outFile string

var pathFlag = cli.StringFlag{
	Name:        "path, p",
	Value:       ".",
	Usage:       "`FILEPATH` of world",
	EnvVar:      "MCPETOOL_WORLD",
	Destination: &worldPath,
}
var inFlag = cli.StringFlag{
	Name:        "in, i",
	Value:       "-",
	Usage:       "Input `FILE` path",
	Destination: &inFile,
}
var outFlag = cli.StringFlag{
	Name:        "out, o",
	Value:       "-",
	Usage:       "Output `FILE` path",
	Destination: &outFile,
}
var dumpFlag = cli.BoolFlag{
	Name:  "dump, d",
	Usage: "Hexdump format",
}
var base64Flag = cli.BoolFlag{
	Name:  "base64",
	Usage: "Base64 format",
}
var jsonFlag = cli.BoolFlag{
	Name:  "json, j",
	Usage: "JSON format",
}
var yamlFlag = cli.BoolFlag{
	Name:  "yaml, yml, y",
	Usage: "YAML format",
}
var binaryFlag = cli.BoolFlag{
	Name:  "binary",
	Usage: "Raw binary",
}

// Read from file or from stdin if inFile is "-"
func readInput(inFile string) ([]byte, error) {
	var inData []byte
	var err error
	if inFile == "-" {
		inData, err = ioutil.ReadAll(os.Stdin)
	} else {
		inData, err = ioutil.ReadFile(inFile)
	}
	return inData, err
}

// Write to file or to stdout if outFile is "-"
func writeOutput(outFile string, outData []byte) error {
	var err error
	if outFile == "-" {
		err = binary.Write(os.Stdout, binary.LittleEndian, outData)
	} else {
		err = ioutil.WriteFile(outFile, outData, 0644)
	}
	return err
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
		pathFlag,
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
