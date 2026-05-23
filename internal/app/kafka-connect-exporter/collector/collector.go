package collector

import (
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/vinted/kafka-connect-exporter/internal/app/kafka-connect-exporter/client"
)

type collector struct {
	client                   client.Client
	URI                      string
	up                       prometheus.Gauge
	connectorsCount          prometheus.Gauge
	isConnectorRunning       *prometheus.Desc
	areConnectorTasksRunning *prometheus.Desc
}

type connectors []string

var supportedSchema = map[string]bool{
	"http":  true,
	"https": true,
}

func NewCollector(uri, nameSpace, user, pass string, tlsConfig *client.TLSConfig) Collector {
	parseURI, err := url.Parse(uri)
	if err != nil {
		logrus.Errorf("%v", err)
		os.Exit(1)
	}
	if !supportedSchema[parseURI.Scheme] {
		logrus.Error("schema not supported")
		os.Exit(1)
	}

	logrus.Infoln("Collecting data from:", uri)

	// Optionally provide kafka connect basic auth credentials
	var authCredentials *client.AuthCredentials = nil
	if user != "" && pass != "" {
		authCredentials = &client.AuthCredentials{
			User:     user,
			Password: pass,
		}
	}

	return &collector{
		client: client.NewClient(uri, authCredentials, tlsConfig),
		URI:    uri,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "up",
			Help:      "was the last scrape of kafka connect successful?",
		}),
		connectorsCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Subsystem: "connectors",
			Name:      "count",
			Help:      "number of deployed connectors",
		}),
		isConnectorRunning: prometheus.NewDesc(
			prometheus.BuildFQName(nameSpace, "connector", "state_running"),
			"is the connector running?",
			[]string{"connector", "consumer_group", "state", "worker_id"},
			nil,
		),
		areConnectorTasksRunning: prometheus.NewDesc(
			prometheus.BuildFQName(nameSpace, "connector", "tasks_state"),
			"the state of tasks. 0-failed, 1-running, 2-unassigned, 3-paused",
			[]string{"connector", "consumer_group", "state", "worker_id", "id"},
			nil,
		),
	}
}
