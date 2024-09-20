package tracing

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type Settings struct {
	ServiceName        string
	Host               string
	Port               int
	SamplingRatio      float64
	Enabled            bool
	BatchScheduleDelay time.Duration
	MaxExportBatchSize int
	KeepAliveTime      time.Duration
	KeepAliveTimeout   time.Duration
}

type TracerShutdown func(ctx context.Context) error

func Configure(
	ctx context.Context,
	settings Settings,
) (TracerShutdown, error) {
	if !settings.Enabled {
		return func(_ context.Context) error { return nil }, nil
	}

	traceResource, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(settings.ServiceName),
		),
	)
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", settings.Host, settings.Port)),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithDialOption(
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                settings.KeepAliveTime,
				Timeout:             settings.KeepAliveTimeout,
				PermitWithoutStream: true,
			})),
	)
	if err != nil {
		return nil, err
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(
			trace.ParentBased(
				trace.TraceIDRatioBased(settings.SamplingRatio),
			),
		),
		trace.WithBatcher(
			traceExporter,
			trace.WithMaxExportBatchSize(settings.MaxExportBatchSize),
			trace.WithBatchTimeout(settings.BatchScheduleDelay),
		),
		trace.WithResource(traceResource),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		jaeger.Jaeger{},
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return traceProvider.Shutdown, nil
}
