package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

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
	probeURL      = "https://www.google.com/"
	interval      = 60
	timeout       = 60
)

var (
	logger           log.Logger
	requestsDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: namespace,
		Name:      "probe_duration_seconds",
		Help:      "Duration of HTTP probes.",
		Buckets:   []float64{1, 2, 4, 8, 16, 32},
	})
	// One for total is not needed because a histogram already has a counter
	requestsFailed = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "probes_failed_total",
		Help:      "Amount of probes performed.",
	})
)

func init() {
	// version.NewCollector(namespace) doesn't give us good information at the moment because our
	// build pipeline doesn't pass version, nor commit SHA as build arguments
	prometheus.MustRegister(version.NewCollector(namespace))
	prometheus.MustRegister(requestsFailed)
	prometheus.MustRegister(requestsDuration)
}

func main() {
	logger = promlog.New(&promlog.Config{})
	level.Info(logger).Log("msg", "Starting http-prober", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	go func() {
		t := time.NewTicker(interval * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				err := Probe()
				if err != nil {
					level.Info(logger).Log("msg", "Failed probe", "err", err)
				}
			}
		}
	}()

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

func Probe() error {
	req, err := http.NewRequest("GET", probeURL, nil)
	if err != nil {
		return fmt.Errorf("Failed building request for url: %s", probeURL)
	}
	client := http.Client{
		Timeout: timeout * time.Second,
	}

	start := time.Now()
	res, err := client.Do(req)
	duration := float64(time.Since(start)) / float64(time.Second)
	if err != nil {
		requestsFailed.Inc()
		requestsDuration.Observe(duration)
		return fmt.Errorf("Failed request for url: %s", probeURL)
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		requestsFailed.Inc()
	}

	level.Info(logger).Log("msg", "Probed", "url", probeURL, "code", res.StatusCode)
	requestsDuration.Observe(duration)

	return nil
}
