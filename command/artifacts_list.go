package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/atlas-go/v1"
	"github.com/mitchellh/cli"
)

var (
	ARTIFACT_TYPES = []string{
		"amazon.ami",
		"amazon.image",
		"azure.image",
		"cloudstack.image",
		"digitalocean.image",
		"docker.image",
		"googlecompute.image",
		"hyperv.image",
		"oneandone.image",
		"openstack.image",
		"parallels.image",
		"profitbricks.image",
		"qemu.image",
		"triton.image",
		"virtualbox.image",
		"vmware.image",
		"custom.image",
		"vagrant.box",
	}
)

type ArtifactsListCommand struct {
	Meta
}

func ArtifactsListFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &ArtifactsListCommand{Meta: meta}, nil
	}
}

func (c *ArtifactsListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("list")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	fargs := flags.Args()
	if len(fargs) == 0 {
		c.Ui.Error("missing username and name")
		return cli.RunResultHelp
	} else if len(fargs) == 1 {
		c.Ui.Error("missing name")
		return cli.RunResultHelp
	}
	username := fargs[0]
	name := fargs[1]
	artifactType := ""
	if len(fargs) > 2 {
		artifactType = fargs[2]
	}

	client := atlas.DefaultClient()
	if err := c.list(client, username, name, artifactType); err != nil {
		c.Ui.Error(fmt.Sprintf("error listing artifacts: %s", err))
		return 1
	}
	return 0
}

func (c *ArtifactsListCommand) Synopsis() string {
	return "List terraform environments in atlas"
}

func (c *ArtifactsListCommand) Help() string {
	helpText := `
Usage: hatlas artifacts list [options] <username> <name> [<type>]

  Lists all artifacts for a particular build.

  <username> and <name> are required to search through the build artifacts.

  If [<type>] is specified, artifacts list will be filtered by artifact type.

Artifacts List Options:
`
	return strings.TrimSpace(helpText)
}

func (c *ArtifactsListCommand) list(client *atlas.Client, username string, name string, artifactType string) error {
	slugs := []string{}

	types := ARTIFACT_TYPES
	if artifactType != "" {
		types = []string{artifactType}
	}

	for _, artifactType := range types {
		versions, err := client.ArtifactSearch(&atlas.ArtifactSearchOpts{
			User: username,
			Name: name,
			Type: artifactType,
		})
		if err != nil {
			return err
		}
		for _, version := range versions {
			slugs = append(slugs, version.Slug)
		}
	}

	c.Ui.Info(strings.Join(slugs, "\n"))
	return nil
}
