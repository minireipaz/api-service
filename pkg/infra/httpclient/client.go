package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"minireipaz/pkg/config"
	"net/http"
	"net/url"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	DoRequest(method, url, authToken string, body interface{}) ([]byte, error)
}

type ClientImpl struct {
	client *http.Client
}

func NewClientImpl(timeout time.Duration) *ClientImpl {
	log.Printf("WARN | Client Http not used Timeout context %v", timeout)
	return &ClientImpl{
		client: &http.Client{
			Timeout:   http.DefaultClient.Timeout,
			Transport: http.DefaultTransport,
		},
	}
}

func (c *ClientImpl) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *ClientImpl) DoRequest(method, url, authToken string, body interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error marshalling data: %v", err)
	}
	// not used with context
	req, err := http.NewRequest(method, url, NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	c.setHeaders(req, authToken)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error response: %v", err)
	}
	defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	bodyBytes, _ := io.ReadAll(resp.Body)
	// 	log.Printf("ERROR | failed to response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	// 	return nil, fmt.Errorf("error response: %v", err)
	// }

	return io.ReadAll(resp.Body)
}

func (c *ClientImpl) setHeaders(req *http.Request, token string) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
}

func NewBuffer(data []byte) io.Reader {
	return bytes.NewBuffer(data)
}

func getActionsURL(endpoint string) (string, error) {
	baseURI := fmt.Sprintf("%s%s", config.GetEnv("URI_ACTIONS", "http://localhost:4040"), endpoint)
	return validateURL(baseURI)
}

func validateURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL")
	}

	return parsedURL.String(), nil
}
