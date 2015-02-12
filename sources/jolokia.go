package sources

import (
    "github.com/golang/glog"
    "encoding/json"
    "net/http"
    "bytes"
    "time"
    "fmt")



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
    Value     StatsValue
    Error     string
}

type JolokiaContainer struct {
    Name        string        `json:"name,omitempty"`
    Host        string        `json:"host"`
    JolokiaPort int           `json:"jolokiaPort"`
    Stats       *StatsEntry   `json:"stats,omitempty"`
}

func (self *JolokiaContainer) GetStats() (*StatsEntry, error) {
    url := fmt.Sprintf("http://") + "/jolokia/"
    glog.V(2).Infof("Requesting jolokia stats from %s", url)

    jolokiaRequests := []JolokiaRequest{
        JVMRequest,
    }

    if amqRequests, err := GetAMQRequests(url); err == nil {
        jolokiaRequests = append(jolokiaRequests, amqRequests...)
    }

    reqBody, err := json.Marshal(jolokiaRequests)
    if err != nil {
        return &StatsEntry{}, err
    }

    glog.V(2).Infof("Sending Jolokia request: %v", string(reqBody))

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
    if err != nil {
        return &StatsEntry{}, err
    }
    var jolokiaResponses []JolokiaResponse
    err = PostRequestAndGetValue(&http.Client{}, req, &jolokiaResponses)
    if err != nil {
        glog.Errorf("failed to get stats from jolokia url: %s - %s\n", url, err)
        return &StatsEntry{}, nil
    }

    glog.V(2).Infof("Received jolokia response: %v", jolokiaResponses)

    jolokiaResponseStats := make(map[string]StatsValue)

    // Retrieve mbean stats
    for _, jolokiaResponse := range jolokiaResponses {
        jolokiaResponseStats[jolokiaResponse.Request["mbean"].(string)] = jolokiaResponse.Value
    }

    jolokiaStats := &StatsEntry{
        Timestamp: time.Unix(0, int64(jolokiaResponses[0].Timestamp)*int64(time.Second)),
        Stats:     jolokiaResponseStats,
    }

    glog.V(2).Infof("Retrieved jolokia stats: %v", jolokiaStats)

    return jolokiaStats, nil
}
