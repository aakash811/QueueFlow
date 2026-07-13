package tracing

import (
	"context"
	"os"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var Tracer = otel.Tracer("queueflow")

func InitTracer(serviceName string) func(context.Context) error {
	ctx := context.Background()

	traceExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
		otlptracegrpc.WithInsecure(),
	)

	if err != nil {
		logger.Log.Fatal(
			"failed to create OTLP trace exporter",
			zap.Error(err),
		)
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentNameKey.String(config.AppConfig.AppEnv),
			attribute.String("service.version", "v2"),
		)),
	)

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	logger.Log.Info(
		"tracing initialized",
		zap.String("service", serviceName),
		zap.String("endpoint", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")),
	)

	return tracerProvider.Shutdown
}

func GetCorrelationID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

func InjectTraceContext(ctx context.Context, headers map[string]string) {
	propagator := propagation.TraceContext{}
	carrier := propagation.MapCarrier(headers)
	propagator.Inject(ctx, carrier)
}

func ExtractTraceContext(ctx context.Context, headers map[string]string) context.Context {
	propagator := propagation.TraceContext{}
	carrier := propagation.MapCarrier(headers)
	return propagator.Extract(ctx, carrier)
}
