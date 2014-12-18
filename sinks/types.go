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
	colTimestamp     string = "time"
	colPodName       string = "pod"
	colPodStatus     string = "pod_status"
	colPodIP         string = "pod_ip"
	colLabels        string = "labels"
	colHostName      string = "hostname"
	colContainerName string = "container_name"
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
