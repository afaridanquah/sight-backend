package internal

import (
	"context"

	"github.com/afaridanquah/verifylab-backend/internal"
	"github.com/afaridanquah/verifylab-backend/internal/envvar"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
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

	jaegerEndpoint, _ := conf.Get("JAEGER_ENDPOINT")

	jaegerExporter, err := otlptracehttp.New(context.Background(), otlptracehttp.WithEndpoint(jaegerEndpoint))

	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "jaeger.New")
	}

	tp := trace.NewTracerProvider(
		trace.WithSyncer(jaegerExporter),
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resource.NewSchemaless(attribute.KeyValue{
			Key:   semconv.ServiceNameKey,
			Value: attribute.StringValue("verifylab"),
		})),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return promExporter, nil
}
