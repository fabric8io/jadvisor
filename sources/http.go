package sources

import (
	"encoding/json"
	"net/http"
)

func PostRequestAndGetValue(client *http.Client, req *http.Request, value interface{}) error {
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	dec := json.NewDecoder(response.Body)
	dec.UseNumber()
	err = dec.Decode(value)
	if err != nil {
		return err
	}
	return nil
}
