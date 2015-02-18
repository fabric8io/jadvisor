package sources

import (
	"flag"
	"time"

	kube_api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/types"
)

var (
	argMaster         = flag.String("kubernetes_master", "https://localhost:8443", "Kubernetes master address")
	argMasterVersion  = flag.String("kubernetes_version", "v1beta2", "Kubernetes api version")
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

type Container interface {
	GetName() string
	GetStats() (*StatsEntry, error)
}

type StatsEntry struct {
	// The time of this stat point.
	Timestamp time.Time             `json:"timestamp"`
	Stats     map[string]StatsValue `json:"stats,omitempty"`
}

type StatsValue map[string]interface{}

func newJolokiaContainer() *JolokiaContainer {
	return &JolokiaContainer{Stats: &StatsEntry{}}
}

func newDmrContainer() *DmrContainer {
	return &DmrContainer{Stats: &StatsEntry{}}
}

type ContainerData struct {
	Pods []Pod
}

type Source interface {
	GetData() (ContainerData, error)
}

func NewSource() (Source, error) {
	return newKubeSource()
}

type Environment interface {
	GetHost(pod *kube_api.Pod, port kube_api.Port) string
	GetPort(pod *kube_api.Pod, port kube_api.Port) int
}
