package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.up.Set(0)

	response, err := c.client.Get("/connectors")
	if err != nil {
		logrus.Errorf("Can't scrape kafka connect: %v", err)
		return
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			logrus.Errorf("Can't close connection to kafka connect: %v", err)
		}
	}()

	output, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logrus.Errorf("Can't scrape kafka connect: %v", err)
		return
	}

	var connectorsList connectors
	if err := json.Unmarshal(output, &connectorsList); err != nil {
		logrus.Errorf("Can't scrape kafka connect: %v", err)
		return
	}

	c.up.Set(1)
	c.connectorsCount.Set(float64(len(connectorsList)))

	ch <- c.up
	ch <- c.connectorsCount

	for _, connector := range connectorsList {

		connectorStatusResponse, err := c.client.Get("/connectors/" + connector + "/status")
		if err != nil {
			logrus.Errorf("Can't get /status for: %v", err)
			continue
		}

		connectorStatusOutput, err := ioutil.ReadAll(connectorStatusResponse.Body)
		if err != nil {
			logrus.Errorf("Can't read Body for: %v", err)
			continue
		}

		var connectorStatus status
		if err := json.Unmarshal(connectorStatusOutput, &connectorStatus); err != nil {
			logrus.Errorf("Can't decode response for: %v", err)
			continue
		}

		var isRunning float64 = 0
		if strings.ToLower(connectorStatus.Connector.State) == "running" {
			isRunning = 1
		}

		ch <- prometheus.MustNewConstMetric(
			c.isConnectorRunning, prometheus.GaugeValue, isRunning,
			connectorStatus.Name, connectorStatus.ConsumerGroup(), strings.ToLower(connectorStatus.Connector.State), connectorStatus.Connector.WorkerId,
		)

		for _, connectorTask := range connectorStatus.Tasks {

			var state float64
			switch taskState := strings.ToLower(connectorTask.State); taskState {
			case "running":
				state = 1
			case "unassigned":
				state = 2
			case "paused":
				state = 3
			default:
				state = 0
			}

			ch <- prometheus.MustNewConstMetric(
				c.areConnectorTasksRunning, prometheus.GaugeValue, state,
				connectorStatus.Name, connectorStatus.ConsumerGroup(), strings.ToLower(connectorTask.State), connectorTask.WorkerId, fmt.Sprintf("%d", int(connectorTask.Id)),
			)
		}

		err = connectorStatusResponse.Body.Close()
		if err != nil {
			logrus.Errorf("Can't close connection to connector: %v", err)
		}
	}
}
