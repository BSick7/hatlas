package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/BSick7/hatlas/structs"
	"github.com/BSick7/hatlas/terraform"
	"github.com/mitchellh/cli"
	"strings"
)

type TerraOutputsCommand struct {
	Meta
}

func TerraOutputsFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &TerraOutputsCommand{Meta: meta}, nil
	}
}

func (c *TerraOutputsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("outputs")
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
	if err := c.getOutputs(client, env); err != nil {
		c.Ui.Error(fmt.Sprintf("error getting state [%s]: %s", env, err))
		return 1
	}
	return 0
}

func (c *TerraOutputsCommand) Synopsis() string {
	return "Download state file for environment"
}

func (c *TerraOutputsCommand) Help() string {
	helpText := `
Usage: hatlas terra outputs [options] <environment>

  Downloads terraform state outputs for <environment>.

  The available options allow the emitting of selected
  values within the state file based on terraform syntax.

Terra State Options:
`
	return strings.TrimSpace(helpText)
}

func (c *TerraOutputsCommand) getOutputs(client *terraform.AtlasClient, env string) error {
	stateRaw, err := client.GetTerraformState(env)
	if err != nil {
		return err
	}

	ts := &structs.TerraformState{}
	if err := json.Unmarshal(stateRaw, ts); err != nil {
		return err
	}

	root := ts.GetRootModule()
	if root == nil {
		return nil
	}

	m := map[string]string{}
	for k, o := range root.Outputs {
		if o.Type == "string" {
			m[k] = fmt.Sprintf("%q", o.Value)
		} else if o.Type == "list" || o.Type == "map" {
			m[k] = fmt.Sprintf("%+v", o.Value)
		}
	}
	maxKeyLength := getMaxKeyLength(m)

	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("%s%s = %s\n", k, pad(" ", maxKeyLength-len(k)), v))
	}
	c.Ui.Info(buf.String())

	return nil
}

func getMaxKeyLength(m map[string]string) int {
	max := -1
	for k := range m {
		if len(k) > max {
			max = len(k)
		}
	}
	return max
}

func pad(pad string, length int) string {
	str := ""
	for {
		if len(str) == length {
			return str
		}
		str += pad
	}
}
