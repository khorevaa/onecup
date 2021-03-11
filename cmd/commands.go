package cmd

import (
	"github.com/urfave/cli/v2"
)

var Commands = []Command{

	&execCommand{},
	//&commandWithSub{
	//	sub: []Command{
	//		&subCommand{},
	//	},
	//},
}

type Command interface {
	Cmd() *cli.Command
}
