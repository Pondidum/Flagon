package command

import (
	"github.com/mitchellh/cli"
)

func Commands(ui cli.Ui) map[string]cli.CommandFactory {

	return map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			cmd := &VersionCommand{}
			cmd.Meta = NewMeta(ui, cmd)

			return cmd, nil
		},
	}
}
