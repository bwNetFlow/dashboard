package influx

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/counters"
)

// Connector provides export features to Influx
type Connector struct {
	URL          string
	Username     string
	Password     string
	ExportFreq   int
	influxClient client.Client
	running      bool
}

// Initialize established a connection to Influx DB
func (connector *Connector) Initialize() {
	var err error
	connector.influxClient, err = client.NewHTTPClient(client.HTTPConfig{
		Addr:     connector.URL,
		Username: connector.Username,
		Password: connector.Password,
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
		return
	}

	connector.running = true
	go connector.startPushCycle()

}

// Close closes the connection to influx
func (connector *Connector) Close() {
	connector.running = false
	connector.influxClient.Close()
}

func (connector *Connector) startPushCycle() {
	start := time.Now()
	for connector.running {
		nextRun := start.Add(time.Duration(connector.ExportFreq) * time.Second)
		time.Sleep(nextRun.Sub(time.Now()))

		connector.push()
		start = time.Now()
	}
}

func (connector *Connector) push() {
	fmt.Print("Push stuff to Influx")

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "bwnetflow",
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}

	// Transform counter values to influx points
	pts := transformCounter(counters.Msgcount)
	bp.AddPoints(pts)
	pts = transformCounter(counters.FlowNumber)
	bp.AddPoints(pts)
	pts = transformCounter(counters.FlowBytes)
	bp.AddPoints(pts)
	pts = transformCounter(counters.FlowPackets)
	bp.AddPoints(pts)

	// Write the batch
	err = connector.influxClient.Write(bp)
	if err != nil {
		fmt.Println("Error writing to influx: ", err.Error())
	}
}
