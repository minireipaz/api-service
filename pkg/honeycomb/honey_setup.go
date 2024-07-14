package honeycomb

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func SetupHoneyComb(ctx context.Context) (*sdktrace.TracerProvider, *otlptrace.Exporter, error) {
	tp, exp, err := InitializeGlobalTracerProvider(ctx)
	if err != nil {
		log.Panicf("ERROR | Failed to initialize OpenTelemetry: %v", err)
	}
	return tp, exp, err
}
