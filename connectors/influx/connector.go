package influx

import (
	"fmt"
	"time"

	"github.com/bwNetFlow/dashboard/counters"
	client "github.com/influxdata/influxdb1-client/v2"
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
	databases    []string

	prevVals map[string]map[uint32]uint64
	// prevTimes map[string]map[uint32]time.Time
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

	connector.lookupDatabases()

	// Create Database if not exists
	db := connector.Database
	connector.createDb(db)

	connector.running = true

	connector.prevVals = make(map[string]map[uint32]uint64)
	// connector.prevTimes = make(map[string]map[uint32]time.Time)

	go connector.startPushCycle()

}

// Close closes the connection to influx
func (connector *Connector) Close() {
	connector.running = false
	connector.influxClient.Close()
}

func (connector *Connector) lookupDatabases() {
	dbQuery := client.Query{
		Command: "SHOW DATABASES",
	}
	response, err := connector.influxClient.Query(dbQuery)
	if err != nil || response.Error() != nil {
		fmt.Println("Error querying InfluxDB databases: ", err.Error())
		return
	}
	connector.databases = []string{}
	if len(response.Results) > 0 {
		for _, result := range response.Results {
			if len(result.Series) > 0 || len(result.Series[0].Values) > 0 {
				// found databases, save them
				dbs := result.Series[0].Values
				for _, db := range dbs {
					dbstr := fmt.Sprintf("%s", db[0])
					connector.databases = append(connector.databases, dbstr)
				}
			}
		}
	}
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
		// copy cids
		cids := counter.GetCids()

		// walk through each cid
		for _, cid := range cids {
			connector.pushCounterCustomer(counter, cid)
		}
	}
}

func (connector *Connector) pushCounter(counter counters.Counter) {
	db := connector.Database

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}

	// get measurements as points
	fields := counter.GetFields()
	for hash := range fields {
		pts := connector.transformCounter(counter, []uint32{hash})
		bp.AddPoints(pts)
	}

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
	hashes := counter.GetFieldHashesByCid(cid)
	pts := connector.transformCounter(counter, hashes)
	bp.AddPoints(pts)

	// Write the batch
	err = connector.influxClient.Write(bp)
	if err != nil {
		fmt.Println("Error writing to influx: ", err.Error())
	}
}

func (connector *Connector) createDb(dbName string) {
	// check if db exists
	for _, db := range connector.databases {
		if db == dbName {
			// found. don't create.
			return
		}
	}

	// if not exist, create
	query := client.NewQuery("CREATE DATABASE \""+dbName+"\"  WITH DURATION 3d", "", "")
	response, err := connector.influxClient.Query(query)
	if err != nil {
		fmt.Printf("Error creating database %s in influx: %v\n", dbName, err.Error())
	}
	if response.Error() != nil {
		fmt.Printf("Error creating database %s in influx: %v\n", dbName, response.Error())
	}

}

func (connector *Connector) transformCounter(counter counters.Counter, hashes []uint32) []*client.Point {
	pts := make([]*client.Point, 0)
	for _, hash := range hashes {
		items := counter.GetField(hash)
		labels := items.Label
		val := items.Value
		now := time.Now()

		prevVal := uint64(0)
		if _, ok := connector.prevVals[counter.GetName()]; ok {
			if _, ok2 := connector.prevVals[counter.GetName()][hash]; ok2 {
				prevVal = connector.prevVals[counter.GetName()][hash]
			}
		} else {
			connector.prevVals[counter.GetName()] = make(map[uint32]uint64)
		}
		/*
			prevTime := now.Add(time.Duration(connector.ExportFreq) * time.Second * -1)
			if _, ok := connector.prevTimes[counter.Name]; ok {
				if _, ok2 := connector.prevTimes[counter.Name][hash]; ok2 {
					prevTime = connector.prevTimes[counter.Name][hash]
				} else {
					fmt.Printf("unknown hash\n")
				}
			} else {
				connector.prevTimes[counter.Name] = make(map[uint32]time.Time)
			}*/
		valDiff := float64(val - prevVal)
		// timeDiff := now.Sub(prevTime).Seconds()
		// rate := valDiff / timeDiff
		connector.prevVals[counter.GetName()][hash] = val
		//connector.prevTimes[counter.Name][hash] = now

		if valDiff == 0 {
			// skip this point! no changes
			continue
		}

		tags := map[string]string{}
		for k, v := range labels.Fields {
			tags[k] = v
		}
		fields := map[string]interface{}{
			"count": int64(valDiff),
		}
		pt, err := client.NewPoint(counter.GetName(), tags, fields, now)
		if err != nil {
			fmt.Println("Error: ", err.Error())
		}
		pts = append(pts, pt)
	}
	return pts
}
