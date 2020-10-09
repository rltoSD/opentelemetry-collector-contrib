package cortexexporter

import (
	"go.opentelemetry.io/collector/config/configmodels"
	prexp "go.opentelemetry.io/collector/exporter/prometheusexporter"
	prw "go.opentelemetry.io/collector/exporter/prometheusremotewriteexporter"
)

// Config defines configuration for Remote Write exporter.
type Config struct {
	// // squash ensures fields are correctly decoded in embedded struct.
	configmodels.ExporterSettings `mapstructure:",squash"`
	// exporterhelper.TimeoutSettings `mapstructure:",squash"`
	// exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	// exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// // prefix attached to each exported metric name
	// // See: https://prometheus.io/docs/practices/naming/#metric-names
	// Namespace string `mapstructure:"namespace"`

	PrwConfig PrwConfig `mapstructure:"prometheusremotewrite"`

	PrexpConfig PrexpConfig `mapstructure:"prometheus"`

	// AWS Sig V4 configuration options
	AuthSettings AuthSettings `mapstructure:"aws_auth"`

	// HTTPClientSettings confighttp.HTTPClientSettings `mapstructure:",squash"`
}

type PrwConfig struct {
	prw.Config `mapstructure:",squash"`
}

type PrexpConfig struct {
	prexp.Config `mapstructure:",squash"`
}

// AuthSettings defines AWS authentication configurations for SigningRoundTripper
type AuthSettings struct {
	Enabled bool `mapstructure:"enabled"`
	// region string for AWS Sig V4
	Region string `mapstructure:"region"`
	// service string for AWS Sig V4
	Service string `mapstructure:"service"`
	// whether AWS Sig v4 debug information should be printed
	Debug bool `mapstructure:"debug"`
}
