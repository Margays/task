package telemetry

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	tracer        trace.Tracer
	shutdownFuncs func(context.Context) error
}

func New(ctx context.Context) (*Telemetry, error) {
	spanExporter, err := autoexport.NewSpanExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating span exporter: %w", err)
	}

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter),
	)
	otel.SetTracerProvider(tracerProvider)

	return &Telemetry{
		tracer: tracerProvider.Tracer("github.com/go-task/task/v3"),
		shutdownFuncs: func(ctx context.Context) error {
			return tracerProvider.Shutdown(ctx)
		},
	}, nil
}

func (t *Telemetry) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if t.shutdownFuncs == nil {
		return nil
	}

	var err error
	for _, shutdownFunc := range []func(context.Context) error{t.shutdownFuncs} {
		err = errors.Join(err, shutdownFunc(shutdownCtx))
	}

	t.shutdownFuncs = nil

	return err
}

func (t *Telemetry) Tracer() trace.Tracer {
	return t.tracer
}
