package connectors

import (
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/connectors/influx"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/connectors/prometheus"
)

// NewInfluxConnector instantiates a new InfluxConnector
func NewInfluxConnector(url string, username string, password string, exportFreq int) influx.Connector {
	return influx.Connector{
		URL:        url,
		Username:   username,
		Password:   password,
		ExportFreq: exportFreq,
	}
}

// NewPrometheusConnector instantiates a new PrometheusConnector
func NewPrometheusConnector(addr string) prometheus.Connector {
	return prometheus.Connector{
		Addr: addr,
	}
}
