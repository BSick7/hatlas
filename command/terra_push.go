package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BSick7/hatlas/structs"
	"github.com/BSick7/hatlas/terraform"
	"github.com/mitchellh/cli"
)

type TerraPushCommand struct {
	Meta
}

func TerraPushFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &TerraPushCommand{Meta: meta}, nil
	}
}

func (c *TerraPushCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("config")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	fargs := flags.Args()
	if len(fargs) < 1 {
		c.Ui.Error("subcommand not specified")
		return cli.RunResultHelp
	}
	subcommand := fargs[0]
	switch subcommand {
	case "vars":
		return c.handleVars(fargs[1:])
	case "state":
		return c.handleState(fargs[1:])
	default:
		c.Ui.Error(fmt.Sprintf("unknown subcommand %q", subcommand))
		return cli.RunResultHelp
	}
}

func (c *TerraPushCommand) Synopsis() string {
	return "Upload vars to environment"
}

func (c *TerraPushCommand) Help() string {
	helpText := `
Usage: hatlas terra push vars|state <environment> <var-file>|<state-file>

  Updates terraform environment

  When 'vars' is specified, command expects <var-file> to push to atlas.
  When 'state' is specified, command expects <state-file> to push to atlas.
`
	return strings.TrimSpace(helpText)
}

func (c *TerraPushCommand) handleVars(args []string) int {
	if len(args) < 1 {
		c.Ui.Error("terraform environment not specified")
		return cli.RunResultHelp
	}
	if len(args) < 2 {
		c.Ui.Error("var file not specified")
		return cli.RunResultHelp
	}
	env := args[0]
	varfile := args[1]

	vars, err := c.getVars(varfile)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error reading vars file [%s]: %s", varfile, err))
		return 1
	}

	client := terraform.NewAtlasClient(nil)
	response, err := client.UpdateVariables(env, vars)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error pushing vars [%s]: %s", env, err))
		return 1
	}
	c.Ui.Info(string(response))
	return 0
}

func (c *TerraPushCommand) handleState(args []string) int {
	if len(args) < 1 {
		c.Ui.Error("terraform environment not specified")
		return cli.RunResultHelp
	}
	if len(args) < 2 {
		c.Ui.Error("state file not specified")
		return cli.RunResultHelp
	}
	env := args[0]
	statefile := args[1]

	raw, err := ioutil.ReadFile(statefile)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("unable to read state file: %s", err))
		return 1
	}

	client := terraform.NewAtlasClient(nil)
	response, err := client.UpdateState(env, raw)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error pushing state [%s]: %s", env, err))
		return 1
	}
	c.Ui.Info(string(response))
	return 0
}

func (c *TerraPushCommand) getVars(varfile string) (*structs.TerraformRawConfig, error) {
	raw, err := ioutil.ReadFile(varfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read var file: %s", err)
	}

	trc, err := structs.NewTerraformRawConfigFromJson(raw)
	if err != nil {
		return nil, fmt.Errorf("cannot read hcl files yet: %s", varfile)
	}

	return trc, nil
}
