package terraform

import (
	"encoding/json"
	"fmt"
	"time"
)

type ListTerraformsResponse struct {
	States []struct {
		UpdatedAt   time.Time `json:"updated_at"`
		Environment struct {
			Username string `json:"username"`
			Name     string `json:"name"`
		} `json:"environment"`
	} `json:"states"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

func (res *ListTerraformsResponse) Names() []string {
	names := []string{}
	for _, state := range res.States {
		names = append(names, fmt.Sprintf("%s/%s", state.Environment.Username, state.Environment.Name))
	}
	return names
}

func (c *AtlasClient) ListTerraforms(username string) (*ListTerraformsResponse, error) {
	path := "/api/v1/terraform/state"
	if username != "" {
		path = path + "/" + username
	}

	payload, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var res ListTerraformsResponse
	if err := json.Unmarshal(payload.Data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *AtlasClient) GetTerraformState(env string) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/terraform/state/%s", env)
	payload, err := c.get(path)
	if err != nil {
		return nil, err
	}
	return payload.Data, nil
}
