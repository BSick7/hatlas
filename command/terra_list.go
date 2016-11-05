package command

import (
	"fmt"
	"github.com/BSick7/hatlas/terraform"
	"github.com/mitchellh/cli"
	"strings"
)

type TerraListCommand struct {
	Meta
}

func TerraListFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &TerraListCommand{Meta: meta}, nil
	}
}

func (c *TerraListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("list")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	username := ""
	if fargs := flags.Args(); len(fargs) > 0 {
		username = fargs[0]
	}

	client := terraform.NewAtlasClient(nil)
	if err := c.list(client, username); err != nil {
		c.Ui.Error(fmt.Sprintf("error listing states: %s", err))
		return 1
	}
	return 0
}

func (c *TerraListCommand) Synopsis() string {
	return "List terraform environments in atlas"
}

func (c *TerraListCommand) Help() string {
	helpText := `
Usage: hatlas terra list [options] <username>

  Lists all terraform environments available.

  If <username> is specified, list will only
  return environments owned by <username>.

Terra List Options:
`
	return strings.TrimSpace(helpText)
}

func (c *TerraListCommand) list(client *terraform.AtlasClient, username string) error {
	states, err := client.ListTerraforms(username)
	if err != nil {
		return err
	}
	c.Ui.Info(strings.Join(states.Names(), "\n"))
	return nil
}
