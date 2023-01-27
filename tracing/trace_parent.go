package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
)

func WithTraceParent(ctx context.Context, parent string) context.Context {
	if parent == "" {
		return ctx
	}

	// https://github.com/open-telemetry/opentelemetry-go/blob/main/propagation/trace_context.go#29
	// the 'traceparent' key is a private constant in the otel library so this
	// is using an internal detail but it's probably fine
	carrier := NewCliCarrier()
	carrier.Set("traceparent", parent)

	prop := otel.GetTextMapPropagator()
	return prop.Extract(ctx, carrier)
}

type CliCarrier struct {
	data map[string]string
}

func NewCliCarrier() *CliCarrier {
	return &CliCarrier{
		data: map[string]string{},
	}
}

func (c *CliCarrier) Keys() []string {
	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

func (c *CliCarrier) Get(key string) string {
	val, found := c.data[key]
	if found {
		return val
	}

	return ""
}

func (c *CliCarrier) Set(key string, value string) {
	c.data[key] = value
}
