package telemetry

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// TracerProvider ...
type TracerProvider struct {
	name string
	*sdktrace.TracerProvider
}

// ITracerProvider ...
type ITracerProvider struct {
	GetName string
}

// GetName ...
func (t *TracerProvider) GetName() string {
	return t.name
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exePath := filepath.Dir(ex)

	viper.AddConfigPath(exePath)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
	viper.SetDefault("TELEMETRY_EXPORTER_JAEGER_AGENT_HOST", "localhost")
	viper.SetDefault("TELEMETRY_EXPORTER_JAEGER_AGENT_PORT", "6831")
	viper.SetDefault("TELEMETRY_HTTPCLIENT_TLS", true)
	viper.SetDefault("TELEMETRY_DISABLED", false)
	viper.SetDefault("TELEMETRY_ENVIRONMENT", "test")

	err = viper.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug("Config file not found.")
		} else {
			log.Errorf("Config file corrupted. Cause: %v", err)
			return
		}
	} else {
		log.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}
}

// NewJaegerProvider ...
func NewJaegerProvider() TracerProvider {
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(viper.GetString("TELEMETRY_EXPORTER_JAEGER_AGENT_HOST")),
			jaeger.WithAgentPort(viper.GetString("TELEMETRY_EXPORTER_JAEGER_AGENT_PORT")),
		),
	)
	if err != nil {
		log.Fatal(err)
	}
	return TracerProvider{
		name:           "jaeger",
		TracerProvider: getTracerProvider(exporter),
	}
}

func getTracerProvider(exporter sdktrace.SpanExporter) *sdktrace.TracerProvider {
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(service),
			attribute.String("environment", viper.GetString("TELEMETRY_ENVIRONMENT")),
			attribute.Int64("ID", int64(os.Getpid())),
		)),
	)
	return tracerProvider
}
