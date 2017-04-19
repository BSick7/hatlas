package terraform

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-rootcerts"
)

type AtlasClient struct {
	config     *AtlasConfig
	httpClient *retryablehttp.Client
}

func NewAtlasClient(config *AtlasConfig) *AtlasClient {
	if config == nil {
		config = DefaultAtlasConfig()
	}

	return &AtlasClient{
		config: config,
	}
}

func (c *AtlasClient) get(path string, query map[string]string) (*Payload, error) {
	return c.do("GET", path, query, nil)
}

func (c *AtlasClient) put(path string, query map[string]string, payload *Payload) (*Payload, error) {
	return c.do("PUT", path, query, payload)
}

func (c *AtlasClient) do(verb string, path string, query map[string]string, data *Payload) (*Payload, error) {
	u, err := c.config.Url(path, query)
	if err != nil {
		return nil, fmt.Errorf("failed to create Atlas URL: %s", err)
	}
	req, err := retryablehttp.NewRequest(verb, u.String(), c.createBody(data))
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %v", err)
	}

	req.Header.Set(atlasTokenHeader, c.config.AccessToken)

	// Request the url
	client, err := c.http()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if os.Getenv("HTTP_DEBUG") != "" && data != nil {
		log.Println(string(data.Data))
	}

	// Handle the common status codes
	switch resp.StatusCode {
	case http.StatusOK:
		// Handled after
	case http.StatusNoContent:
		return nil, nil
	case http.StatusNotFound:
		return nil, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("HTTP remote state endpoint requires auth")
	case http.StatusForbidden:
		return nil, fmt.Errorf("HTTP remote state endpoint invalid auth")
	case http.StatusInternalServerError:
		return nil, fmt.Errorf("HTTP remote state internal server error")
	default:
		return nil, fmt.Errorf(
			"Unexpected HTTP response code: %d\n\nBody: %s",
			resp.StatusCode, c.readBody(resp.Body))
	}

	// Read response
	payload, err := NewPayloadFromResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read remote state: %v", err)
	}

	return payload, nil
}

func (c *AtlasClient) createBody(payload *Payload) io.ReadSeeker {
	if payload == nil {
		return nil
	}
	return payload.GetReader()
}

func (c *AtlasClient) readBody(b io.Reader) string {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, b); err != nil {
		return fmt.Sprintf("Error reading body: %s", err)
	}

	result := buf.String()
	if result == "" {
		result = "<empty>"
	}

	return result
}

func (c *AtlasClient) http() (*retryablehttp.Client, error) {
	if c.httpClient != nil {
		return c.httpClient, nil
	}
	tlsConfig := &tls.Config{}
	err := rootcerts.ConfigureTLS(tlsConfig, &rootcerts.Config{
		CAFile: os.Getenv("ATLAS_CAFILE"),
		CAPath: os.Getenv("ATLAS_CAPATH"),
	})
	if err != nil {
		return nil, err
	}
	rc := retryablehttp.NewClient()

	rc.CheckRetry = func(resp *http.Response, err error) (bool, error) {
		if err != nil {
			// don't bother retrying if the certs don't match
			if err, ok := err.(*url.Error); ok {
				if _, ok := err.Err.(x509.UnknownAuthorityError); ok {
					return false, nil
				}
			}
			// continue retrying
			return true, nil
		}
		return retryablehttp.DefaultRetryPolicy(resp, err)
	}

	t := cleanhttp.DefaultTransport()
	t.TLSClientConfig = tlsConfig
	rc.HTTPClient.Transport = t

	if os.Getenv("HTTP_DEBUG") == "" {
		rc.Logger = log.New(ioutil.Discard, "", 0)
	}

	c.httpClient = rc
	return rc, nil
}
