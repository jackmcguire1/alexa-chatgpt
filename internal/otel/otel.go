package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func SetupXrayOtel() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	idg := xray.NewIDGenerator()

	// Create and start new OTLP trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("0.0.0.0:4317"), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0)),
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithIDGenerator(idg),
	)

	//tp, err := xrayconfig.NewTracerProvider(ctx)
	//if err != nil {
	//	fmt.Printf("error creating tracer provider: %v", err)
	//}

	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})
	return tp, nil
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
