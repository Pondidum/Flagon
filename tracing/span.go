package tracing

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/spf13/pflag"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func FromMap[V any](prefix string, m map[string]V) []attribute.KeyValue {

	attrs := make([]attribute.KeyValue, 0, len(m))

	for k, v := range m {
		attrs = append(attrs, asAttribute(k, v))
	}

	return attrs
}

func asAttribute(key string, v any) attribute.KeyValue {

	switch val := v.(type) {

	case string:
		return attribute.String(key, val)

	case *string:
		return attribute.String(key, *val)

	case bool:
		return attribute.Bool(key, val)

	case *bool:
		return attribute.Bool(key, *val)

	case int:
		return attribute.Int(key, val)

	case *int:
		return attribute.Int(key, *val)

	case int64:
		return attribute.Int64(key, val)

	case *int64:
		return attribute.Int64(key, *val)

	case *time.Time:
		return attribute.Int64(key, val.UnixMilli())

	case time.Time:
		return attribute.Int64(key, val.UnixMilli())

	default:
		// handle the pointer case, so that `Sprint` prints the actual value
		if reflect.ValueOf(v).Kind() == reflect.Pointer {
			v = reflect.Indirect(reflect.ValueOf(v))
		}

		return attribute.String(key, fmt.Sprint(v))

	}
}

func StoreFlags(ctx context.Context, flags *pflag.FlagSet) {
	s := trace.SpanFromContext(ctx)

	flags.VisitAll(func(f *pflag.Flag) {
		s.SetAttributes(attribute.String("flags."+f.Name, f.Value.String()))
	})
}

func Errorf(s trace.Span, format string, a ...interface{}) error {
	return Error(s, fmt.Errorf(format, a...))
}

func Error(s trace.Span, err error) error {
	s.RecordError(err)
	s.SetStatus(codes.Error, err.Error())

	return err
}
