package httpclient

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HttpClientImpl struct{}

func (c *HttpClientImpl) Do(req *http.Request) (*http.Response, error) {
	client := &http.Client{
    Timeout: 15 * time.Second,
  }

	return client.Do(req)
}

func NewBuffer(data []byte) io.Reader {
	return bytes.NewBuffer(data)
}
