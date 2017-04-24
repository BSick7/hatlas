package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hashicorp/atlas-go/v1"
	"github.com/mitchellh/cli"
)

type ArtifactsPushCommand struct {
	Meta
}

func ArtifactsPushFactory(meta Meta) func() (cli.Command, error) {
	return func() (cli.Command, error) {
		return &ArtifactsPushCommand{Meta: meta}, nil
	}
}

func (c *ArtifactsPushCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("meta")
	flags.Usage = func() {
		c.Ui.Error(c.Help())
	}
	if err := flags.Parse(args); err != nil {
		return 1
	}

	fargs := flags.Args()
	if len(fargs) == 0 {
		c.Ui.Error("missing username, name, and artifacts-json-file")
		return cli.RunResultHelp
	} else if len(fargs) == 1 {
		c.Ui.Error("missing name and artifacts-json-file")
		return cli.RunResultHelp
	} else if len(fargs) == 2 {
		c.Ui.Error("missing artifacts-json-file")
		return cli.RunResultHelp
	}

	username := fargs[0]
	name := fargs[1]
	metafile := fargs[2]

	client := atlas.DefaultClient()
	if err := c.push(client, username, name, metafile); err != nil {
		c.Ui.Error(fmt.Sprintf("error pushing artifact metadata: %s", err))
		return 1
	}
	return 0
}

func (c *ArtifactsPushCommand) Synopsis() string {
	return "Pushes artifact metadata to a single build in atlas."
}

func (c *ArtifactsPushCommand) Help() string {
	helpText := `
Usage: hatlas artifacts push [options] <username> <name> <artifacts-json-file>

  Pushes artifact metadata to a single build in atlas.
  If metadata version already exists, hatlas will skip the push.

  <username> and <name> are required to push to the correct build.

  <artifacts-json-file> should be a json file containing a listing of artifact metadata.

Artifacts Push Options:
`
	return strings.TrimSpace(helpText)
}

func (c *ArtifactsPushCommand) push(client *atlas.Client, username string, name string, metafile string) error {
	if err := c.ensureArtifact(client, username, name); err != nil {
		return err
	}

	versions, err := c.getMeta(metafile)
	if err != nil {
		return err
	}

	for _, version := range versions {
		c.pushSingleArtifact(client, version)
	}

	return nil
}

func (c *ArtifactsPushCommand) getMeta(metafile string) ([]*atlas.ArtifactVersion, error) {
	raw, err := ioutil.ReadFile(metafile)
	if err != nil {
		return nil, fmt.Errorf("unable to read meta file: %s", err)
	}

	versions := []*atlas.ArtifactVersion{}
	if err := json.Unmarshal(raw, &versions); err != nil {
		return nil, fmt.Errorf("error reading artifact meta file %q: %s", metafile, err)
	}

	return versions, nil
}

func (c *ArtifactsPushCommand) ensureArtifact(client *atlas.Client, username string, name string) error {
	if _, err := client.Artifact(username, name); err == nil {
		return nil
	} else if err != atlas.ErrNotFound {
		return fmt.Errorf("error finding artifact: %s", err)
	}

	// Artifact doesn't exist, create it
	if _, err := client.CreateArtifact(username, name); err != nil {
		return fmt.Errorf("error creating artifact: %s", err)
	}
	return nil
}

func (c *ArtifactsPushCommand) pushSingleArtifact(client *atlas.Client, version *atlas.ArtifactVersion) error {
	if version.File {
		return fmt.Errorf("hatlas does not support pushing artifacts that contain files")
	}

	_, err := client.UploadArtifact(&atlas.UploadArtifactOpts{
		User:     version.User,
		Name:     version.Name,
		Type:     version.Type,
		ID:       version.ID,
		Metadata: version.Metadata,
	})
	return err
}
