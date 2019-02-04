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
	Database     string
	ExportFreq   int
	PerCid       bool
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
	for connector.running {
		time.Sleep(time.Duration(connector.ExportFreq))
		connector.push()
	}
}

func (connector *Connector) push() {
	generalCounters := []counters.Counter{
		counters.Msgcount,
		counters.KafkaOffsets,
	}
	var customerCounters []counters.Counter
	if connector.PerCid {
		customerCounters = []counters.Counter{
			counters.FlowNumber,
			counters.FlowBytes,
			counters.FlowPackets,
			counters.HostBytes,
			counters.HostConnections,
		}
	} else {
		generalCounters = append(generalCounters, []counters.Counter{
			counters.FlowNumber,
			counters.FlowBytes,
			counters.FlowPackets,
			counters.HostBytes,
			counters.HostConnections,
		}...)
	}

	for _, counter := range generalCounters {
		connector.pushCounter(counter)
	}
	for _, counter := range customerCounters {
		for cid := range counter.CustomerIndex {
			connector.pushCounterCustomer(counter, cid)
		}
	}
}

func (connector *Connector) pushCounter(counter counters.Counter) {
	// Create Database if not exists
	db := connector.Database
	connector.createDb(db)

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}

	// get measurements as points
	counter.Access.Lock()
	for hash := range counter.Fields {
		pts := transformCounter(counter, []uint32{hash})
		bp.AddPoints(pts)
	}
	counter.Access.Unlock()

	// Write the batch
	err = connector.influxClient.Write(bp)
	if err != nil {
		fmt.Println("Error writing to influx: ", err.Error())
	}
}

func (connector *Connector) pushCounterCustomer(counter counters.Counter, cid string) {
	// Create Database if not exists
	db := connector.Database + "-" + cid
	connector.createDb(db)

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}

	// get measurements as points
	counter.Access.Lock()
	hashes := counter.CustomerIndex[cid]
	pts := transformCounter(counter, hashes)
	bp.AddPoints(pts)
	counter.Access.Unlock()

	// Write the batch
	err = connector.influxClient.Write(bp)
	if err != nil {
		fmt.Println("Error writing to influx: ", err.Error())
	}
}

func (connector *Connector) createDb(dbName string) {
	query := client.NewQuery("CREATE DATABASE \""+dbName+"\"  WITH DURATION 3d", "", "")
	response, err := connector.influxClient.Query(query)
	if err != nil {
		fmt.Printf("Error creating database %s in influx: %v\n", dbName, err.Error())
	}
	if response.Error() != nil {
		fmt.Printf("Error creating database %s in influx: %v\n", dbName, response.Error())
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
