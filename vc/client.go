package vc

import (
	"crypto/tls"
	"net/http"
	"time"
)

const (
	LOCALHOSTURL string = "http://127.0.0.1:5000/"
)

// The virtual control server will need to be controlled one of two ways:
//
// - https with a authorization header
// - http local host without auth.
//
// Controlling local host :
// -- requires the /VirtualControl/config/api/ be removed from the route
// -- Local host port 5000
// -- use http NOT https
//
// Controlling external:
// -- full path urls https://[ServerURL]/VirtualControl/config/api/
// -- use of https, I'm not writing this for unsecured servers.
// -- accept self signed certs, this si way too common not too!
// -- header  "Authorization: [Token]"
func createLocalClient() *http.Client {
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

// Creates a remote client with SSL and auth header
func createRemoteClient(token string) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second}

	client.Transport = &headerTransport{
		transport: client.Transport,
		header:    http.Header{"xContent-Type": []string{"application/json"}, "Authorization": []string{token}},
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
