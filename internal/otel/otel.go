package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func SetupXrayOtel(ctx context.Context) (*sdktrace.TracerProvider, error) {
	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		return nil, err
	}

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
