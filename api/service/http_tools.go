package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func query(url string, dst interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, dst)
}
