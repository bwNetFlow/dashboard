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
	connector.pushCounter(counters.Msgcount)
	connector.pushCounter(counters.KafkaOffsets)
	connector.pushCounter(counters.FlowNumber)
	connector.pushCounter(counters.FlowBytes)
	connector.pushCounter(counters.FlowPackets)
	connector.pushCounter(counters.HostBytes)
	connector.pushCounter(counters.HostConnections)
}

func (connector *Connector) pushCounter(counter counters.Counter) {
	for cid := range counter.CustomerIndex {
		counter.Access.Lock()
		connector.pushCounterCustomer(counter, cid)
		counter.Access.Unlock()
	}
}

func (connector *Connector) pushCounterCustomer(counter counters.Counter, cid string) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "bwnetflow-" + cid,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}

	// get measurements as points
	hashes := counter.CustomerIndex[cid]
	pts := transformCounter(counter, hashes)
	bp.AddPoints(pts)

	// Write the batch
	err = connector.influxClient.Write(bp)
	if err != nil {
		fmt.Println("Error writing to influx: ", err.Error())
	}
}

func transformCounter(counter counters.Counter, hashes []uint32) []*client.Point {
	pts := make([]*client.Point, 0)
	for _, hash := range hashes {
		items := counter.Fields[hash]
		labels := items.Label
		val := items.Value

		tags := map[string]string{}
		for k, v := range labels.Fields {
			tags[k] = v
		}
		fields := map[string]interface{}{
			"count": int64(val),
		}
		pt, err := client.NewPoint(counter.Name, tags, fields, time.Now())
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		pts = append(pts, pt)
	}
	return pts
}
