package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(host string, port int, key string, json bool) (string, error) {
	// Get the value for the key
	client := &http.Client{
		CheckRedirect: http.DefaultClient.CheckRedirect,
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/api/lookup/%s", host, port, key), nil)
	if err != nil {
		return "", err
	}
	if json {
		req.Header.Set("Accept", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error: %s", resp.Status)
	}

	io, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(io), nil
}
