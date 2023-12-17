package vc

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		header: http.Header{
			//"xContent-Type": []string{"application/json"},
			"Authorization": []string{token}},
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

func (vc *VC) getBody(url string, result any) (err VirtualControlError) {

	resp, err := vc.client.Get(vc.url + url)
	if err != nil {
		return NewServerError(500, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return NewServerError(resp.StatusCode, errors.New("FAILED TO GET DEVICE INFO"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return NewServerError(resp.StatusCode, err)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return NewServerError(resp.StatusCode, err)
	}

	return nil
}

func addFormField(writer *multipart.Writer, key string, value string) {

	fieldWriter, err := writer.CreateFormField(key)
	if err != nil {
		return
	}
	_, err = fieldWriter.Write([]byte(value))
	if err != nil {
		return
	}
}

func addFormFile(filename string, key string, writer *multipart.Writer) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	part, err := writer.CreateFormFile(key, filepath.Base(filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}
	return nil
}
