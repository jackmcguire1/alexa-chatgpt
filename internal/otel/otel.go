package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Common resource for all providers
func newResource() (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("alexa-chatgpt"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)
}

func SetupXrayOtel() *sdktrace.TracerProvider {
	ctx := context.Background()

	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}

	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
	return tp
}

func GetXRayTraceID(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ""
	}

	traceID := spanCtx.TraceID().String() // 32-char hex
	// Convert to AWS X-Ray trace ID: 1-<8_hex_digits>-<24_hex_digits>
	return fmt.Sprintf("1-%s-%s", traceID[0:8], traceID[8:])
}
