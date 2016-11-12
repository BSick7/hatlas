package main

import (
	"github.com/BSick7/hatlas/command"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

var Version string

func main() {
	c := cli.NewCLI("hatlas", Version)
	c.Args = os.Args[1:]
	metaPtr := &command.Meta{
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Version: Version,
	}
	meta := *metaPtr

	c.Commands = map[string]cli.CommandFactory{
		"terra": command.TerraCommandFactory(meta),
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
