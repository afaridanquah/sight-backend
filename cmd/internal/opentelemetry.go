package internal

import (
	"context"
	"time"

	"bitbucket.org/msafaridanquah/verifylab-service/internal"
	"bitbucket.org/msafaridanquah/verifylab-service/internal/envvar"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
)

func NewOTExporter(conf *envvar.Configuration) (*prometheus.Exporter, error) {
	promExporter, err := prometheus.New()

	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "prometheus.New")
	}

	// jaegerEndpoint, _ := conf.Get("JAEGER_ENDPOINT")

	jaegerExporter, err := otlptracehttp.New(
		context.Background(),
		// otlptracehttp.WithInsecure(),
	)

	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "jaeger.New")
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(
			jaegerExporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay*time.Millisecond),
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("verifylab-service"),
			)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return promExporter, nil
}
