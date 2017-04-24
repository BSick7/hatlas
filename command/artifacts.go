package command

import (
	"strings"

	"github.com/mitchellh/cli"
)

type ArtifactsCommand struct {
	Meta
}

func ArtifactsCommandFactory(meta Meta) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &ArtifactsCommand{
			Meta: meta,
		}, nil
	}
}

func (c *ArtifactsCommand) Subcommands() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"list": ArtifactsListFactory(c.Meta),
		"meta": ArtifactsMetaFactory(c.Meta),
		"push": ArtifactsPushFactory(c.Meta),
	}
}

func (c *ArtifactsCommand) Name() string {
	return "artifacts"
}

func (c *ArtifactsCommand) Run(args []string) int {
	return RunSubcommand(c, c.Meta, args)
}

func (c *ArtifactsCommand) Synopsis() string {
	return "Run against artifacts atlas API"
}

func (c *ArtifactsCommand) Help() string {
	helpText := `
Usage: hatlas artifacts <subcommand>

  Introspect artifacts against the Atlas API.

Available subcommands:

  list       List artifacts
  meta       Dump artifact metadata
  push       Push artifact metadata
`
	return strings.TrimSpace(helpText)
}
