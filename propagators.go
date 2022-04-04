package telemetry

import "go.opentelemetry.io/otel/propagation"

func getDefaultTextMapPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
}
