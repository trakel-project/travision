package rpc

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client defines the interface for go client that wants to connect to a hyperchain RPC endpoint
type client interface {
	Send(body []byte) ([]byte, error)
	Close()
}

// httpClient connects to a hyperchain RPC server over HTTP.
type httpClient struct {
	endpoint *url.URL    // HTTP-RPC server endpoint
	client   *http.Client // reuse connection
}

// NewHTTPClient create a new RPC client that connection to
// a hyperchain RPC server over HTTP.
func newHTTPClient(endpoint string, timeout time.Duration) (client, error) {
	url, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: timeout,
	}
	return &httpClient{endpoint: url, client: client}, nil
}

// Send will serialize the given req to JSON and sends it to the RPC server.
// If receive response with statusOK(200), return []byte of response body.
func (c *httpClient) Send(body []byte) ([]byte, error) {

	resp, err := c.client.Post(c.endpoint.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ioutil.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("http failed: %s", resp.Status)
}

// Close is not necessary for httpClient
func (c *httpClient) Close() {
}
