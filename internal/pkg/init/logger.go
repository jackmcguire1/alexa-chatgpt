package init

import (
	"context"
	"log/slog"
	"os"

	otelsetup "github.com/jackmcguire1/alexa-chatgpt/internal/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

// SetupLogger creates and returns a JSON logger
func SetupLogger() *slog.Logger {
	jsonLogH := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	return slog.New(jsonLogH)
}

// SetupTracing initializes and returns the OpenTelemetry tracer
func SetupTracing(ctx context.Context, logger *slog.Logger) *trace.TracerProvider {
	tracer, err := otelsetup.SetupXrayOtel(ctx)
	if err != nil {
		logger.With("error", err).Error("failed to setup tracer")
		panic(err)
	}
	return tracer
}