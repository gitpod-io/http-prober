package main

import (
	"net/http"
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
)

const (
	namespace     = "http_prober"
	metricsPath   = "/metrics"
	listenAddress = ":8080"
)

var (
	buckets = []float64{1, 2, 4, 8, 16, 32, 64}
	logger  log.Logger
)

func init() {
	prometheus.MustRegister(version.NewCollector(namespace))
}

func main() {
	logger = promlog.New(&promlog.Config{})
	level.Info(logger).Log("msg", "Starting http-prober", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>HTTP-prober Exporter</title></head>
			<body>
			<h1>HTTP-prober Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	level.Info(logger).Log("msg", "Listening on", "address", listenAddress)
	server := &http.Server{Addr: listenAddress}

	if err := web.ListenAndServe(server, "", logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
