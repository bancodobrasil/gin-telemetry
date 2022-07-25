module github.com/bancodobrasil/gin-telemetry

go 1.16

require (
	github.com/gin-gonic/gin v1.8.1
	github.com/go-playground/assert/v2 v2.0.1
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/viper v1.12.0
	go.opentelemetry.io/contrib/propagators/b3 v1.6.0
	go.opentelemetry.io/otel v1.8.0
	go.opentelemetry.io/otel/exporters/jaeger v1.8.0
	go.opentelemetry.io/otel/sdk v1.8.0
	go.opentelemetry.io/otel/trace v1.8.0
)

require (
	github.com/go-playground/validator/v10 v10.11.0 // indirect
	github.com/goccy/go-json v0.9.10 // indirect
	github.com/pelletier/go-toml/v2 v2.0.2 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/subosito/gotenv v1.4.0 // indirect
	github.com/uptrace/opentelemetry-go-extra/otellogrus v0.1.14
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.0.0-20220708220712-1185a9018129 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
)
