package sinks

import (
	"flag"
	"fmt"
)

var argSink = flag.String("sink", "memory", "Backend storage. Options are [memory | influxdb]")

type Data interface{}

type Sink interface {
	StoreData(data Data) error
}

const (
	statsTable               string = "stats"
	colTimestamp             string = "time"
	colPodName               string = "pod"
	colPodStatus             string = "pod_status"
	colPodIP                 string = "pod_ip"
	colLabels                string = "labels"
	colHostName              string = "hostname"
	colContainerName         string = "container_name"
	colHeapUsageCommitted    string = "heap_committed"
	colHeapUsageInit         string = "heap_init"
	colHeapUsageMax          string = "heap_max"
	colHeapUsageUsed         string = "heap_used"
	colNonHeapUsageCommitted string = "non_heap_committed"
	colNonHeapUsageInit      string = "non_heap_init"
	colNonHeapUsageMax       string = "non_heap_max"
	colNonHeapUsageUsed      string = "non_heap_used"
)

func NewSink() (Sink, error) {
	switch *argSink {
	case "memory":
		return NewMemorySink(), nil
	case "influxdb":
		return NewInfluxdbSink()
	default:
		return nil, fmt.Errorf("Invalid sink specified - %s", *argSink)
	}
}
