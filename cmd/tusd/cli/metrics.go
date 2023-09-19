package cli

import (
	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/hooks"
	"github.com/tus/tusd/v2/pkg/prometheuscollector"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricsOpenConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "tusd_connections_open",
	Help: "Current number of open connections.",},
	[]string{"computername"},
	
)

func SetupMetrics(mux *http.ServeMux, handler *handler.Handler) {
	prometheus.MustRegister(MetricsOpenConnections)
	prometheus.MustRegister(hooks.MetricsHookErrorsTotal)
	prometheus.MustRegister(hooks.MetricsHookInvocationsTotal)
	prometheus.MustRegister(prometheuscollector.New(handler.Metrics))

	stdout.Printf("Using %s as the metrics path.\n", Flags.MetricsPath)
	mux.Handle(Flags.MetricsPath, promhttp.Handler())
}
