package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client interface {
	Call(req Request) (interface{}, error)
}

type Response struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}

type Request struct {
	// URL path
	Path []string
	// The body sent with the POST request
	Body interface{}
}

type upstashClient struct {
	url        string
	httpClient HTTPClient
	token      string
}

func New(
	// The Upstash endpoint you want to use
	url string,

	// Requests to the Upstash API must provide an API token.
	token string,

	// Use a custom HTTPClient.
	// This is mainly used for testing
	// Omit to use the default `net/http` implementation
	httpClient HTTPClient) Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	return &upstashClient{
		url,
		httpClient,
		token,
	}
}

// JSON marshal the body if present
func marshalBody(body interface{}) (io.Reader, error) {
	var payload io.Reader = nil
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		payload = bytes.NewBuffer(b)
	}
	return payload, nil
}

// Perform a request and return its response
func (c *upstashClient) request(path []string, body interface{}) (*http.Response, error) {
	payload, err := marshalBody(body)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal request body: %w", err)
	}
	url := fmt.Sprintf("%s/%s", c.url, strings.Join(path, "/"))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to perform request: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var responseBody map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&responseBody)
		if err != nil {
			return nil, fmt.Errorf("Unable to decode response body of bad response: %s: %w", res.Status, err)
		}

		// Try to prettyprint the response body
		// If that is not possible we return the raw body
		pretty, err := json.MarshalIndent(responseBody, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("Response returned status code %d: %+v, path: %s", res.StatusCode, responseBody, path)
		}
		return nil, fmt.Errorf("Response returned status code %d: %+v, path: %s", res.StatusCode, string(pretty), path)
	}

	return res, nil

}

// Call the API and unmarshal its response directly
func (c *upstashClient) Call(req Request) (interface{}, error) {

	httpResponse, err := c.request(req.Path, req.Body)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	var response Response
	err = json.NewDecoder(httpResponse.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal response: %w", err)
	}
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response.Result, nil
}
