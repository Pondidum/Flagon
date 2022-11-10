package command

import (
	"context"
	"flagon/backends"
	"flagon/backends/launchdarkly"
	"flagon/tracing"
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
)

type StateCommand struct {
	Meta

	backend string
	output  string

	userKey        string
	userAttributes []string
}

func (c *StateCommand) Name() string {
	return "state"
}

func (c *StateCommand) Synopsis() string {
	return "Checks the state of a feature flag"
}

func (c *StateCommand) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet(c.Name(), pflag.ContinueOnError)

	flags.StringVar(&c.backend, "backend", "launchdarkly", "which flag service to use")
	flags.StringVar(&c.output, "output", "json", "specifies the output format")

	flags.StringVar(&c.userKey, "user", "", "The key/id of the user to query a flag against")
	flags.StringSliceVar(&c.userAttributes, "attr", []string{}, "key=value pairs of additional properties for the user")

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

	defaultValue, err := strconv.ParseBool(args[1])
	if err != nil {
		return tracing.Error(span, err)
	}

	flag := backends.Flag{
		Key:          args[0],
		DefaultValue: defaultValue,
	}

	attrs, err := parseKeyValuePairs(c.userAttributes)
	if err != nil {
		return tracing.Error(span, err)
	}

	user := backends.User{
		Key:        c.userKey,
		Attributes: attrs,
	}

	value, err := backend.State(ctx, flag, user)
	if err != nil {
		return tracing.Error(span, err)
	}

	return print(c.Ui, c.output, map[string]interface{}{
		"flag":    flag.Key,
		"default": flag.DefaultValue,
		"state":   value,
	})
}

func (c *StateCommand) createBackend(ctx context.Context) (backends.Backend, error) {
	ctx, span := c.tr.Start(ctx, "create_backend")
	defer span.End()

	switch c.backend {
	case "launchdarkly":

		cfg, err := launchdarkly.ConfigFromEnvironment()
		if err != nil {
			return nil, tracing.Error(span, err)
		}

		return launchdarkly.CreateBackend(ctx, cfg)

	default:
		return nil, fmt.Errorf("unsupported backend: %s", c.backend)
	}
}
