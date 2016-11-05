package command

import (
	"fmt"
	"github.com/mitchellh/cli"
)

type ParentCommand interface {
	Name() string
	Subcommands() map[string]cli.CommandFactory
}

func RunSubcommand(c ParentCommand, meta Meta, args []string) int {
	if len(args) < 1 {
		meta.Ui.Error(fmt.Sprintf("missing %s subcommand", c.Name()))
		return cli.RunResultHelp
	}

	subc, subargs := args[0], args[1:]
	factory, ok := c.Subcommands()[subc]
	if !ok {
		meta.Ui.Error(fmt.Sprintf("unknown %s subcommand", c.Name()))
		return cli.RunResultHelp
	}

	cmd, err := factory()
	if err != nil {
		meta.Ui.Error(fmt.Sprintf("error creating %s subcommand", c.Name()))
		return 1
	}

	code := cmd.Run(subargs)
	if code == cli.RunResultHelp {
		meta.Ui.Error(cmd.Help())
		return 1
	}
	return code
}
