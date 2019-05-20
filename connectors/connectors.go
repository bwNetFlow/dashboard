package connectors

import (
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/connectors/influx"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/connectors/prometheus"
)

// NewInfluxConnector instantiates a new InfluxConnector
func NewInfluxConnector(url string, username string, password string, database string, exportFreq int, perCid bool) influx.Connector {
	return influx.Connector{
		URL:        url,
		Username:   username,
		Password:   password,
		Database:   database,
		ExportFreq: exportFreq,
		PerCid:     perCid,
	}
}

// NewPrometheusConnector instantiates a new PrometheusConnector
func NewPrometheusConnector(addr string) prometheus.Connector {
	return prometheus.Connector{
		Addr: addr,
	}
}
