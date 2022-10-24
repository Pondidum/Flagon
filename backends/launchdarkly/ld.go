package launchdarkly

import (
	"context"
	"flagon/backends"
	"flagon/tracing"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"gopkg.in/launchdarkly/go-sdk-common.v2/lduser"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
	ld "gopkg.in/launchdarkly/go-server-sdk.v5"
	"gopkg.in/launchdarkly/go-server-sdk.v5/interfaces"
)

var tr = otel.Tracer("backend.launch_darkly")

type LaunchDarklyBackend struct {
	client *ld.LDClient
}

func CreateBackend(ctx context.Context, cfg LaunchDarklyConfiguration) (*LaunchDarklyBackend, error) {
	ctx, span := tr.Start(ctx, "create_backend")
	defer span.End()

	client, err := ld.MakeClient(cfg.ApiKey, cfg.Timeout)
	if err != nil {
		return nil, tracing.Error(span, err)
	}

	client.GetDataSourceStatusProvider().WaitFor(interfaces.DataSourceStateValid, 5*time.Second)
	return &LaunchDarklyBackend{
		client: client,
	}, nil

}

func (ldb *LaunchDarklyBackend) State(ctx context.Context, flag backends.Flag, user backends.User) (bool, error) {
	ctx, span := tr.Start(ctx, "state")
	defer span.End()

	u := createUser(ctx, user)

	variation, err := ldb.client.BoolVariation(flag.Key, u, flag.DefaultValue)
	if err != nil {
		return flag.DefaultValue, err
	}

	return variation, nil
}

func createUser(ctx context.Context, user backends.User) lduser.User {
	builder := lduser.NewUserBuilder(user.Key)

	for key, value := range user.Attributes {

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
