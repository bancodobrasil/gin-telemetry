package tracing

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type configuration struct {
	Provider    trace.TracerProvider
	Propagators propagation.TextMapPropagator
}

// Option is ...
type Option interface {
	apply(*configuration)
}

type option func(*configuration)

func (ic option) apply(c *configuration) {
	ic.apply(c)
}

func withProvider(provider trace.TracerProvider) Option {
	return option(func(c *configuration) {
		if provider != nil {
			c.Provider = provider
		}
	})
}

func withPropagators(propagators propagation.TextMapPropagator) Option {
	return option(func(c *configuration) {
		if propagators != nil {
			c.Propagators = propagators
		}
	})
}
