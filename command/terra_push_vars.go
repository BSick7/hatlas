package command

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/BSick7/hatlas/structs"
	"github.com/BSick7/hatlas/terraform"
	"github.com/mitchellh/cli"
)

type TerraPushVarsCommand struct {
	Meta
}

func TerraPushVarsFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &TerraPushVarsCommand{Meta: meta}, nil
	}
}

func (c *TerraPushVarsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("config")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	fargs := flags.Args()
	if len(fargs) < 1 {
		c.Ui.Error("terraform environment not specified")
		return cli.RunResultHelp
	}
	if len(fargs) < 2 {
		c.Ui.Error("var file not specified")
		return cli.RunResultHelp
	}
	env := fargs[0]
	varfile := fargs[1]

	vars, err := c.getVars(varfile)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("error reading vars file [%s]: %s", varfile, err))
		return 1
	}

	client := terraform.NewAtlasClient(nil)
	if err := client.UpdateVariables(env, vars); err != nil {
		c.Ui.Error(fmt.Sprintf("error pushing vars [%s]: %s", env, err))
		return 1
	}

	return 0
}

func (c *TerraPushVarsCommand) Synopsis() string {
	return "Upload vars to environment"
}

func (c *TerraPushVarsCommand) Help() string {
	helpText := `
Usage: hatlas terra push-vars <environment> <var-file>

  Upload terraform vars from <var-file> to <environment>.
`
	return strings.TrimSpace(helpText)
}

func (c *TerraPushVarsCommand) getVars(varfile string) (*structs.TerraformRawConfig, error) {
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
