package cmd

import (
	"github.com/urfave/cli/v2"
)

var Commands = []Command{

	&execCommand{},
	&appCommand{},
	//&commandWithSub{
	//	sub: []Command{
	//		&subCommand{},
	//	},
	//},
}

type Command interface {
	Cmd() *cli.Command
}
