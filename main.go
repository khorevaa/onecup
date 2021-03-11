package main

import (
	"fmt"
	"github.com/khorevaa/onecup/cmd"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// nolint: gochecknoglobals
var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {

	app := &cli.App{
		Name:    "onecup",
		Version: buildVersion(),
		Authors: []*cli.Author{
			{
				Name: "Aleksey Khorev",
			},
		},
		Usage:       "Application for automate update configuration for 1C. Enterprise",
		Copyright:   "(c) 2021 Khorevaa",
		Description: "Application for automate update configuration for 1C. Enterprise",
	}

	for _, command := range cmd.Commands {
		app.Commands = append(app.Commands, command.Cmd())
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func buildVersion() string {
	var result = version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}
