package sources

import (
	"flag"

	kube_api "github.com/GoogleCloudPlatform/kubernetes/pkg/api")

var jubeEnv = flag.Bool("jube", false, "Are we running in Jube?")

func newEnvironment() Environment {
	isJube := *jubeEnv // TODO -- any better way then flag?

	if isJube {
		return Jube{}
	} else {
		return Kubernetes{}
	}
}

type Jube struct {
}

func (self Jube) GetHost(pod *kube_api.Pod, port kube_api.Port) string {
	return pod.Status.Host;
}

func (self Jube) GetPort(pod *kube_api.Pod, port kube_api.Port) int {
	return port.HostPort
}

type Kubernetes struct {
}

func (self Kubernetes) GetHost(pod *kube_api.Pod, port kube_api.Port) string {
	return pod.Status.PodIP;
}

func (self Kubernetes) GetPort(pod *kube_api.Pod, port kube_api.Port) int {
	return port.ContainerPort
}
