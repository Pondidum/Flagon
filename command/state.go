package command

import (
	"bufio"
	"context"
	"flagon/backends"
	"flagon/tracing"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/mitchellh/cli"
	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/attribute"
)

func NewStateCommand(ui cli.Ui) (*StateCommand, error) {
	cmd := &StateCommand{
		readFile: func(f string) (io.ReadCloser, error) {
			return os.Open(f)
		},
	}
	cmd.Meta = NewMeta(ui, cmd)

	return cmd, nil
}

type StateCommand struct {
	Meta

	userKey        string
	userAttributes []string

	userAttributesFile string

	readFile func(filePath string) (io.ReadCloser, error)
}

func (c *StateCommand) Name() string {
	return "state"
}

func (c *StateCommand) Synopsis() string {
	return "Checks the state of a feature flag"
}

func (c *StateCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)

	flags.StringVar(&c.userKey, "user", "", "The key/id of the user to query a flag against")
	flags.StringSliceVar(&c.userAttributes, "attr", []string{}, "key=value pairs of additional properties for the user")
	flags.StringVar(&c.userAttributesFile, "attr-file", "flagon.attrs", "a file containing additional properties for the user")

	return flags
}

func (c *StateCommand) RunContext(ctx context.Context, args []string) error {
	ctx, span := c.tr.Start(ctx, "run")
	defer span.End()

	if len(args) < 1 || len(args) > 2 {
		return fmt.Errorf("this command takes one to two arguments: flagKey and flagDefault")
	}

	backend, err := c.createBackend(ctx)
	if err != nil {
		return tracing.Error(span, err)
	}
	defer backend.Close(ctx)

	flag := backends.Flag{
		Key: args[0],
	}

	if len(args) > 1 {
		defaultValue, err := strconv.ParseBool(args[1])
		if err != nil {
			return tracing.Error(span, err)
		}

		flag.DefaultValue = defaultValue
	}

	span.SetAttributes(
		attribute.String("flag.key", flag.Key),
		attribute.Bool("flag.default", flag.DefaultValue),
	)

	lines := []string{}
	if c.userAttributesFile != "" {
		f, err := c.readFile(c.userAttributesFile)
		if err == nil {
			defer f.Close()
			s := bufio.NewScanner(f)
			for s.Scan() {
				lines = append(lines, s.Text())
			}
		}
	}

	attrs, err := parseKeyValuePairs(append(lines, c.userAttributes...))
	if err != nil {
		return tracing.Error(span, err)
	}

	if key, found := attrs["user-key"]; found {
		delete(attrs, "user-key")
		if c.userKey == "" {
			c.userKey = key
		}
	}

	user := backends.User{
		Key:        c.userKey,
		Attributes: attrs,
	}
	span.SetAttributes(attribute.String("user.key", user.Key))
	span.SetAttributes(tracing.FromMap("user.", user.Attributes)...)

	if flag, err = backend.State(ctx, flag, user); err != nil {
		return tracing.Error(span, err)
	}

	if err := c.print(flag); err != nil {
		return tracing.Error(span, err)
	}

	span.SetAttributes(attribute.Bool("flag.value", flag.Value))

	if flag.Value {
		return nil
	}

	return &SilentError{}
}
