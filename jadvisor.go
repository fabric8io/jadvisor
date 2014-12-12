package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/jimmidyson/jadvisor/sinks"
	"github.com/jimmidyson/jadvisor/sources"
)

var argPollDuration = flag.Duration("poll_duration", 10*time.Second, "Polling duration")

func main() {
	flag.Parse()
	glog.Infof(strings.Join(os.Args, " "))
	glog.Infof("jAdvisor version %v", jadvisorVersion)
	err := doWork()
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func doWork() error {
	source, err := sources.NewSource()
	if err != nil {
		return err
	}
	sink, err := sinks.NewSink()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(*argPollDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			data, err := source.GetInfo()
			if err != nil {
				return err
			}

			glog.V(2).Infof("Got stats: %v", data)

			if err := sink.StoreData(data); err != nil {
				return err
			}
		}
	}
	return nil
}
