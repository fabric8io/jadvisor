package sources

import (
	"flag"
	"time"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/types"
)

var (
	argMaster         = flag.String("kubernetes_master", "https://localhost:8443", "Kubernetes master address")
	argMasterInsecure = flag.Bool("kubernetes_insecure", false, "Trust Kubernetes master certificate (if using https)")
)

// PodState is the state of a pod, used as either input (desired state) or output (current state)
type Pod struct {
	Namespace  string            `json:"namespace,omitempty"`
	Name       string            `json:"name,omitempty"`
	ID         types.UID         `json:"id,omitempty"`
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
	List   JolokiaRequestType = "list"
	Exec   JolokiaRequestType = "exec"
	Write  JolokiaRequestType = "write"
)

type JolokiaRequest struct {
	Type      JolokiaRequestType `json:"type"`
	MBean     string             `json:"mbean,omitempty"`
	Attribute interface{}        `json:"attribute,omitempty"`
	Path      string             `json:"path,omitempty"`
	MaxDepth  uint               `json:"maxDepth,omitempty"`
}

type JolokiaResponse struct {
	Status    uint32
	Timestamp uint32
	Request   map[string]interface{}
	Value     JolokiaValue
	Error     string
}

type JolokiaStats struct {
	// The time of this stat point.
	Timestamp time.Time               `json:"timestamp"`
	Stats     map[string]JolokiaValue `json:"stats,omitempty"`
}

type JolokiaValue map[string]interface{}

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
