package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

type keyContext string

const (
	keyTracer     keyContext = "gin-telemetry-trace"
	keyPropagator keyContext = "gin-telemetry-propagator"
	name          string     = "https://github.com/bancodobrasil/gin-telemetry"
	environment   string     = "production"
	id            int64      = 1
)

var (
	service string
)

// Middleware ...
func Middleware(serviceName string, opts ...Option) gin.HandlerFunc {
	log.Info("Configuring gin-telemetry middleware ...")

	if serviceName == "" {
		hostname, _ := os.Hostname()
		service = fmt.Sprintf("service-%s", hostname)
	} else {
		service = serviceName
	}

	cfg := configuration{}

	for _, opt := range opts {
		opt.apply(&cfg)
	}

	// setting default provider to Jaeger
	if cfg.Provider == nil {
		cfg.Provider = NewJaegerProvider()
	}

	tracer := cfg.Provider.Tracer(
		name,
		trace.WithInstrumentationVersion("0.0.1"),
	)

	if cfg.Propagator == nil {
		cfg.Propagator = getDefaultTextMapPropagator()
	}

	log.Info("Gin-telemetry successfully configured!")
	return func(c *gin.Context) {
		tracedCtx := c.Request.Context()
		tracedCtx = context.WithValue(tracedCtx, keyTracer, tracer)
		tracedCtx = context.WithValue(tracedCtx, keyPropagator, cfg.Propagator)
		defer func() {
			c.Request = c.Request.WithContext(tracedCtx)
		}()

		ctx := cfg.Propagator.Extract(tracedCtx, propagation.HeaderCarrier(c.Request.Header))
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

		c.Request = c.Request.WithContext(ctx)
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

// GetTracer ...
func GetTracer(ctx context.Context) trace.Tracer {
	tracerInterface := ctx.Value(keyTracer)
	return tracerInterface.(trace.Tracer)
}

// Inject ...
func Inject(ctx context.Context, headers http.Header) {
	propagatorInteface := ctx.Value(keyPropagator)
	propagator := propagatorInteface.(propagation.TextMapPropagator)
	propagator.Inject(ctx, propagation.HeaderCarrier(headers))
}
