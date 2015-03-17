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

type DmrAttributeRequest struct {
	Operation 	string	`json:"operation"`
	Name      	string	`json:"name"`
	Pretty		int		`json:"json.pretty"`
}

type DmrResourceRequest struct {
	Operation 			string		`json:"operation"`
	IncludeRuntime      bool		`json:"include-runtime"`
	Address				[]string	`json:"address"`
	Pretty				int			`json:"json.pretty"`
}

type DmrResponse struct {
	Outcome 				string		`json:"outcome"`
	Result      			interface{}	`json:"result"`
	FailureDescription      string		`json:"failure-description"`
	RolledBacked      		bool		`json:"rolled-back"`
}

type WebResult struct {
	BytesReceived		StringInt			`json:"bytesReceived"`
	BytesSent			StringInt			`json:"bytesSent"`
	EnableLookups		bool				`json:"enable-lookups"`
	Enabled				bool				`json:"enabled"`
	ErrorCount			StringInt			`json:"errorCount"`
	Executor			string				`json:"executor"`
	MaxConnections		int					`json:"max-connections"`
	MaxPostSize			int64				`json:"max-post-size"`
	MaxSavePostSize		int64				`json:"max-save-post-size"`
	MaxTime				StringInt			`json:"maxTime"`
	Name				string				`json:"name"`
	ProcessingTime		StringInt			`json:"processingTime"`
	Protocol			string				`json:"protocol"`
	ProxyName			string				`json:"proxy-name"`
	ProxyPort			string				`json:"proxy-port"`
	RedirectPort		int					`json:"redirect-port"`
	RequestCount		StringInt			`json:"requestCount"`
	Scheme				string				`json:"scheme"`
	Secure				bool				`json:"secure"`
	SocketBinding		string				`json:"socket-binding"`
	SSL					string				`json:"ssl"`
	VirtualServer		string				`json:"virtual-server"`
}

func (self *DmrContainer) GetName() string {
	return self.Name
}

func (self *DmrContainer) GetStats() (*StatsEntry, error) {
	dmrRequest := DmrResourceRequest{
		Operation: "read-resource",
		IncludeRuntime: true,
		Address: []string{"subsystem", "web", "connector", "http"},
		Pretty: 1,
	}

	wr := WebResult{}
	dmrResponse := DmrResponse{
		Result: &wr,
	}

	err := self.getStats(&dmrRequest, &dmrResponse)
	if err != nil {
		return nil, err
	}

	glog.Infof("outcome: %s, result: %s, failure: %s", dmrResponse.Outcome, dmrResponse.Result, dmrResponse.FailureDescription)

	return &StatsEntry{}, nil
}

func (self *DmrContainer) getStats(request interface{}, result interface{}) error {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("http://%s:%d/management", self.Host, self.DmrPort)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	err = PostRequestAndGetValue(&http.Client{}, req, result)
	if err != nil {
		return err
	}

	return nil
}
