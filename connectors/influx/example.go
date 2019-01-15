package influx

import (
	"fmt"
	"time"

	"github.com/influxdata/influxdb1-client/v2"
)

func influxTest() {

	fmt.Println("Let's start to connect to influx")
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "admin",
		Password: "secret",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "bwnetflow",
		Precision: "s",
	})
	if err != nil {
		fmt.Println("Error creating NewBatchPoints: ", err.Error())
	}

	// Create a point and add to batch
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	pt, err := client.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)

	// Write the batch
	err = c.Write(bp)
	if err != nil {
		fmt.Println("Error writing to influx: ", err.Error())
	}
}
