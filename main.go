package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BaritoLog/unclog/cmds"
	"github.com/urfave/cli"
)

const (
	Name    = "unclog"
	Version = "0.7.0"
)

func main() {
	app := cli.App{
		Name:    Name,
		Usage:   "Provide kafka producer or consumer for Barito project",
		Version: Version,
		Commands: []cli.Command{
			// {
			// 	Name:      "producer",
			// 	ShortName: "p",
			// 	Usage:     "start unclog as producer",
			// 	Action:    cmds.ActionBaritoProducerService,
			// },
			{
				Name:      "consumer",
				ShortName: "c",
				Usage:     "start unclog as consumer",
				Action:    cmds.ActionBaritoConsumerService,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(fmt.Sprintf("Some error occurred: %s", err.Error()))
	}
}
