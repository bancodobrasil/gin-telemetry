package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

type keyContext string

const (
	keyTracer     keyContext = "gin-telemetry-trace"
	keyRootSpan   keyContext = "gin-telemetry-root-span"
	keyPropagator keyContext = "gin-telemetry-propagator"
	name          string     = "https://github.com/bancodobrasil/gin-telemetry"
	version       string     = "0.0.1-rc2"
)

var (
	service string
	// MiddlewareDisabled ...
	MiddlewareDisabled bool
)

// Middleware ...
func Middleware(serviceName string, opts ...Option) gin.HandlerFunc {
	MiddlewareDisabled = viper.GetViper().GetBool("TELEMETRY_DISABLED")

	if MiddlewareDisabled {
		log.Debug("*** Gin-telemetry disabled ***")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	log.Debug("Configuring gin-telemetry middleware ...")

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
		trace.WithInstrumentationVersion(version),
	)

	if cfg.Propagator == nil {
		cfg.Propagator = getDefaultTextMapPropagator()
	}

	// Instrument logrus
	log.AddHook(otellogrus.NewHook(
		otellogrus.WithLevels(
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		),
	))

	log.Debug("Gin-telemetry successfully configured!")

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

		spanCtx := context.WithValue(ctx, keyRootSpan, span)
		c.Request = c.Request.WithContext(spanCtx)
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
func GetTracer(ctx context.Context) *trace.Tracer {
	tracerInterface := ctx.Value(keyTracer)
	if tracerInterface == nil {
		return nil
	}
	return tracerInterface.(*trace.Tracer)
}

// GetRootSpan ...
func GetRootSpan(ctx context.Context) *trace.Span {
	span := ctx.Value(keyRootSpan)
	if span == nil {
		return nil
	}
	return span.(*trace.Span)
}

// Inject ...
func Inject(ctx context.Context, headers http.Header) {
	propagatorInteface := ctx.Value(keyPropagator)
	propagator := propagatorInteface.(propagation.TextMapPropagator)
	propagator.Inject(ctx, propagation.HeaderCarrier(headers))
}
