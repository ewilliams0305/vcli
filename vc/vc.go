package vc

import (
	"net/http"
	"time"
)

const (
	HOSTURL    string = "http://127.0.0.1:5000/"
	DEVICEINFO string = "DeviceInfo"
)

type VirtualControl interface {
	GetDeviceInfo() (DeviceInfo, error)
}

var client *http.Client = createClient()

func createClient() *http.Client {
	tr := &http.Transport{}

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second}

	client.Transport = &headerTransport{
		transport: client.Transport,
		header:    http.Header{"xContent-Type": []string{"application/json"}},
	}

	return client
}

// headerTransport is a custom transport that adds headers to each request
type headerTransport struct {
	transport http.RoundTripper
	header    http.Header
}

// RoundTrip adds the custom header to the request and delegates the actual
// request to the underlying transport.
func (t *headerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add custom headers to the request
	for key, values := range t.header {
		req.Header[key] = values
	}

	return t.transport.RoundTrip(req)
}

func makeUrl(url string) string {
	return HOSTURL + url
}
