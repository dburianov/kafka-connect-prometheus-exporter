package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/vinted/kafka-connect-exporter/internal/app/kafka-connect-exporter/client"
	"github.com/vinted/kafka-connect-exporter/internal/app/kafka-connect-exporter/collector"
)

const (
	nameSpace  = "kafka_connect"
	version    = "dev"
	versionUrl = "https://github.com/vinted/kafka-connect-exporter"
)

var (
	showVersion           = flag.Bool("version", false, "show version and exit")
	listenAddress         = flag.String("listen-address", ":8080", "Address on which to expose metrics.")
	metricsPath           = flag.String("telemetry-path", "/metrics", "Path under which to expose metrics.")
	scrapeURI             = flag.String("scrape-uri", "http://127.0.0.1:8080", "URI on which to scrape kafka connect.")
	user                  = flag.String("user", "", "Optional username for authenticating to kafka-connect")
	pass                  = flag.String("pass", "", "Optional password for authenticating to kafka-connect")
	tlsEnabled            = flag.Bool("tls.enabled", false, "Connect to Kafka using TLS")
	tlsServerName         = flag.String("tls.server-name", "", "Used to verify the hostname on the returned certificates unless tls.insecure-skip-tls-verify is given. The kafka server's name should be given")
	tlsCAFile             = flag.String("tls.ca-file", "", "The optional certificate authority file for Kafka TLS client authentication")
	tlsCertFile           = flag.String("tls.cert-file", "", "The optional certificate file for Kafka client authentication")
	tlsKeyFile            = flag.String("tls.key-file", "", "The optional key file for Kafka client authentication")
	tlsInsecureSkipVerify = flag.Bool("tls.insecure-skip-tls-verify", false, "If true, the server's certificate will not be checked for validity")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("kafka_connect_exporter\n url: %s\n version: %s\n", versionUrl, version)
		os.Exit(2)
	}

	logrus.Infoln("Starting kafka_connect_exporter")

	prometheus.Unregister(prometheus.NewGoCollector())
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))

	// Create TLS configuration
	var tlsConfig *client.TLSConfig = nil
	if *tlsEnabled {
		tlsConfig = &client.TLSConfig{
			Enabled:               *tlsEnabled,
			ServerName:            *tlsServerName,
			CAFile:                *tlsCAFile,
			CertFile:              *tlsCertFile,
			KeyFile:               *tlsKeyFile,
			InsecureSkipTLSVerify: *tlsInsecureSkipVerify,
		}
	}

	prometheus.MustRegister(collector.NewCollector(*scrapeURI, nameSpace, *user, *pass, tlsConfig))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	logrus.Fatal(http.ListenAndServe(*listenAddress, nil))
}
