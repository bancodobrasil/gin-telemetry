package tracing

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// TracerProvider is ...
type TracerProvider struct {
	name string
	*sdktrace.TracerProvider
}

// ITracerProvider is ...
type ITracerProvider struct {
	GetName string
}

// TracerProvider.GetName is ...
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
	viper.SetDefault("JAEGER_URL", "http://0.0.0.0:14268")
	if err := viper.ReadInConfig(); err == nil {
		log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}

// NewJaegerProvider is ...
func NewJaegerProvider() TracerProvider {
	jaegerURL := viper.GetString("JAEGER_URL")
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(fmt.Sprintf("%s/api/traces", jaegerURL))))
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
