package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-squads/unclog-worker/cmds"
	"github.com/go-squads/unclog-worker/config"
	"github.com/go-squads/unclog-worker/migration"
	"github.com/urfave/cli"
)

const (
	NAME    = "unclog"
	VERSION = "0.0.0"
)

func main() {
	err := setupConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	app := setupCLIApp()
	startCLI(&app)
}

func startCLI(app *cli.App) {
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(fmt.Sprintf("Some error occurred: %s", err.Error()))
	}
}

func setupCLIApp() cli.App {
	app := cli.App{
		Name:    NAME,
		Usage:   "Provide kafka stream processor or stream processor log count for Project Unclog",
		Version: VERSION,
		Commands: []cli.Command{
			{
				Name:      "streamprocessor",
				ShortName: "sp",
				Usage:     "Start unclog-worker as stream processor",
				Action:    cmds.ActionStreamProcessorService,
			},
			{
				Name:      "streamprocessorlogcount",
				ShortName: "splc",
				Usage:     "Start unclog-worker as stream processor log counter",
				Action:    cmds.ActionStreamProcessorLogCountService,
			},
			{
				Name:      "migrator",
				ShortName: "migrate",
				Usage:     "run the database migration",
				Action:    migration.RunMigration,
			},
		},
	}

	return app
}

func setupConfiguration() error {
	err := config.SetupConfig()
	if err != nil {
		return err
	}

	return nil
}
