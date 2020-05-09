package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/midnightfreddie/McpeTool/api"
	"github.com/midnightfreddie/McpeTool/world"
	"github.com/urfave/cli/v2"
)

const appVersion = "0.3.2"
const jsonComment = "MCPE Tool v" + appVersion

var worldPath, inFile, outFile string

var pathFlag = cli.StringFlag{
	Name:        "path",
	Aliases:     []string{"p"},
	Value:       ".",
	Usage:       "`FILEPATH` of world",
	EnvVars:     []string{"MCPETOOL_WORLD"},
	Destination: &worldPath,
}
var inFlag = cli.StringFlag{
	Name:        "in",
	Aliases:     []string{"i"},
	Value:       "-",
	Usage:       "Input `FILE` path",
	Destination: &inFile,
}
var outFlag = cli.StringFlag{
	Name:        "out",
	Aliases:     []string{"o"},
	Value:       "-",
	Usage:       "Output `FILE` path",
	Destination: &outFile,
}
var dumpFlag = cli.BoolFlag{
	Name:    "dump",
	Aliases: []string{"d"},
	Usage:   "Hexdump format",
}
var base64Flag = cli.BoolFlag{
	Name:  "base64",
	Usage: "Base64 format",
}
var jsonFlag = cli.BoolFlag{
	Name:    "json",
	Aliases: []string{"j"},
	Usage:   "JSON format",
}
var yamlFlag = cli.BoolFlag{
	Name:    "yaml",
	Aliases: []string{"yml", "y"},
	Usage:   "YAML format",
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
	cli.AppHelpTemplate = fmt.Sprintf(`%s
WEBSITE: https://github.com/midnightfreddie/McpeTool
`, cli.AppHelpTemplate)

	app := cli.NewApp()
	app.Name = "mcpetool"
	app.Version = appVersion
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Jim Nelson",
			Email: "jim@jimnelson.us",
		},
	}
	app.Copyright = "(c) 2018, 2020 Jim Nelson"
	app.Usage = "Reads and writes a Minecraft Bedrock Edition world directory."
	app.Commands = []*cli.Command{
		&levelDatCommand,
		&dbCommand,
		{
			Name:    "api",
			Aliases: []string{"www"},
			Usage:   "Open world, start API at http://127.0.0.1:8080 . Control-c to exit.",
			Flags: []cli.Flag{
				&pathFlag,
				&cli.StringFlag{
					Name:    "addr",
					Value:   "127.0.0.1",
					Usage:   "`ADDRESS` on which to bind",
					EnvVars: []string{"MCPETOOL_ADDR"},
				},
				&cli.StringFlag{
					Name:    "port",
					Value:   "8080",
					Usage:   "`PORT` on which to listen",
					EnvVars: []string{"MCPETOOL_PORT"},
				},
			},
			Action: func(c *cli.Context) error {
				world, err := world.OpenWorld(worldPath)
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				defer world.Close()
				fmt.Println("Starting API server for world at " + worldPath)
				fmt.Println("REST at http://" + c.String("addr") + ":" + c.String("port") + "/api/v1/db")
				fmt.Println("Press control-C to exit")
				err = api.Server(&world, c.String("addr"), c.String("port"))
				if err != nil {
					return cli.NewExitError(err, 1)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
