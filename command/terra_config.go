package command

import (
	"fmt"
	"github.com/BSick7/hatlas/terraform"
	"github.com/mitchellh/cli"
	"strings"
)

type TerraConfigCommand struct {
	Meta
}

func TerraConfigFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &TerraConfigCommand{Meta: meta}, nil
	}
}

func (c *TerraConfigCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("config")
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
	if err := c.getConfig(client, env); err != nil {
		c.Ui.Error(fmt.Sprintf("error getting config [%s]: %s", env, err))
		return 1
	}
	return 0
}

func (c *TerraConfigCommand) Synopsis() string {
	return "Download config for environment"
}

func (c *TerraConfigCommand) Help() string {
	helpText := `
Usage: hatlas terra config <environment>

  Downloads terraform config for <environment>.

  The available options allow the emitting of selected
  values within the state file based on terraform syntax.
`
	return strings.TrimSpace(helpText)
}

func (c *TerraConfigCommand) getConfig(client *terraform.AtlasClient, env string) error {
	stateRaw, err := client.GetTerraformConfig(env)
	if err != nil {
		return err
	}
	c.Ui.Info(string(stateRaw))
	return nil
}