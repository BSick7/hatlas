package terraform

import (
	"fmt"
)

func (c *AtlasClient) GetState(env string) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/terraform/state/%s", env)
	payload, err := c.get(path)
	if err != nil {
		return nil, err
	}
	return payload.Data, nil
}
