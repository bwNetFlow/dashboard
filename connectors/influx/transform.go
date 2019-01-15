package influx

import (
	"fmt"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"omi-gitlab.e-technik.uni-ulm.de/bwnetflow/kafka/consumer_dashboard/counters"
)

func transformCounter(counter counters.Counter) []*client.Point {
	pts := make([]*client.Point, 0)
	counter.Access.Lock()
	for _, items := range counter.Fields {
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
	counter.Access.Unlock()
	return pts
}
