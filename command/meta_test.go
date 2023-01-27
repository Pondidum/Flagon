package command

import (
	"context"
	"flagon/backends"
	"flagon/tracing"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestPrinting(t *testing.T) {

	input := backends.Flag{
		Key:          "the-flag-key",
		DefaultValue: false,
		Value:        true,
	}

	ui := cli.NewMockUi()
	m := NewMeta(ui, &VersionCommand{})

	t.Run("Json", func(t *testing.T) {

		m.output = "json"
		ui.OutputWriter.Reset()
		ui.ErrorWriter.Reset()

		assert.NoError(t, m.print(input))
		assert.Equal(t,
			`{"key":"the-flag-key","defaultValue":false,"value":true}`,
			strings.TrimSpace(ui.OutputWriter.String()),
		)

	})

	t.Run("Template - Good", func(t *testing.T) {

		m.output = "template={{.Value}}"
		ui.OutputWriter.Reset()
		ui.ErrorWriter.Reset()

		assert.NoError(t, m.print(input))
		assert.Equal(t,
			`true`,
			strings.TrimSpace(ui.OutputWriter.String()),
		)

	})

	t.Run("Template - Bad Casing", func(t *testing.T) {

		m.output = "template={{.value}}"
		ui.OutputWriter.Reset()
		ui.ErrorWriter.Reset()

		assert.ErrorContains(t, m.print(input), "<.value>")
	})
}

func TestTraceParent(t *testing.T) {

	exporter := tracing.NewMemoryExporter()

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSyncer(exporter),
	)

	otel.SetTextMapPropagator(propagation.TraceContext{})

	m := &Meta{
		tr:  tp.Tracer("mock"),
		cmd: &MockCommand{},
	}

	os.Setenv(TraceParentEnvVar, "00-7107538ee3f6bc77ada1b2d34a412e1d-bfe6177cefb76eb2-01")
	assert.Equal(t, 0, m.Run([]string{}))

	assert.Equal(t, "7107538ee3f6bc77ada1b2d34a412e1d", exporter.Spans[0].SpanContext().TraceID().String())
}

func NewMockCommand(ui cli.Ui) *MockCommand {
	cmd := &MockCommand{}
	cmd.Meta = NewMeta(ui, cmd)
	return cmd
}

type MockCommand struct {
	Meta
}

func (c *MockCommand) Name() string {
	return "mock"
}
func (c *MockCommand) Synopsis() string {
	return "mock"
}
func (c *MockCommand) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)
}

func (c *MockCommand) RunContext(ctx context.Context, args []string) error {
	return nil
}
