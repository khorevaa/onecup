package cmd

import (
	"github.com/khorevaa/onecup/app"
	"github.com/urfave/cli/v2"
)

type appCommand struct {
	configFile string
	simulate   bool
}

func (c *appCommand) run(context *cli.Context) error {

	return app.Run()
}

func (c *appCommand) Flags() []cli.Flag {

	return []cli.Flag{
		//&cli.StringFlag{
		//	Destination: &c.configFile, Name: "config",
		//	Value: "", Usage: "job config file `FILE`"},
		//&cli.BoolFlag{
		//	Destination: &c.simulate, Name: "simulate", Aliases: []string{"s"},
		//	Value: false, Usage: "simulate job work"},
	}
}

func (c *appCommand) Cmd() *cli.Command {

	cmd := &cli.Command{
		//Category:    "some_category",
		Name:        "app",
		Usage:       "run GUI",
		Description: "This command help execute any job from config file",
		Action:      c.run,
		Flags:       c.Flags(),
	}

	return cmd
}
