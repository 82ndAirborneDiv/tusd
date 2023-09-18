// package prometheuscollector allows to expose metrics for Prometheus.
//
// Using the provided collector, you can easily expose metrics for tusd in the
// Prometheus exposition format (https://prometheus.io/docs/instrumenting/exposition_formats/):
//
//	handler, err := handler.NewHandler(…)
//	collector := prometheuscollector.New(handler.Metrics)
//	prometheus.MustRegister(collector)

//  attempt to add COMPUTERNAME env var to metrics

package prometheuscollector

import (
	//	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tus/tusd/v2/pkg/handler"
	"os"
	"strconv"
	"sync/atomic"
)

var (
	requestsTotalDesc = prometheus.NewDesc(
		"tusd_requests_total",
		"Total number of requests served by tusd per method.",
		[]string{"method", "computername"}, nil)
	errorsTotalDesc = prometheus.NewDesc(
		"tusd_errors_total",
		"Total number of errors per status.",
		[]string{"status", "code", "computername"}, nil)
	bytesReceivedDesc = prometheus.NewDesc(
		"tusd_bytes_received",
		"Number of bytes received for uploads.",
		[]string{"computername"}, nil)
	uploadsCreatedDesc = prometheus.NewDesc(
		"tusd_uploads_created",
		"Number of created uploads.",
		[]string{"computername"}, nil)
	uploadsFinishedDesc = prometheus.NewDesc(
		"tusd_uploads_finished",
		"Number of finished uploads.",
		[]string{"computername"}, nil)
	uploadsTerminatedDesc = prometheus.NewDesc(
		"tusd_uploads_terminated",
		"Number of terminated uploads.",
		[]string{"computername"}, nil)
)

type Collector struct {
	metrics handler.Metrics
}

// New creates a new collector which read froms the provided Metrics struct.
func New(metrics handler.Metrics) Collector {
	return Collector{
		metrics: metrics,
	}
}

func (Collector) Describe(descs chan<- *prometheus.Desc) {
	descs <- requestsTotalDesc
	descs <- errorsTotalDesc
	descs <- bytesReceivedDesc
	descs <- uploadsCreatedDesc
	descs <- uploadsFinishedDesc
	descs <- uploadsTerminatedDesc
}

func (c Collector) Collect(metrics chan<- prometheus.Metric) {
	computerName := os.Getenv("COMPUTERNAME")

	for method, valuePtr := range c.metrics.RequestsTotal {
		metrics <- prometheus.MustNewConstMetric(
			requestsTotalDesc,
			prometheus.CounterValue,
			float64(atomic.LoadUint64(valuePtr)),
			method,
			computerName,
		)
	}

	for httpError, valuePtr := range c.metrics.ErrorsTotal.Load() {
		metrics <- prometheus.MustNewConstMetric(
			errorsTotalDesc,
			prometheus.CounterValue,
			float64(atomic.LoadUint64(valuePtr)),
			strconv.Itoa(httpError.StatusCode),
			httpError.ErrorCode,
			computerName,
		)
	}

	metrics <- prometheus.MustNewConstMetric(
		bytesReceivedDesc,
		prometheus.CounterValue,
		float64(atomic.LoadUint64(c.metrics.BytesReceived)),
		computerName,
	)

	metrics <- prometheus.MustNewConstMetric(
		uploadsFinishedDesc,
		prometheus.CounterValue,
		float64(atomic.LoadUint64(c.metrics.UploadsFinished)),
		computerName,
	)

	metrics <- prometheus.MustNewConstMetric(
		uploadsCreatedDesc,
		prometheus.CounterValue,
		float64(atomic.LoadUint64(c.metrics.UploadsCreated)),
		computerName,
	)

	metrics <- prometheus.MustNewConstMetric(
		uploadsTerminatedDesc,
		prometheus.CounterValue,
		float64(atomic.LoadUint64(c.metrics.UploadsTerminated)),
		computerName,
	)
}
