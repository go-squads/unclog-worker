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
	Version = "0.7.0"
)

func main() {
	app := cli.App{
		Name:    Name,
		Usage:   "Provide kafka stream processor for Unclog project",
		Version: Version,
		Action:  cmds.ActionStreamProcessorService,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(fmt.Sprintf("Some error occurred: %s", err.Error()))
	}
}
