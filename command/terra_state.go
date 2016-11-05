package command

import (
	"fmt"
	"github.com/BSick7/hatlas/terraform"
	"github.com/mitchellh/cli"
	"strings"
)

type TerraStateCommand struct {
	Meta
}

func TerraStateFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &TerraStateCommand{Meta: meta}, nil
	}
}

func (c *TerraStateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("state")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	fargs := flags.Args()
	if len(fargs) < 1 {
		c.Ui.Error("missing terraform environment")
		return cli.RunResultHelp
	}
	env := fargs[0]

	client := terraform.NewAtlasClient(nil)
	if err := c.getState(client, env); err != nil {
		c.Ui.Error(fmt.Sprintf("error getting state [%s]: %s", env, err))
		return 1
	}
	return 0
}

func (c *TerraStateCommand) Synopsis() string {
	return "Download state file for environment"
}

func (c *TerraStateCommand) Help() string {
	helpText := `
Usage: hatlas terra state [options] <environment>

  Downloads terraform state file for <environment>.

  The available options allow the emitting of selected
  values within the state file based on terraform syntax.

Terra State Options:
`
	return strings.TrimSpace(helpText)
}

func (c *TerraStateCommand) getState(client *terraform.AtlasClient, env string) error {
	stateRaw, err := client.GetTerraformState(env)
	if err != nil {
		return err
	}
	c.Ui.Info(string(stateRaw))
	return nil
}
