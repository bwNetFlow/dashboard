package main

import (
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	msgcount := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "kafka_messages_total",
			Help: "Number of Kafka messages",
		})

	msgcount.Inc()
	msgcount.Inc()

	var x int
	x = 44
	println(x)

	s := reflect.ValueOf(x).Elem()
	println(s.String())
	//metric := s.FieldByName("valInt").Interface()
	//fmt.Println(metric)
}
