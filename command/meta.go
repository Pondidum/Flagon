package command

import (
	"context"
	"flagon/backends"
	"flagon/backends/launchdarkly"
	"flagon/tracing"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/posener/complete"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Meta struct {
	Ui cli.Ui

	tr  trace.Tracer
	cmd NamedCommand

	backend string
	output  string

	ldFlags launchdarkly.LaunchDarklyConfiguration
}

type NamedCommand interface {
	Name() string
	Synopsis() string

	Flags() *pflag.FlagSet
	RunContext(ctx context.Context, args []string) error
}

type FlagGroup struct {
	*pflag.FlagSet
	Name string
}

func NewMeta(ui cli.Ui, cmd NamedCommand) Meta {
	return Meta{
		Ui:  ui,
		cmd: cmd,
		tr:  otel.Tracer(cmd.Name()),

		ldFlags: launchdarkly.LaunchDarklyConfiguration{},
	}
}

func (m *Meta) AutocompleteFlags() complete.Flags {
	// return m.cmd.Flags().Autocomplete()
	return nil
}

func (m *Meta) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (m *Meta) Help() string {
	sb := strings.Builder{}

	sb.WriteString(m.cmd.Synopsis())
	sb.WriteString("\n\n")

	for _, group := range m.allFlags() {
		sb.WriteString(group.Name)
		sb.WriteString(" flags")
		sb.WriteString("\n\n")
		sb.WriteString(group.FlagUsages())
		sb.WriteString("\n")
	}

	return sb.String()
}

func combineFlags(groups []FlagGroup) *pflag.FlagSet {

	if len(groups) == 0 {
		return nil
	}

	flags := pflag.NewFlagSet(groups[0].Name, pflag.ContinueOnError)

	for _, f := range groups {
		flags.AddFlagSet(f.FlagSet)
	}

	return flags
}

func newFlagGroup(name string) FlagGroup {
	return FlagGroup{
		FlagSet: pflag.NewFlagSet(name, pflag.ContinueOnError),
		Name:    name,
	}
}

func (m *Meta) allFlags() []FlagGroup {

	common := newFlagGroup("Common")

	common.StringVar(&m.backend, "backend", "launchdarkly", "which flag service to use")
	common.StringVar(&m.output, "output", "json", "specifies the output format")

	return []FlagGroup{
		{Name: "Command", FlagSet: m.cmd.Flags()},
		common,
		{Name: "LaunchDarkly Backend", FlagSet: m.ldFlags.Flags()},
		// other backend flags here
	}
}

func (m *Meta) createBackend(ctx context.Context) (backends.Backend, error) {
	ctx, span := m.tr.Start(ctx, "create_backend")
	defer span.End()

	span.SetAttributes(attribute.String("backend", m.backend))

	switch m.backend {
	case "launchdarkly":

		cfg := launchdarkly.DefaultConfig()
		cfg.OverrideFrom(launchdarkly.ConfigFromEnvironment())
		cfg.OverrideFrom(m.ldFlags)

		return launchdarkly.CreateBackend(ctx, cfg)

	default:
		return nil, fmt.Errorf("unsupported backend: %s", m.backend)
	}
}

func (m *Meta) Run(args []string) int {
	ctx := context.Background()

	ctx, span := m.tr.Start(ctx, m.cmd.Name())
	defer span.End()

	f := combineFlags(m.allFlags())

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
