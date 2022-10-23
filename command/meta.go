package command

import (
	"bufio"
	"context"
	"flagon/tracing"
	"io"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type Meta struct {
	Ui cli.Ui

	tr  trace.Tracer
	cmd NamedCommand
}

func NewMeta(ui cli.Ui, cmd NamedCommand) Meta {
	return Meta{
		Ui:  ui,
		cmd: cmd,
		tr:  otel.Tracer(cmd.Name()),
	}
}

type NamedCommand interface {
	Name() string
	Synopsis() string

	Flags() *pflag.FlagSet
	RunContext(ctx context.Context, args []string) error
}

func (m *Meta) Flags(c NamedCommand) *pflag.FlagSet {
	f := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
	f.Usage = func() { m.Ui.Output(m.Help()) }

	// add common flags etc here

	// Create an io.Writer that writes to our UI properly for errors.
	// This is kind of a hack, but it does the job. Basically: create
	// a pipe, use a scanner to break it into lines, and output each line
	// to the UI. Do this forever.
	errR, errW := io.Pipe()
	errScanner := bufio.NewScanner(errR)
	go func() {
		for errScanner.Scan() {
			m.Ui.Error(errScanner.Text())
		}
	}()
	f.SetOutput(errW)

	return f
}

func (m *Meta) AutocompleteFlags() complete.Flags {
	// return m.cmd.Flags().Autocomplete()
	return nil
}

func (m *Meta) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (m *Meta) Help() string {
	return m.cmd.Synopsis() + "\n\n" + m.cmd.Flags().FlagUsages()
}

func (m *Meta) Run(args []string) int {
	ctx := context.Background()

	ctx, span := m.tr.Start(ctx, m.cmd.Name())
	defer span.End()

	f := m.cmd.Flags()

	if err := f.Parse(args); err != nil {
		tracing.Error(span, err)
		m.Ui.Error(err.Error())

		return 1
	}

	tracing.StoreFlags(ctx, f)

	if err := m.cmd.RunContext(ctx, f.Args()); err != nil {
		tracing.Error(span, err)
		m.Ui.Error(err.Error())

		return 1
	}

	return 0
}
