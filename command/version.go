package command

import (
	"context"
	"flagon/version"
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/spf13/pflag"
)

type VersionCommand struct {
	Meta

	printLog bool
	short    bool
}

func (c *VersionCommand) Name() string {
	return "version"
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the version number"
}

func (c *VersionCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	flags.BoolVar(&c.printLog, "changelog", false, "print the changelog")
	flags.BoolVar(&c.short, "short", false, "show only the version, not the sha")

	return flags
}

func (c *VersionCommand) RunContext(ctx context.Context, args []string) error {

	change := version.Changelog()

	if c.short {
		c.Ui.Output(change[0].Version)
	} else {
		c.Ui.Output(fmt.Sprintf(
			"%s - %s",
			change[0].Version,
			version.VersionNumber(),
		))
	}

	if c.printLog {
		out, _ := glamour.Render(change[0].Log, "dark")
		c.Ui.Output(strings.TrimSpace(out))
	}

	return nil
}
