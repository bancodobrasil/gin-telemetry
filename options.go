package telemetry

import (
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type configuration struct {
	Provider   trace.TracerProvider
	Propagator propagation.TextMapPropagator
}

// Option ...
type Option interface {
	apply(*configuration)
}

type optionFunc func(*configuration)

func (o optionFunc) apply(c *configuration) {
	o(c)
}

// WithProvider ...
func WithProvider(provider trace.TracerProvider) Option {
	return optionFunc(func(c *configuration) {
		if provider != nil {
			c.Provider = provider
		}
	})
}

// WithPropagators ...
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(c *configuration) {
		if propagators != nil {
			c.Propagator = propagators
		}
	})
}
