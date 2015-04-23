package sinks

import (
	"container/list"
	"flag"
	"time"

	"github.com/fabric8io/jadvisor/sources"
	"github.com/golang/glog"
)

var argMaxStorageDuration = flag.Duration("sink_memory_ttl", 1*time.Hour, "Time duration for which stats should be cached if the memory sink is used")

type MemorySink struct {
	containersData     *list.List
	oldestData         time.Time
	maxStorageDuration time.Duration
}

type entry struct {
	timestamp time.Time
	data      interface{}
}

func (self *MemorySink) reapOldData() {
	if self.containersData.Len() == 0 || time.Since(self.oldestData) < self.maxStorageDuration {
		return
	}
	// TODO(vishh): Reap old data.
}

func (self *MemorySink) handlePods(pods []sources.Pod) {
	for _, pod := range pods {
		for _, container := range pod.Containers {
			stats, err := container.GetStats()

			if err != nil {
				glog.Errorf("Error getting container [%s] stats: %s", container.GetName(), err)
			} else {
				for mbean, stats := range stats.Stats {
					glog.Infof("%s -> %s", mbean, stats)
				}
			}
		}
	}
}

func (self *MemorySink) StoreData(input Data) error {
	glog.Info("Storing data ...")
	if data, ok := input.(sources.ContainerData); ok {

		self.handlePods(data.Pods)

		for _, value := range data.Pods {
			self.containersData.PushFront(entry{time.Now(), value})
			if self.containersData.Len() == 1 {
				self.oldestData = time.Now()
			}
		}
	}
	self.reapOldData()
	return nil
}

func NewMemorySink() Sink {
	return &MemorySink{
		containersData:     list.New(),
		oldestData:         time.Now(),
		maxStorageDuration: *argMaxStorageDuration,
	}
}
