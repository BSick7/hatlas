package command

import (
	"github.com/mitchellh/cli"
	"strings"
)

type TerraCommand struct {
	Meta
}

func TerraCommandFactory(meta Meta) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &TerraCommand{
			Meta: meta,
		}, nil
	}
}

func (c *TerraCommand) Subcommands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"list":   TerraListFactory(c.Meta),
		"state":  TerraStateFactory(c.Meta),
		"config": TerraConfigFactory(c.Meta),
	}
}

func (c *TerraCommand) Name() string {
	return "terra"
}

func (c *TerraCommand) Run(args []string) int {
	return RunSubcommand(c, c.Meta, args)
}

func (c *TerraCommand) Synopsis() string {
	return "Run against terraform atlas API"
}

func (c *TerraCommand) Help() string {
	helpText := `
Usage: hatlas terra <subcommand>

  Introspect terraform against the Atlas API.

Available subcommands:

  list      List environments
  state     Introspect state file
  config    Introspect configuration
`
	return strings.TrimSpace(helpText)
}
