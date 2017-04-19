package terraform

import (
	"net/url"
	"os"
)

const (
	// defaultAtlasServer is used when no address is given
	defaultAtlasServer = "https://atlas.hashicorp.com/"
	atlasTokenHeader   = "X-Atlas-Token"
)

type AtlasConfig struct {
	Address     string
	AccessToken string
}

func DefaultAtlasConfig() *AtlasConfig {
	server := os.Getenv("ATLAS_ADDR")
	if server == "" {
		server = defaultAtlasServer
	}

	return &AtlasConfig{
		Address:     server,
		AccessToken: os.Getenv("ATLAS_TOKEN"),
	}
}

func (ac *AtlasConfig) Url(path string, query map[string]string) (*url.URL, error) {
	u, err := url.Parse(ac.Address)
	if err != nil {
		return nil, err
	}
	u.Path = path
	params := url.Values{}
	for k, v := range query {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()
	return u, nil
}
