package cmd

import (
	"github.com/urfave/cli/v2"
)

type updateCommand struct {
}

func (c *updateCommand) run(context *cli.Context) error {

	return nil
}

func (c *updateCommand) Flags() []cli.Flag {

	return []cli.Flag{
		&cli.StringFlag{
			Destination: &c.port, Name: "port",
			Value: ":3001", Usage: "port to listen on"},
		&cli.StringFlag{
			Destination: &c.appServer, Name: "server",
			Value: "localhost:1545", Usage: "ras client address with port"},
		&cli.BoolFlag{
			Destination: &c.debug, Name: "debug",
			Value: false, Usage: "debug mode"},
	}
}

func (c *updateCommand) Cmd() *cli.Command {

	cmd := &cli.Command{
		//Category:    "some_category",
		Name:        "some_command",
		Usage:       "Usage decription",
		Description: `Full usage description `,
		Action:      c.run,
		Flags:       c.Flags(),
		},
	}

	return cmd
}
