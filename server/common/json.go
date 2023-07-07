package common

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

var httpClient = &http.Client{Timeout: 20 * time.Second}

func GetJson(url string, retries int, target interface{}) (int, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		if retries > 0 {
			return GetJson(url, retries-1, target)
		}
		return 0, err
	}
	if resp.StatusCode != 200 {
		if retries > 0 {
			return GetJson(url, retries-1, target)
		}
		return resp.StatusCode, errors.New(resp.Status)
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(target)
}
