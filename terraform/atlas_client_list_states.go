package terraform

import (
	"encoding/json"
	"fmt"
	"time"
)

type StateListResponse struct {
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

func (res *StateListResponse) Names() []string {
	names := []string{}
	for _, state := range res.States {
		names = append(names, fmt.Sprintf("%s/%s", state.Environment.Username, state.Environment.Name))
	}
	return names
}

func (c *AtlasClient) ListStates(username string) (*StateListResponse, error) {
	path := "/api/v1/terraform/state"
	if username != "" {
		path = path + "/" + username
	}

	payload, err := c.get(path)
	if err != nil {
		return nil, err
	}

	var res StateListResponse
	if err := json.Unmarshal(payload.Data, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
