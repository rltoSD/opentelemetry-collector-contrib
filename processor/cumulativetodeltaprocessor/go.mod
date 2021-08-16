module github.com/open-telemetry/opentelemetry-collector-contrib/processor/cumulativetodeltaprocessor

go 1.16

require (
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.31.1-0.20210810171211-8038673eba9e
	go.opentelemetry.io/collector/model v0.31.1-0.20210810171211-8038673eba9e
	go.uber.org/zap v1.19.0
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/metrics => ./../../internal/aws/metrics