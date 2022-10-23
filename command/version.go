package command

import (
	"context"
	"flagon/version"

	"github.com/spf13/pflag"
)

type VersionCommand struct {
	Meta
}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the version number"
}

func (c *VersionCommand) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
}

func (c *VersionCommand) RunContext(ctx context.Context, args []string) error {
	c.Ui.Output(version.VersionNumber())
	return nil
}
