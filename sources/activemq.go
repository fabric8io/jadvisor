package sources

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

const (
	AMQDefaultDomain string = "org.apache.activemq"
)

func GetAMQRequests(url string) ([]JolokiaRequest, error) {

	listRequest := JolokiaRequest{
		Type:     List,
		Path:     AMQDefaultDomain,
		MaxDepth: 1,
	}

	reqBody, err := json.Marshal(listRequest)
	if err != nil {
		return []JolokiaRequest{}, err
	}

	glog.V(2).Infof("Sending Jolokia request: %v", string(reqBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return []JolokiaRequest{}, err
	}
	var jolokiaResponse JolokiaResponse
	err = PostRequestAndGetValue(&http.Client{}, req, &jolokiaResponse)
	if err != nil {
		glog.Errorf("failed to list AMQ mbeans from jolokia url: %s - %s\n", url, err)
		return []JolokiaRequest{}, err
	}

	var requests []JolokiaRequest
	for key, _ := range jolokiaResponse.Value {
		requests = append(requests, JolokiaRequest{
			Type:  Read,
			MBean: fmt.Sprintf("%s:%s", AMQDefaultDomain, key),
		})
	}

	return requests, nil
}
