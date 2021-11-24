// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusreceiver

import (
	"testing"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/model/pdata"
)

const targetExternalLabels = `
# HELP go_threads Number of OS threads created
# TYPE go_threads gauge
go_threads 19`

func TestExternalLabels(t *testing.T) {
	targets := []*testData{
		{
			name: "target1",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetExternalLabels},
			},
			validateFunc: verifyExternalLabels,
		},
	}

	mp, cfg, err := setupMockPrometheus(targets...)
	cfg.GlobalConfig.ExternalLabels = labels.FromStrings("key", "value")
	require.Nilf(t, err, "Failed to create Prometheus config: %v", err)

	testComponentCustomConfig(t, targets, mp, cfg)
}

func verifyExternalLabels(t *testing.T, td *testData, rms []*pdata.ResourceMetrics) {
	verifyNumScrapeResults(t, td, rms)
	require.Greater(t, len(rms), 0, "At least one resource metric should be present")

	wantAttributes := td.attributes
	metrics1 := rms[0].InstrumentationLibraryMetrics().At(0).Metrics()
	ts1 := metrics1.At(0).Gauge().DataPoints().At(0).Timestamp()
	doCompare(t, "scrape-externalLabels", wantAttributes, rms[0], []testExpectation{
		assertMetricPresent("go_threads",
			compareMetricType(pdata.MetricDataTypeGauge),
			[]dataPointExpectation{
				{
					numberPointComparator: []numberPointComparator{
						compareTimestamp(ts1),
						compareDoubleValue(19),
						compareAttributes(map[string]string{"key": "value"}),
					},
				},
			}),
	})
}

const targetLabelLimit1 = `
# HELP test_gauge0 This is my gauge
# TYPE test_gauge0 gauge
test_gauge0{label1="value1",label2="value2"} 10
`

func verifyLabelLimitTarget1(t *testing.T, td *testData, rms []*pdata.ResourceMetrics) {
	//each sample in the scraped metrics is within the configured label_limit, scrape should be successful
	verifyNumScrapeResults(t, td, rms)
	require.Greater(t, len(rms), 0, "At least one resource metric should be present")

	want := td.attributes
	metrics1 := rms[0].InstrumentationLibraryMetrics().At(0).Metrics()
	ts1 := metrics1.At(0).Gauge().DataPoints().At(0).Timestamp()

	doCompare(t, "scrape-labelLimit", want, rms[0], []testExpectation{
		assertMetricPresent("test_gauge0",
			compareMetricType(pdata.MetricDataTypeGauge),
			[]dataPointExpectation{
				{
					numberPointComparator: []numberPointComparator{
						compareTimestamp(ts1),
						compareDoubleValue(10),
						compareAttributes(map[string]string{"label1": "value1", "label2": "value2"}),
					},
				},
			},
		),
	})
}

const targetLabelLimit2 = `
# HELP test_gauge0 This is my gauge
# TYPE test_gauge0 gauge
test_gauge0{label1="value1",label2="value2",label3="value3"} 10
`

func verifyFailedScrape(t *testing.T, _ *testData, rms []*pdata.ResourceMetrics) {
	//Scrape should be unsuccessful since limit is exceeded in target2
	for _, rm := range rms {
		metrics := getMetrics(rm)
		assertUp(t, 0, metrics)
	}
}

func TestLabelLimitConfig(t *testing.T) {
	targets := []*testData{
		{
			name: "target1",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetLabelLimit1},
			},
			validateFunc: verifyLabelLimitTarget1,
		},
		{
			name: "target2",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetLabelLimit2},
			},
			validateFunc: verifyFailedScrape,
		},
	}

	mp, cfg, err := setupMockPrometheus(targets...)
	require.Nilf(t, err, "Failed to create Prometheus config: %v", err)

	// set label limit in scrape_config
	for _, scrapeCfg := range cfg.ScrapeConfigs {
		scrapeCfg.LabelLimit = 5
	}

	testComponentCustomConfig(t, targets, mp, cfg)
}

const targetLabelLimits1 = `
# HELP test_gauge0 This is my gauge
# TYPE test_gauge0 gauge
test_gauge0{label1="value1",label2="value2"} 10

# HELP test_counter0 This is my counter
# TYPE test_counter0 counter
test_counter0{label1="value1",label2="value2"} 1

# HELP test_histogram0 This is my histogram
# TYPE test_histogram0 histogram
test_histogram0_bucket{label1="value1",label2="value2",le="0.1"} 1000
test_histogram0_bucket{label1="value1",label2="value2",le="0.5"} 1500
test_histogram0_bucket{label1="value1",label2="value2",le="1"} 2000
test_histogram0_bucket{label1="value1",label2="value2",le="+Inf"} 2500
test_histogram0_sum{label1="value1",label2="value2"} 5000
test_histogram0_count{label1="value1",label2="value2"} 2500

# HELP test_summary0 This is my summary
# TYPE test_summary0 summary
test_summary0{label1="value1",label2="value2",quantile="0.1"} 1
test_summary0{label1="value1",label2="value2",quantile="0.5"} 5
test_summary0{label1="value1",label2="value2",quantile="0.99"} 8
test_summary0_sum{label1="value1",label2="value2"} 5000
test_summary0_count{label1="value1",label2="value2"} 1000
`

func verifyLabelConfigTarget1(t *testing.T, td *testData, rms []*pdata.ResourceMetrics) {
	verifyNumScrapeResults(t, td, rms)
	require.Greater(t, len(rms), 0, "At least one resource metric should be present")

	want := td.attributes
	metrics1 := rms[0].InstrumentationLibraryMetrics().At(0).Metrics()
	ts1 := metrics1.At(0).Gauge().DataPoints().At(0).Timestamp()

	e1 := []testExpectation{
		assertMetricPresent("test_counter0",
			compareMetricType(pdata.MetricDataTypeSum),
			[]dataPointExpectation{
				{
					numberPointComparator: []numberPointComparator{
						compareStartTimestamp(ts1),
						compareTimestamp(ts1),
						compareDoubleValue(1),
						compareAttributes(map[string]string{"label1": "value1", "label2": "value2"}),
					},
				},
			}),
		assertMetricPresent("test_gauge0",
			compareMetricType(pdata.MetricDataTypeGauge),
			[]dataPointExpectation{
				{
					numberPointComparator: []numberPointComparator{
						compareTimestamp(ts1),
						compareDoubleValue(10),
						compareAttributes(map[string]string{"label1": "value1", "label2": "value2"}),
					},
				},
			}),
		assertMetricPresent("test_histogram0",
			compareMetricType(pdata.MetricDataTypeHistogram),
			[]dataPointExpectation{
				{
					histogramPointComparator: []histogramPointComparator{
						compareHistogramStartTimestamp(ts1),
						compareHistogramTimestamp(ts1),
						compareHistogram(2500, 5000, []uint64{1000, 500, 500, 500}),
						compareHistogramAttributes(map[string]string{"label1": "value1", "label2": "value2"}),
					},
				},
			}),
		assertMetricPresent("test_summary0",
			compareMetricType(pdata.MetricDataTypeSummary),
			[]dataPointExpectation{
				{
					summaryPointComparator: []summaryPointComparator{
						compareSummaryStartTimestamp(ts1),
						compareSummaryTimestamp(ts1),
						compareSummary(1000, 5000, [][]float64{{0.1, 1}, {0.5, 5}, {0.99, 8}}),
						compareSummaryAttributes(map[string]string{"label1": "value1", "label2": "value2"}),
					},
				},
			}),
	}
	doCompare(t, "scrape-label-config-test", want, rms[0], e1)
}

const targetLabelNameLimit = `
# HELP test_gauge0 This is my gauge
# TYPE test_gauge0 gauge
test_gauge0{label1="value1",labelNameExceedingLimit="value2"} 10

# HELP test_counter0 This is my counter
# TYPE test_counter0 counter
test_counter0{label1="value1",label2="value2"} 1
`

func TestLabelNameLimitConfig(t *testing.T) {
	targets := []*testData{
		{
			name: "target1",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetLabelLimits1},
			},
			validateFunc: verifyLabelConfigTarget1,
		},
		{
			name: "target2",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetLabelNameLimit},
			},
			validateFunc: verifyFailedScrape,
		},
	}

	mp, cfg, err := setupMockPrometheus(targets...)
	require.Nilf(t, err, "Failed to create Prometheus config: %v", err)

	// set label name limit in scrape_config
	for _, scrapeCfg := range cfg.ScrapeConfigs {
		scrapeCfg.LabelNameLengthLimit = 20
	}

	testComponentCustomConfig(t, targets, mp, cfg)
}

const targetLabelValueLimit = `
# HELP test_gauge0 This is my gauge
# TYPE test_gauge0 gauge
test_gauge0{label1="value1",label2="label-value-exceeding-limit"} 10

# HELP test_counter0 This is my counter
# TYPE test_counter0 counter
test_counter0{label1="value1",label2="value2"} 1
`

func TestLabelValueLimitConfig(t *testing.T) {
	targets := []*testData{
		{
			name: "target1",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetLabelLimits1},
			},
			validateFunc: verifyLabelConfigTarget1,
		},
		{
			name: "target2",
			pages: []mockPrometheusResponse{
				{code: 200, data: targetLabelValueLimit},
			},
			validateFunc: verifyFailedScrape,
		},
	}

	mp, cfg, err := setupMockPrometheus(targets...)
	require.Nilf(t, err, "Failed to create Prometheus config: %v", err)

	//set label value length limit in scrape_config
	for _, scrapeCfg := range cfg.ScrapeConfigs {
		scrapeCfg.LabelValueLengthLimit = 25
	}

	testComponentCustomConfig(t, targets, mp, cfg)
}
