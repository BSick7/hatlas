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
	return &AtlasConfig{
		Address:     defaultAtlasServer,
		AccessToken: os.Getenv("ATLAS_TOKEN"),
	}
}

func (ac *AtlasConfig) Url(path string) (*url.URL, error) {
	u, err := url.Parse(ac.Address)
	if err != nil {
		return nil, err
	}
	u.Path = path
	return u, nil
}
