// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"
)

// Type is the component type name.
const Type configmodels.Type = "hostmetricsreceiver"

// MetricIntf is an interface to generically interact with generated metric.
type MetricIntf interface {
	Name() string
	New() pdata.Metric
	Init(metric pdata.Metric)
}

// Intentionally not exposing this so that it is opaque and can change freely.
type metricImpl struct {
	name     string
	initFunc func(pdata.Metric)
}

// Name returns the metric name.
func (m *metricImpl) Name() string {
	return m.name
}

// New creates a metric object preinitialized.
func (m *metricImpl) New() pdata.Metric {
	metric := pdata.NewMetric()
	m.Init(metric)
	return metric
}

// Init initializes the provided metric object.
func (m *metricImpl) Init(metric pdata.Metric) {
	m.initFunc(metric)
}

type metricStruct struct {
	SystemCPULoadAverage15m MetricIntf
	SystemCPULoadAverage1m  MetricIntf
	SystemCPULoadAverage5m  MetricIntf
	SystemCPUTime           MetricIntf
	SystemMemoryUsage       MetricIntf
}

// Names returns a list of all the metric name strings.
func (m *metricStruct) Names() []string {
	return []string{
		"system.cpu.load_average.15m",
		"system.cpu.load_average.1m",
		"system.cpu.load_average.5m",
		"system.cpu.time",
		"system.memory.usage",
	}
}

var metricsByName = map[string]MetricIntf{
	"system.cpu.load_average.15m": Metrics.SystemCPULoadAverage15m,
	"system.cpu.load_average.1m":  Metrics.SystemCPULoadAverage1m,
	"system.cpu.load_average.5m":  Metrics.SystemCPULoadAverage5m,
	"system.cpu.time":             Metrics.SystemCPUTime,
	"system.memory.usage":         Metrics.SystemMemoryUsage,
}

func (m *metricStruct) ByName(n string) MetricIntf {
	return metricsByName[n]
}

func (m *metricStruct) FactoriesByName() map[string]func() pdata.Metric {
	return map[string]func() pdata.Metric{
		Metrics.SystemCPULoadAverage15m.Name(): Metrics.SystemCPULoadAverage15m.New,
		Metrics.SystemCPULoadAverage1m.Name():  Metrics.SystemCPULoadAverage1m.New,
		Metrics.SystemCPULoadAverage5m.Name():  Metrics.SystemCPULoadAverage5m.New,
		Metrics.SystemCPUTime.Name():           Metrics.SystemCPUTime.New,
		Metrics.SystemMemoryUsage.Name():       Metrics.SystemMemoryUsage.New,
	}
}

// Metrics contains a set of methods for each metric that help with
// manipulating those metrics.
var Metrics = &metricStruct{
	&metricImpl{
		"system.cpu.load_average.15m",
		func(metric pdata.Metric) {
			metric.SetName("system.cpu.load_average.15m")
			metric.SetDescription("Average CPU Load over 15 minutes.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeDoubleGauge)
		},
	},
	&metricImpl{
		"system.cpu.load_average.1m",
		func(metric pdata.Metric) {
			metric.SetName("system.cpu.load_average.1m")
			metric.SetDescription("Average CPU Load over 1 minute.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeDoubleGauge)
		},
	},
	&metricImpl{
		"system.cpu.load_average.5m",
		func(metric pdata.Metric) {
			metric.SetName("system.cpu.load_average.5m")
			metric.SetDescription("Average CPU Load over 5 minutes.")
			metric.SetUnit("1")
			metric.SetDataType(pdata.MetricDataTypeDoubleGauge)
		},
	},
	&metricImpl{
		"system.cpu.time",
		func(metric pdata.Metric) {
			metric.SetName("system.cpu.time")
			metric.SetDescription("Total CPU seconds broken down by different states.")
			metric.SetUnit("s")
			metric.SetDataType(pdata.MetricDataTypeDoubleSum)
			metric.DoubleSum().SetIsMonotonic(true)
			metric.DoubleSum().SetAggregationTemporality(pdata.AggregationTemporalityCumulative)
		},
	},
	&metricImpl{
		"system.memory.usage",
		func(metric pdata.Metric) {
			metric.SetName("system.memory.usage")
			metric.SetDescription("Bytes of memory in use.")
			metric.SetUnit("By")
			metric.SetDataType(pdata.MetricDataTypeIntSum)
			metric.IntSum().SetIsMonotonic(false)
			metric.IntSum().SetAggregationTemporality(pdata.AggregationTemporalityCumulative)
		},
	},
}

// M contains a set of methods for each metric that help with
// manipulating those metrics. M is an alias for Metrics
var M = Metrics

// Labels contains the possible metric labels that can be used.
var Labels = struct {
	// Cpu (CPU number starting at 0.)
	Cpu string
	// CPUState (Breakdown of CPU usage by type.)
	CPUState string
	// MemState (Breakdown of memory usage by type.)
	MemState string
}{
	"cpu",
	"state",
	"state",
}

// L contains the possible metric labels that can be used. L is an alias for
// Labels.
var L = Labels

// LabelCPUState are the possible values that the label "cpu.state" can have.
var LabelCPUState = struct {
	Idle      string
	Interrupt string
	Nice      string
	Softirq   string
	Steal     string
	System    string
	User      string
	Wait      string
}{
	"idle",
	"interrupt",
	"nice",
	"softirq",
	"steal",
	"system",
	"user",
	"wait",
}

// LabelMemState are the possible values that the label "mem.state" can have.
var LabelMemState = struct {
	Buffered          string
	Cached            string
	Inactive          string
	Free              string
	SlabReclaimable   string
	SlabUnreclaimable string
	Used              string
}{
	"buffered",
	"cached",
	"inactive",
	"free",
	"slab_reclaimable",
	"slab_unreclaimable",
	"used",
}