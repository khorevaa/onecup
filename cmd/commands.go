package cmd

import (
	"github.com/urfave/cli/v2"
)

var Commands = []Command{

	&updateCommand{},
	//&commandWithSub{
	//	sub: []Command{
	//		&subCommand{},
	//	},
	//},
}

type Command interface {
	Cmd() *cli.Command
}
