module github.com/bancodobrasil/gin-telemetry

go 1.16

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/go-playground/assert/v2 v2.0.1
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/viper v1.10.1
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.1.15
	go.opentelemetry.io/contrib/propagators/b3 v1.6.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/exporters/jaeger v1.8.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0
)
