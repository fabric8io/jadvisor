package sources

import (
	"flag"
	"time"
)

var (
	argMaster = flag.String("kubernetes_master", "", "Kubernetes master IP")
)

// PodState is the state of a pod, used as either input (desired state) or output (current state)
type Pod struct {
	Namespace  string            `json:"namespace,omitempty"`
	Name       string            `json:"name,omitempty"`
	ID         string            `json:"id,omitempty"`
	Hostname   string            `json:"hostname,omitempty"`
	Containers []*Container      `json:"containers"`
	Status     string            `json:"status,omitempty"`
	PodIP      string            `json:"podIP,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type Container struct {
	Name        string        `json:"name,omitempty"`
	JolokiaPort int           `json:"jolokiaPort"`
	Stats       *JolokiaStats `json:"stats,omitempty"`
}

type JolokiaRequestType string

const (
	Search JolokiaRequestType = "search"
	Read   JolokiaRequestType = "read"
)

type JolokiaRequest struct {
	Type      JolokiaRequestType `json:"type"`
	MBean     string             `json:"mbean"`
	Attribute string             `json:"attribute,omitempty"`
	Path      string             `json:"path,omitempty"`
}

type JolokiaResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     map[string]interface{}
	Error     string
}

type JolokiaStats struct {
	// The time of this stat point.
	Timestamp time.Time   `json:"timestamp"`
	Memory    MemoryStats `json:"mem,omitempty"`
}

type MemoryStats struct {
	HeapUsage struct {
		Committed uint64 `json:"committed"`
		Init      uint64 `json:"init"`
		Max       int64  `json:"max"`
		Used      uint64 `json:"used"`
	} `json:"HeapMemoryUsage"`
	NonHeapUsage struct {
		Committed uint64 `json:"committed"`
		Init      uint64 `json:"init"`
		Max       int64  `json:"max"`
		Used      uint64 `json:"used"`
	} `json:"NonHeapMemoryUsage"`
}

func newContainer() *Container {
	return &Container{Stats: &JolokiaStats{}}
}

type ContainerData struct {
	Pods []Pod
}

type Source interface {
	GetInfo() (ContainerData, error)
}

func NewSource() (Source, error) {
	return newKubeSource()
}
