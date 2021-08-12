// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for the Dynatrace exporter.
type Config struct {
	config.ExporterSettings       `mapstructure:",squash"`
	confighttp.HTTPClientSettings `mapstructure:",squash"`

	exporterhelper.QueueSettings               `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings               `mapstructure:"retry_on_failure"`
	exporterhelper.ResourceToTelemetrySettings `mapstructure:"resource_to_telemetry_conversion"`

	// Dynatrace API token with metrics ingest permission
	APIToken string `mapstructure:"api_token"`

	// Tags will be added to all exported metrics
	Tags []string `mapstructure:"tags"`

	// String to prefix all metric names
	Prefix string `mapstructure:"prefix"`
}

// Sanitize ensures an API token has been provided
func (c *Config) Sanitize() error {
	c.APIToken = strings.TrimSpace(c.APIToken)

	if c.APIToken == "" {
		return errors.New("missing api_token")
	}

	if c.Endpoint == "" {
		return errors.New("missing endpoint")
	}

	if !(strings.HasPrefix(c.Endpoint, "http://") || strings.HasPrefix(c.Endpoint, "https://")) {
		return errors.New("endpoint must start with https:// or http://")
	}

	if c.HTTPClientSettings.Headers == nil {
		c.HTTPClientSettings.Headers = make(map[string]string)
	}

	c.HTTPClientSettings.Headers["Content-Type"] = "text/plain; charset=UTF-8"
	c.HTTPClientSettings.Headers["Authorization"] = fmt.Sprintf("Api-Token %s", c.APIToken)
	c.HTTPClientSettings.Headers["User-Agent"] = "opentelemetry-collector"

	return nil
}
