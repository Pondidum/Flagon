package command

import (
	"context"
	"flagon/backends"
	"flagon/tracing"
	"fmt"
	"strconv"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/attribute"
)

type StateCommand struct {
	Meta

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
	span.SetAttributes(
		attribute.String("flag.key", flag.Key),
		attribute.Bool("flag.default", flag.DefaultValue),
	)

	attrs, err := parseKeyValuePairs(c.userAttributes)
	if err != nil {
		return tracing.Error(span, err)
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

	return nil
}
