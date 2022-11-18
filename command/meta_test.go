package command

import (
	"flagon/backends"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
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
