// Package otel provides otel support.
package otel

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/msafaridanquah/sight-backend/foundation/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// Config defines the information needed to init tracing.
type Config struct {
	ServiceName    string
	Host           string
	ExcludedRoutes map[string]struct{}
	Probability    float64
}

// InitTracing configures open telemetry to be used with the service.
func InitTracing(log *logger.Logger, cfg Config) (*sdktrace.TracerProvider, func(ctx context.Context), error) {

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			otlptracegrpc.WithInsecure(),
			otlptracegrpc.WithEndpoint(cfg.Host),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("creating new exporter: %w", err)
	}

	// var traceProvider *trace.TracerProvider

	tp := sdktrace.NewTracerProvider(
		// sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithBatchTimeout(sdktrace.DefaultScheduleDelay*time.Millisecond),
		),
		sdktrace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(cfg.ServiceName),
			),
		),
	)

	teardown := func(ctx context.Context) {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Info(ctx, "Error shutting down tracer provider: %v", err)
		}
	}

	// traceProvider = tp

	// We must set this provider as the global provider for things to work,
	// but we pass this provider around the program where needed to collect
	// our traces.
	otel.SetTracerProvider(tp)

	// Extract incoming trace contexts and the headers we set in outgoing requests.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, teardown, nil
}

// InjectTracing initializes the request for tracing by writing otel related
// information into the response and saving the tracer and trace id in the
// context for later use.
func InjectTracing(ctx context.Context, tracer trace.Tracer) context.Context {
	ctx = setTracer(ctx, tracer)

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()
	if traceID == "00000000000000000000000000000000" {
		traceID = uuid.NewString()
	}
	ctx = setTraceID(ctx, traceID)

	return ctx
}

// AddSpan adds an otel span to the existing trace.
func AddSpan(ctx context.Context, spanName string, keyValues ...attribute.KeyValue) (context.Context, trace.Span) {
	v, ok := ctx.Value(tracerKey).(trace.Tracer)

	if !ok || v == nil {
		return ctx, trace.SpanFromContext(ctx)
	}

	ctx, span := v.Start(ctx, spanName)
	span.SetAttributes(keyValues...)

	return ctx, span
}

// AddTraceToRequest adds the current trace id to the request so it
// can be delivered to the service being called.
func AddTraceToRequest(ctx context.Context, r *http.Request) {
	hc := propagation.HeaderCarrier(r.Header)
	otel.GetTextMapPropagator().Inject(ctx, hc)
}

func NewOTELSpan(ctx context.Context, otelName, name string, keyValues ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := otel.Tracer(otelName).Start(ctx, name)

	span.SetAttributes(keyValues...)

	return ctx, span
}
