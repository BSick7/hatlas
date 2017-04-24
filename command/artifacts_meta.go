package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/atlas-go/v1"
	"github.com/mitchellh/cli"
)

type ArtifactsMetaCommand struct {
	Meta
}

func ArtifactsMetaFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &ArtifactsMetaCommand{Meta: meta}, nil
	}
}

func (c *ArtifactsMetaCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("meta")
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
	if err := c.meta(client, username, name, artifactType); err != nil {
		c.Ui.Error(fmt.Sprintf("error listing artifacts: %s", err))
		return 1
	}
	return 0
}

func (c *ArtifactsMetaCommand) Synopsis() string {
	return "Dump metadata from all artifacts for a build"
}

func (c *ArtifactsMetaCommand) Help() string {
	helpText := `
Usage: hatlas artifacts meta [options] <username> <name> [<type>]

  Dump metadata from all artifacts for a build.

  <username> and <name> are required to pull the correct build.

  If [<type>] is specified, artifacts will be filtered by artifact type.

Artifacts Meta Options:
`
	return strings.TrimSpace(helpText)
}

func (c *ArtifactsMetaCommand) meta(client *atlas.Client, username string, name string, artifactType string) error {
	data := []*atlas.ArtifactVersion{}

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
			data = append(data, version)
		}
	}

	raw, _ := json.MarshalIndent(data, "", "  ")
	c.Ui.Info(string(raw))
	return nil
}
