package sources

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
)

func PostRequestAndGetValue(client *http.Client, req *http.Request, value interface{}) error {
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	glog.V(5).Infof("Received response: %s", response.Body)
	dec := json.NewDecoder(response.Body)
	dec.UseNumber()
	err = dec.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
