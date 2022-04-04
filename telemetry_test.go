package telemetry

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var serverTest *httptest.Server

func TestPropagation(t *testing.T) {
	provider := trace.NewNoopTracerProvider()
	propagator := b3.New()

	r := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	ctx := context.Background()

	sc := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: trace.TraceID{0x01},
		SpanID:  trace.SpanID{0x01},
	})

	ctx = trace.ContextWithRemoteSpanContext(ctx, sc)
	ctx, _ = provider.Tracer(name).Start(ctx, "test")
	propagator.Inject(ctx, propagation.HeaderCarrier(r.Header))

	router := gin.New()

	router.Use(Middleware("middleware-test", WithProvider(provider), WithPropagators(propagator)))
	router.GET("/user/:id", func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		assert.Equal(t, sc.TraceID(), span.SpanContext().TraceID())
		assert.Equal(t, sc.SpanID(), span.SpanContext().SpanID())
	})
	router.ServeHTTP(w, r)

}

func shutdown() {

}
