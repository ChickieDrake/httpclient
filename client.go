/*
Package httpclient provides a simple interface for making http calls
*/
package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const wrongStatusErrorFormat = "apiclient: Received non-200 status code from api endpoint: %d"

// ApiClient provides a wrapper for making http calls to AnkiConnect.
type ApiClient struct {
	httpClient httpClient
	bodyReader bodyReader
}

// New creates a pointer to a new instance of ApiClient.
func New() *ApiClient {
	return &ApiClient{
		httpClient: new(http.Client),
		bodyReader: new(defaultBodyReader),
	}
}

// DoPost takes a well-formatted JSON message and sends it to AnkiConnect.
func (a *ApiClient) DoPost(uri string, body string) (string, error) {
	return a.doHTTP(uri, http.MethodPost, body)
}

// DoGet returns the response body of an HTTP GET request
func (a *ApiClient) DoGet(uri string) (string, error) {
	return a.doHTTP(uri, http.MethodGet, "")
}

func (a *ApiClient) doHTTP(uri string, method string, body string) (string, error) {
	var message string
	var resp *http.Response
	var err error

	if method == http.MethodPost {
		mimeType := "application/json"
		resp, err = a.httpClient.Post(uri, mimeType, bytes.NewBufferString(body))
	} else if method == http.MethodGet {
		resp, err = a.httpClient.Get(uri)
	}
	if err != nil {
		return message, err
	}

	if resp.StatusCode != 200 {
		errMessage := fmt.Sprintf(wrongStatusErrorFormat, resp.StatusCode)
		err = errors.New(errMessage)
		return message, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
		body, readErr := a.bodyReader.ReadAll(resp.Body)
		if readErr != nil {
			err = readErr
			return message, err
		}

		message = string(body)
	}

	return message, err
}

//go:generate mockery -name httpClient -filename mock_http_test.go -structname MockHTTPClient -output . -inpkg
type httpClient interface {
	Post(uri, contentType string, body io.Reader) (*http.Response, error)
	Get(uri string) (*http.Response, error)
}

type bodyReader interface {
	ReadAll(r io.Reader) ([]byte, error)
}

type defaultBodyReader struct{}

func (*defaultBodyReader) ReadAll(r io.Reader) ([]byte, error) {
	return ioutil.ReadAll(r)
}
