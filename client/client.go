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
	Read(req Request) (interface{}, error)
	Write(req Request) (interface{}, error)
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
	edgeUrl    string
	httpClient HTTPClient
	token      string
}

func New(
	// The Upstash endpoint you want to use
	url string,
	edgeUrl string,

	// Requests to the Upstash API must provide an API token.
	token string,

) Client {
	httpClient := &http.Client{}

	return &upstashClient{
		url,
		edgeUrl,
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
func (c *upstashClient) request(method string, path []string, body interface{}) (interface{}, error) {
	payload, err := marshalBody(body)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal request body: %w", err)
	}

	baseUrl := c.url
	if method == "GET" && c.edgeUrl != "" {
		baseUrl = c.edgeUrl
	}

	url := fmt.Sprintf("%s/%s", baseUrl, strings.Join(path, "/"))
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, fmt.Errorf("Unable to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Unable to perform request: %w", err)
	}
	defer res.Body.Close()

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

	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal response: %w", err)
	}
	if response.Error != "" {
		return nil, fmt.Errorf(response.Error)
	}
	return response.Result, nil

}

func (c *upstashClient) Read(req Request) (interface{}, error) {
	return c.request("GET", req.Path, nil)
}

// Call the API and unmarshal its response directly
func (c *upstashClient) Write(req Request) (interface{}, error) {
	return c.request("POST", req.Path, req.Body)
}
