package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HttpClient struct {
	Client *http.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		Client: http.DefaultClient,
	}
}

func (c *HttpClient) DoRequest(method, uri string, headers map[string][]string, body interface{}) (*http.Response, error) {
	var requestBody []byte
	if body != nil {
		var err error
		requestBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, uri, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		if len(value) > 0 {
			req.Header.Add(key, value[0])
		}
	}

	return c.Client.Do(req)
}

func (c *HttpClient) Get(url string, headers map[string][]string, queryParams map[string]string) (*http.Response, error) {
	// Construct the URL with query parameters
	reqURL := url
	if queryParams != nil {
		reqURL += "?"
		for key, value := range queryParams {
			reqURL += key + "=" + value + "&"
		}
		reqURL = reqURL[:len(reqURL)-1] // Remove the last "&" character
	}

	return c.DoRequest("GET", reqURL, headers, nil)
}

func (c *HttpClient) Post(url string, headers map[string][]string, queryParams map[string]string, body interface{}) (*http.Response, error) {
	return c.DoRequest("POST", url, headers, body)
}

func (c *HttpClient) Put(url string, headers map[string][]string, queryParams map[string]string, body interface{}) (*http.Response, error) {
	return c.DoRequest("PUT", url, headers, body)
}

func (c *HttpClient) Patch(url string, headers map[string][]string, queryParams map[string]string, body interface{}) (*http.Response, error) {
	return c.DoRequest("PATCH", url, headers, body)
}

func (c *HttpClient) Delete(url string, headers map[string][]string, queryParams map[string]string) (*http.Response, error) {
	return c.DoRequest("DELETE", url, headers, nil)
}
