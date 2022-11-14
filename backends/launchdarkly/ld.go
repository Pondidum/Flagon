package launchdarkly

import (
	"context"
	"flagon/backends"
	"flagon/tracing"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/interfaces"
	"gopkg.in/launchdarkly/go-server-sdk.v5/ldcomponents"
)

var tr = otel.Tracer("backend.launch_darkly")

type LaunchDarklyBackend struct {
	client *ld.LDClient
}

func CreateBackend(ctx context.Context, cfg LaunchDarklyConfiguration) (*LaunchDarklyBackend, error) {
	ctx, span := tr.Start(ctx, "create_backend")
	defer span.End()

	ldConfig := ld.Config{
		DiagnosticOptOut: true,
	}

	if cfg.Debug {
		ldConfig.Logging = ldcomponents.Logging()
	} else {
		ldConfig.Logging = ldcomponents.NoLogging()
	}

	client, err := ld.MakeCustomClient(cfg.SdkKey, ldConfig, cfg.Timeout)
	if err != nil {
		return nil, tracing.Error(span, err)
	}
	client.GetDataSourceStatusProvider().WaitFor(interfaces.DataSourceStateValid, 5*time.Second)
	return &LaunchDarklyBackend{
		client: client,
	}, nil

}

func (ldb *LaunchDarklyBackend) Close(ctx context.Context) error {
	_, span := tr.Start(ctx, "close")
	defer span.End()

	return ldb.client.Close()
}

func (ldb *LaunchDarklyBackend) State(ctx context.Context, flag backends.Flag, user backends.User) (backends.Flag, error) {
	ctx, span := tr.Start(ctx, "state")
	defer span.End()

	u := createUser(ctx, user)

	span.SetAttributes(attribute.String("flag.key", flag.Key))

	flag.Value = flag.DefaultValue

	variation, detail, err := ldb.client.BoolVariationDetail(flag.Key, u, flag.DefaultValue)
	if err != nil {
		return flag, tracing.Error(span, err)
	}

	span.SetAttributes(attribute.String("reason", detail.Reason.String()))
	span.SetAttributes(attribute.Bool("variation", variation))

	flag.Value = variation

	return flag, nil
}

func createUser(ctx context.Context, user backends.User) lduser.User {
	ctx, span := tr.Start(ctx, "create_user")
	defer span.End()

	span.SetAttributes(attribute.String("attr.key", user.Key))

	builder := lduser.NewUserBuilder(user.Key)

	for key, value := range user.Attributes {

		span.SetAttributes(attribute.String("attr."+key, value))
		cleanKey := strings.ToLower(strings.ReplaceAll(key, "_", ""))

		switch cleanKey {
		case "name":
			builder.Name(value)

		case "firstname":
			builder.FirstName(value)

		case "lastname":
			builder.LastName(value)

		case "email":
			builder.Email(value)

		case "country":
			builder.Country(value)

		case "ip":
			builder.IP(value)

		case "secondary":
			builder.Secondary(value)

		default:
			// note, this is the key as passed in, not the cleankey used for switch
			builder.Custom(key, ldvalue.String(value))
		}
	}

	return builder.Build()
}
