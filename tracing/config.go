package tracing

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func Configure(ctx context.Context, appName string, version string, exporterConfig *ExporterConfig) (func(ctx context.Context) error, error) {

	exporter, err := createExporter(ctx, exporterConfig)
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appName),
			semconv.ServiceVersionKey.String(version),
		)),
	)

	defer otel.SetTracerProvider(tp)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		fmt.Printf("Received %s, stopping\n", s)

		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}

		os.Exit(0)
	}()

	return tp.Shutdown, nil
}
