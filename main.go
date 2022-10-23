package main

import (
	"context"
	"flagon/command"
	"flagon/tracing"
	"flagon/version"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	appName := "flagon"
	appVersion := version.VersionNumber()

	exporter, err := tracing.ConfigFromEnvironment()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error configuring tracing: %s\n", err.Error())
		return 1
	}

	shutdown, err := tracing.Configure(context.Background(), appName, appVersion, exporter)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initialising tracing: %s\n", err.Error())
		return 1
	}

	stdOut, stdErr := configureOutput()

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      stdOut,
			ErrorWriter: stdErr,
		},
	}

	cli := &cli.CLI{
		Name:                       appName,
		Args:                       args,
		Commands:                   command.Commands(ui),
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: false,
		HelpFunc:                   cli.BasicHelpFunc(appName),
		HelpWriter:                 stdOut,
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	err = shutdown(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}

func configureOutput() (stdOut io.Writer, stdErr io.Writer) {
	stdOut = os.Stdout
	stdErr = os.Stderr

	useColor := true
	if os.Getenv("NO_COLOR") != "" || color.NoColor {
		useColor = false
	}

	if useColor {
		if f, ok := stdOut.(*os.File); ok {
			stdOut = colorable.NewColorable(f)
		}
		if f, ok := stdErr.(*os.File); ok {
			stdErr = colorable.NewColorable(f)
		}
	} else {
		stdOut = colorable.NewNonColorable(stdOut)
		stdErr = colorable.NewNonColorable(stdErr)
	}

	return stdOut, stdErr
}
