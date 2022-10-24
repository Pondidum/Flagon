package tracing

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otlphttp "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const OtlpEndpointEnvVar = "OTEL_EXPORTER_OTLP_ENDPOINT"
const OtlpTracesEndpointEnvVar = "OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"
const OtlpHeaders = "OTEL_EXPORTER_OTLP_HEADERS"
const OtelDebugEnvVar = "OTEL_DEBUG"

const OtelTraceExporterEnvVar = "OTEL_TRACE_EXPORTER"

type ExporterConfig struct {
	Endpoint string
	Headers  map[string]string
	Debug    bool

	ExporterType string
}

func DefaultConfig() *ExporterConfig {
	return &ExporterConfig{
		Endpoint:     "localhost:4317",
		Headers:      map[string]string{},
		Debug:        false,
		ExporterType: "",
	}
}
func ConfigFromEnvironment() (*ExporterConfig, error) {

	config := DefaultConfig()

	if val := os.Getenv(OtlpTracesEndpointEnvVar); val != "" {
		config.Endpoint = val
	} else if val := os.Getenv(OtlpEndpointEnvVar); val != "" {
		config.Endpoint = val
	}

	if val, err := strconv.ParseBool(os.Getenv(OtelDebugEnvVar)); err == nil {
		config.Debug = val
	}

	if val := os.Getenv(OtlpHeaders); val != "" {
		for _, pair := range strings.Split(val, ",") {
			index := strings.Index(pair, ":")
			if index == -1 {
				return nil, fmt.Errorf("Unable to parse '%s' as a key:value pair, missing a ':'", pair)
			}

			key := strings.TrimSpace(pair[0:index])
			val := strings.TrimSpace(pair[index+1:])

			config.Headers[key] = val
		}
	}

	if val := os.Getenv(OtelTraceExporterEnvVar); val != "" {
		config.ExporterType = val
	}

	return config, nil
}

func createExporter(ctx context.Context, conf *ExporterConfig) (sdktrace.SpanExporter, error) {

	if conf.ExporterType == "" {
		return &NoopExporter{}, nil
	}

	if conf.ExporterType == "stdout" {
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	if conf.ExporterType == "stderr" {
		return stdouttrace.New(stdouttrace.WithPrettyPrint(), stdouttrace.WithWriter(os.Stderr))
	}

	endpoint := strings.ToLower(conf.Endpoint)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(endpoint, "https://") || strings.HasPrefix(endpoint, "http://") {

		opts := []otlphttp.Option{}

		hostAndPort := u.Host
		if u.Port() == "" {
			if u.Scheme == "https" {
				hostAndPort += ":443"
			} else {
				hostAndPort += ":80"
			}
		}
		opts = append(opts, otlphttp.WithEndpoint(hostAndPort))

		if u.Path == "" {
			u.Path = "/v1/traces"
		}
		opts = append(opts, otlphttp.WithURLPath(u.Path))

		if u.Scheme == "http" {
			opts = append(opts, otlphttp.WithInsecure())
		}

		opts = append(opts, otlphttp.WithHeaders(conf.Headers))

		return otlphttp.New(ctx, opts...)
	} else {
		opts := []otlpgrpc.Option{}

		opts = append(opts, otlpgrpc.WithEndpoint(endpoint))

		isLocal, err := isLoopbackAddress(endpoint)
		if err != nil {
			return nil, err
		}

		if isLocal {
			opts = append(opts, otlpgrpc.WithInsecure())
		}

		opts = append(opts, otlpgrpc.WithHeaders(conf.Headers))

		return otlpgrpc.New(ctx, opts...)
	}

}

func isLoopbackAddress(endpoint string) (bool, error) {
	hpRe := regexp.MustCompile(`^[\w.-]+:\d+$`)
	uriRe := regexp.MustCompile(`^(http|https)`)

	endpoint = strings.TrimSpace(endpoint)

	var hostname string
	if hpRe.MatchString(endpoint) {
		parts := strings.SplitN(endpoint, ":", 2)
		hostname = parts[0]
	} else if uriRe.MatchString(endpoint) {
		u, err := url.Parse(endpoint)
		if err != nil {
			return false, err
		}
		hostname = u.Hostname()
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		return false, err
	}

	allAreLoopback := true
	for _, ip := range ips {
		if !ip.IsLoopback() {
			allAreLoopback = false
		}
	}

	return allAreLoopback, nil
}

type NoopExporter struct{}

func (nsb *NoopExporter) ExportSpans(context.Context, []sdktrace.ReadOnlySpan) error { return nil }
func (nsb *NoopExporter) Shutdown(context.Context) error                             { return nil }
