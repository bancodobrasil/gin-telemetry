package telemetry

import (
	"fmt"
	"net/http"
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
	viper.SetDefault("TELEMETRY_EXPORTER_URL", "http://localhost:14268")
	viper.SetDefault("TELEMETRY_HTTPCLIENT_TLS", true)
	viper.SetDefault("TELEMETRY_DISABLED", false)
	if err := viper.ReadInConfig(); err == nil {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}
}

// NewJaegerProvider ...
func NewJaegerProvider() TracerProvider {
	jaegerURL := viper.GetString("TELEMETRY_EXPORTER_URL")
	exporter, err := jaeger.New(
		jaeger.WithCollectorEndpoint(
			jaeger.WithEndpoint(fmt.Sprintf("%s/api/traces", jaegerURL)),
			jaeger.WithHTTPClient(http.DefaultClient),
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
			attribute.String("environment", environment),
			attribute.Int64("ID", id),
		)),
	)
	return tracerProvider
}
