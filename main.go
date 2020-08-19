package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

const version string = "1.2.3"

var (
	showVersion   = flag.Bool("version", false, "Print version information.")
	listenAddress = flag.String("web.listen-address", ":9324", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	newFormat     = flag.Bool("format.new", false, "New metric format (more convenient / generic)")
)

//immutable

func init() {
	flag.Usage = func() {
		fmt.Println("Usage: ip_exporter [ ... ]\n\nParameters:")
		fmt.Println()
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	startServer()
}

func printVersion() {
	fmt.Println("ip_exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Daniel Czerwonk")
}

func startServer() {
	log.Infof("Starting bird exporter (Version: %s)\n", version)

	if !*newFormat {
		log.Info("INFO: You are using the old metric format. Please consider using the new (more convenient one) by setting -format.new=true.")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Bird Routing Daemon Exporter (Version ` + version + `)</title></head>
			<body>
			<h1>Bird Routing Daemon Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			<h2>More information:</h2>
			<p><a href="https://github.com/czerwonk/ip_exporter">github.com/czerwonk/ip_exporter</a></p>
			</body>
			</html>`))
	})
	http.HandleFunc(*metricsPath, handleMetricsRequest)

	log.Infof("Listening for %s on %s\n", *metricsPath, *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}

func handleMetricsRequest(w http.ResponseWriter, r *http.Request) {
	reg := prometheus.NewRegistry()

	c := NewMetricCollector()
	reg.MustRegister(c)

	promhttp.HandlerFor(reg, promhttp.HandlerOpts{
		ErrorLog:      log.NewErrorLogger(),
		ErrorHandling: promhttp.ContinueOnError}).ServeHTTP(w, r)
}
