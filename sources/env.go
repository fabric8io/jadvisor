package sources

import (
	"flag"

	kube_api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
)

var jubeEnv = flag.Bool("jube", false, "Are we running Jube?")

func newEnvironment() *Environment {
	if isJube() {
		jube := new(Jube)
		env := Environment(jube)
		return &env
	} else {
		kube := new(Kubernetes)
		env := Environment(kube)
		return &env
	}
}

func isJube() bool {
	return *jubeEnv // TODO -- any better way then flag?
}

type Jube struct {
}

func (self *Jube) GetHost(pod *kube_api.Pod, port kube_api.Port) string {
	return pod.Status.Host
}

func (self *Jube) GetPort(pod *kube_api.Pod, port kube_api.Port) int {
	return port.HostPort
}

type Kubernetes struct {
}

func (self *Kubernetes) GetHost(pod *kube_api.Pod, port kube_api.Port) string {
	return pod.Status.PodIP
}

func (self *Kubernetes) GetPort(pod *kube_api.Pod, port kube_api.Port) int {
	return port.ContainerPort
}
