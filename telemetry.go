package tracing

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

const (
	key         = "gin-tracing"
	name        = "https://github.com/bancodobrasil/gin-tracing"
	service     = "trace-demo"
	environment = "production"
	id          = 1
)

func New(service string, opts ...Option) gin.HandlerFunc {
	log.Println("Configuring gin-telemetry middleware...")
	cfg := configuration{}

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	if cfg.Provider == nil {
		cfg.Provider = NewJaegerProvider()
	}

	tracer := cfg.Provider.Tracer(
		name,
		trace.WithInstrumentationVersion("0.0.1"),
	)

	if cfg.Propagators == nil {
		cfg.Propagators = otel.GetTextMapPropagator()
	}
	log.Println("Gin-telemetry successfully configured!")
	return func(c *gin.Context) {
		c.Set(key, tracer)
		tracedCtx := c.Request.Context()

		defer func() {
			c.Request = c.Request.WithContext(tracedCtx)
		}()

		ctx := cfg.Propagators.Extract(tracedCtx, propagation.HeaderCarrier(c.Request.Header))
		opts := []trace.SpanStartOption{
			trace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", c.Request)...),
			trace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(c.Request)...),
			trace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest(service, c.FullPath(), c.Request)...),
			trace.WithSpanKind(trace.SpanKindServer),
		}
		spanName := c.FullPath()
		if spanName == "" {
			spanName = fmt.Sprintf("HTTP %s route not found", c.Request.Method)
		}
		ctx, span := tracer.Start(ctx, spanName, opts...)
		defer span.End()

		c.Request = c.Request.WithContext(tracedCtx)
		c.Next()

		status := c.Writer.Status()
		attrs := semconv.HTTPAttributesFromHTTPStatusCode(status)
		spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
		span.SetAttributes(attrs...)
		span.SetStatus(spanStatus, spanMessage)
		if len(c.Errors) > 0 {
			span.SetAttributes(attribute.String("gin.errors", c.Errors.String()))
		}
	}
}
