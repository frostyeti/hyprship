package telemetry

import (
	"context"
	"fmt"
	"strings"

	"github.com/frostyeti/hyprship/apps/api/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// InitOtel initializes OpenTelemetry providers based on configuration.
// It returns a shutdown function that should be called when the application exits.
func InitOtel(ctx context.Context, cfg *config.Config) (func(context.Context) error, error) {
	if cfg.Otel == nil || !cfg.Otel.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("hyprship-api"),
			semconv.DeploymentEnvironment(cfg.Env),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var traceExporter trace.SpanExporter
	var hasTraceExporter bool

	// Find the first valid trace exporter
	for _, exp := range cfg.Otel.Exporters {
		switch exp {
		case "stdout":
			traceExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
			hasTraceExporter = err == nil
		case "otlp":
			traceExporter, err = otlptracehttp.New(ctx)
			hasTraceExporter = err == nil
		}
		if hasTraceExporter {
			break
		}
	}

	var tp *trace.TracerProvider
	if hasTraceExporter {
		// Use ParentBased(AlwaysSample) for dev, or consider TraceIDRatioBased for prod
		sampler := trace.ParentBased(trace.AlwaysSample())
		if cfg.Env == "production" || cfg.Env == "prod" {
			sampler = trace.ParentBased(trace.TraceIDRatioBased(0.1)) // 10% sampling
		}

		tp = trace.NewTracerProvider(
			trace.WithSampler(sampler),
			trace.WithBatcher(traceExporter),
			trace.WithResource(res),
		)
		otel.SetTracerProvider(tp)
	}

	var meterProvider *metric.MeterProvider
	var hasMetricExporter bool

	for _, exp := range cfg.Otel.Exporters {
		if exp == "prometheus" {
			metricExporter, err := prometheus.New()
			if err == nil {
				meterProvider = metric.NewMeterProvider(
					metric.WithReader(metricExporter),
					metric.WithResource(res),
				)
				hasMetricExporter = true
				break
			}
		}
		// NOTE: if otlp metrics are desired, add otlpmetrichttp setup here
	}

	if hasMetricExporter {
		otel.SetMeterProvider(meterProvider)
	}

	return func(ctx context.Context) error {
		var errs []string
		if tp != nil {
			if err := tp.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("failed to shutdown tracer provider: %v", err))
			}
		}
		if meterProvider != nil {
			if err := meterProvider.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Sprintf("failed to shutdown meter provider: %v", err))
			}
		}

		if len(errs) > 0 {
			return fmt.Errorf("otel shutdown errors: %v", strings.Join(errs, ", "))
		}
		return nil
	}, nil
}
