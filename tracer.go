package logger

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// reference: https://github.com/vincentfree/opentelemetry/blob/main/otelzerolog/example_test.go

type (
	TracerLogger struct {
		trace.Span
		trace.Tracer
	}
)

func NewLoggerWithTraceID(ctx context.Context, componentName, spanName string) *TracerLogger {
	var (
		span   trace.Span
		tracer trace.Tracer
	)

	tracer = otel.Tracer(componentName)

	// Check if the context has a valid trace ID
	if !trace.SpanContextFromContext(ctx).IsValid() {
		// If not, start a new span
		ctx, _ = tracer.Start(ctx, spanName)
	}

	// Start a new span
	ctx, span = tracer.Start(ctx, spanName)

	// Use the span context to generate a trace ID
	traceID := span.SpanContext().TraceID()

	// Log the trace ID
	fmt.Println("Trace ID:", traceID)

	return &TracerLogger{
		Span:   span,
		Tracer: tracer,
	}
}

func (t *TracerLogger) Info(sender, format string, v ...any) {
	t.logWithTraceID(LevelInfo, sender, otelzerolog.AddTracingContext(t.Span), format, v...)
}

func (t *TracerLogger) InfoWithAttributes(sender string, attributes []attribute.KeyValue, format string, v ...any) {
	t.logWithTraceID(LevelInfo, sender, otelzerolog.AddTracingContextWithAttributes(t.Span, attributes), format, v...)
}

func (t *TracerLogger) Error(sender, format string, v ...any) {
	t.logWithTraceID(LevelError, sender, otelzerolog.AddTracingContext(t.Span), format, v...)
}

func (t *TracerLogger) ErrorWithAttributes(sender string, attributes []attribute.KeyValue, format string, v ...any) {
	t.logWithTraceID(LevelError, sender, otelzerolog.AddTracingContextWithAttributes(t.Span, attributes), format, v...)
}

func (t *TracerLogger) Warn(sender, format string, v ...any) {
	t.logWithTraceID(LevelWarn, sender, otelzerolog.AddTracingContext(t.Span), format, v...)
}

func (t *TracerLogger) WarnWithAttributes(sender string, attributes []attribute.KeyValue, format string, v ...any) {
	t.logWithTraceID(LevelWarn, sender, otelzerolog.AddTracingContextWithAttributes(t.Span, attributes), format, v...)
}

func (t *TracerLogger) Debug(sender, format string, v ...any) {
	t.logWithTraceID(LevelDebug, sender, otelzerolog.AddTracingContext(t.Span), format, v...)
}

func (t *TracerLogger) DebugWithAttributes(sender string, attributes []attribute.KeyValue, format string, v ...any) {
	t.logWithTraceID(LevelDebug, sender, otelzerolog.AddTracingContextWithAttributes(t.Span, attributes), format, v...)
}

// Close closes the span
func (t *TracerLogger) Close() {
	t.Span.End()
}

// LogWithTraceID logs at the specified level for the specified sender
func (t *TracerLogger) logWithTraceID(level LogLevel, sender string, fn func(e *zerolog.Event), format string, v ...any) {
	var ev *zerolog.Event
	switch level {
	case LevelDebug:
		ev = global.Debug()
	case LevelInfo:
		ev = global.Info()
	case LevelWarn:
		ev = global.Warn()
	default:
		ev = global.Error()
	}

	ev.Timestamp().Str("sender", sender)

	fn(ev)

	ev.Msg(fmt.Sprintf(format, v...))
}
