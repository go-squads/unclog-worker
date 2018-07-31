package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BaritoLog/unclog-worker/cmds"
	"github.com/urfave/cli"
)

const (
	Name    = "unclog"
	Version = "0.0.0"
)

func main() {
	app := cli.App{
		Name:    Name,
		Usage:   "Provide kafka stream processor or stream processor log count for Project Unclog",
		Version: Version,
		Commands: []cli.Command{
			{
				Name:      "streamprocessor",
				ShortName: "sp",
				Usage:     "start unclog-worker as streamprocessor",
				Action:    cmds.ActionStreamProcessorService,
			},
			{
				Name:      "streamprocessorlogcount",
				ShortName: "splc",
				Usage:     "start unclog-worker as streamprocessorlogcount",
				Action:    cmds.ActionStreamProcessorLogCountService,
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(fmt.Sprintf("Some error occurred: %s", err.Error()))
	}
}
