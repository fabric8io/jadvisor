package sources

import (
	"net/http"
	"fmt"
	"encoding/json"
	"github.com/golang/glog"
	"bytes")

type DmrContainer struct {
	Name    string      `json:"name,omitempty"`
	Host    string      `json:"host"`
	DmrPort int         `json:"dmrPort"`
	Stats   *StatsEntry `json:"stats,omitempty"`
}

type DmrRequest struct {
	operation 	string	`json:"operation"`
	name      	string	`json:"name"`
	pretty		int		`json:"json.pretty"`
}

func (self *DmrContainer) GetName() string {
	return self.Name
}

func (self *DmrContainer) GetStats() (*StatsEntry, error) {
	dmrRequest := DmrRequest{
		operation: "read-attribute",
		name: "server-state",
		pretty: 1,
	}
	jsonMap, err := self.getStats(dmrRequest)
	if err != nil {
		return nil, err
	}
	glog.Info(jsonMap)

	return &StatsEntry{}, nil
}

func (self *DmrContainer) getStats(dmrRequest DmrRequest) (map[string]string, error) {
	reqBody, err := json.Marshal(dmrRequest)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("http://%s:%d/management", self.Host, self.DmrPort)
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	var jsonMap map[string]string
	err = PostRequestAndGetValue(&http.Client{}, req, jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}
