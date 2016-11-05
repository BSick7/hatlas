package terraform

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/go-rootcerts"
	"io"
	"net/http"
	"net/url"
	"os"
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
	u, err := c.config.Url(path, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Atlas URL: %s", err)
	}
	req, err := retryablehttp.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to make HTTP request: %v", err)
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

	// Read in the body
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return nil, fmt.Errorf("Failed to read remote state: %v", err)
	}

	// Create the payload
	payload := &Payload{
		Data: buf.Bytes(),
	}

	if len(payload.Data) == 0 {
		return nil, nil
	}

	// Check for the MD5
	if raw := resp.Header.Get("Content-MD5"); raw != "" {
		md5, err := base64.StdEncoding.DecodeString(raw)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode Content-MD5 '%s': %v", raw, err)
		}

		payload.MD5 = md5
	} else {
		// Generate the MD5
		hash := md5.Sum(payload.Data)
		payload.MD5 = hash[:]
	}

	return payload, nil
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

	c.httpClient = rc
	return rc, nil
}
