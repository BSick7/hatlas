package command

import (
	"fmt"
	"github.com/BSick7/hatlas/terraform"
	"strings"
)

type TerraformCommand struct {
	Meta
}

func (c *TerraformCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("terraform")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	env := ""
	if fargs := flags.Args(); len(fargs) > 0 {
		env = fargs[0]
	}

	client := terraform.NewAtlasClient(nil)

	if env == "" {
		if err := c.listStates(client); err != nil {
			c.Ui.Error(fmt.Sprintf("error listing states: %s", err))
		}
	} else {
		if err := c.getState(client, env); err != nil {
			c.Ui.Error(fmt.Sprintf("error getting state: %s", err))
			return 1
		}
	}

	return 0
}

func (c *TerraformCommand) Synopsis() string {
	return "Run against terraform atlas API"
}

func (c *TerraformCommand) Help() string {
	helpText := `
Usage: hatlas terraform [options] <environment>

  Downloads terraform state file for <environment>.

  The available options allow the emitting of selected
  values within the state file based on terraform syntax.

  If no environment is specified, will list environments.

Terraform Options:
`
	return strings.TrimSpace(helpText)
}

func (c *TerraformCommand) listStates(client *terraform.AtlasClient) error {
	states, err := client.ListStates("")
	if err != nil {
		return err
	}
	c.Ui.Info(strings.Join(states.Names(), "\n"))
	return nil
}

func (c *TerraformCommand) getState(client *terraform.AtlasClient, env string) error {
	stateRaw, err := client.GetState(env)
	if err != nil {
		return err
	}
	c.Ui.Info(string(stateRaw))
	return nil
}
