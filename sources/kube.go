package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	kube_api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	kube_client "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	kube_labels "github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/golang/glog"
)

type KubeSource struct {
	client    *kube_client.Client
	lastQuery time.Time
}

func (self *KubeSource) parsePod(pod *kube_api.Pod) *Pod {
	localPod := Pod{
		Namespace:  pod.Namespace,
		Name:       pod.Name,
		ID:         pod.UID,
		PodIP:      pod.Status.PodIP,
		Hostname:   pod.Status.Host,
		Status:     string(pod.Status.Phase),
		Labels:     make(map[string]string, 0),
		Containers: make([]*Container, 0),
	}
	for key, value := range pod.Labels {
		localPod.Labels[key] = value
	}
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			if port.Name == "jolokia" || port.ContainerPort == 8778 {
				localContainer := newContainer()
				localContainer.Name = container.Name
				localContainer.JolokiaPort = port.ContainerPort
				localPod.Containers = append(localPod.Containers, localContainer)
				break
			}
		}
	}
	glog.V(2).Infof("found pod: %+v", localPod)

	return &localPod
}

func (self *KubeSource) getPods() ([]Pod, error) {
	pods, err := self.client.Pods(kube_api.NamespaceAll).List(kube_labels.Everything())
	if err != nil {
		return nil, err
	}
	glog.V(1).Infof("got pods from api server %+v", pods)
	out := make([]Pod, 0)
	for _, pod := range pods.Items {
		if pod.Status.Phase == kube_api.PodRunning {
			pod := self.parsePod(&pod)
			out = append(out, *pod)
		}
	}

	return out, nil
}

func (self *KubeSource) getStatsFromJolokia(jolokiaUrl string) (*JolokiaStats, error) {
	url := jolokiaUrl + "/jolokia/"
	glog.V(2).Infof("Requesting jolokia stats from %s", url)

	jolokiaRequests := []JolokiaRequest{
		JVMRequest,
	}

	if amqRequests, err := GetAMQRequests(url); err == nil {
		jolokiaRequests = append(jolokiaRequests, amqRequests...)
	}

	reqBody, err := json.Marshal(jolokiaRequests)
	if err != nil {
		return &JolokiaStats{}, err
	}

	glog.V(2).Infof("Sending Jolokia request: %v", string(reqBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return &JolokiaStats{}, err
	}
	var jolokiaResponses []JolokiaResponse
	err = PostRequestAndGetValue(&http.Client{}, req, &jolokiaResponses)
	if err != nil {
		glog.Errorf("failed to get stats from jolokia url: %s - %s\n", url, err)
		return &JolokiaStats{}, nil
	}

	glog.V(2).Infof("Received jolokia response: %v", jolokiaResponses)

	jolokiaResponseStats := make(map[string]JolokiaValue)

	// Retrieve mbean stats
	for _, jolokiaResponse := range jolokiaResponses {
		jolokiaResponseStats[jolokiaResponse.Request["mbean"].(string)] = jolokiaResponse.Value
	}

	jolokiaStats := &JolokiaStats{
		Timestamp: time.Unix(0, int64(jolokiaResponses[0].Timestamp)*int64(time.Second)),
		Stats:     jolokiaResponseStats,
	}

	glog.V(2).Infof("Retrieved jolokia stats: %v", jolokiaStats)

	return jolokiaStats, nil
}

func (self *KubeSource) GetInfo() (ContainerData, error) {
	pods, err := self.getPods()
	if err != nil {
		return ContainerData{}, err
	}
	for _, pod := range pods {
		for _, container := range pod.Containers {
			stats, err := self.getStatsFromJolokia(fmt.Sprintf("http://%s:%d", pod.PodIP, 8778))
			if err != nil {
				return ContainerData{}, err
			}
			container.Stats = stats
		}
	}

	self.lastQuery = time.Now()

	return ContainerData{Pods: pods}, nil
}

func newKubeSource() (*KubeSource, error) {
	if !(strings.HasPrefix(*argMaster, "http://") || strings.HasPrefix(*argMaster, "https://")) {
		*argMaster = "http://" + *argMaster
	}
	if len(*argMaster) == 0 {
		return nil, fmt.Errorf("kubernetes_master flag not specified")
	}
	kubeClient := kube_client.NewOrDie(&kube_client.Config{
		Host:     os.ExpandEnv(*argMaster),
		Version:  "v1beta2",
		Insecure: *argMasterInsecure,
	})

	return &KubeSource{
		client:    kubeClient,
		lastQuery: time.Now(),
	}, nil
}
