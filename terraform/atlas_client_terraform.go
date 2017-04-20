package terraform

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/BSick7/hatlas/structs"
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

func (c *AtlasClient) ListTerraforms(username string, page int) (*ListTerraformsResponse, error) {
	path := "/api/v1/terraform/state"
	query := map[string]string{}
	if username != "" {
		query["username"] = username
	}
	if page > 1 {
		query["page"] = fmt.Sprintf("%d", page)
	}

	payload, err := c.get(path, query)
	if err != nil {
		return nil, err
	}

	var res ListTerraformsResponse
	if payload != nil {
		if err := json.Unmarshal(payload.Data, &res); err != nil {
			return nil, err
		}
	}
	return &res, nil
}

func (c *AtlasClient) GetTerraformState(env string) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/terraform/state/%s", env)
	payload, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}
	return payload.Data, nil
}

func (c *AtlasClient) GetTerraformConfig(env string) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/terraform/configurations/%s/versions/latest", env)
	payload, err := c.get(path, nil)
	if err != nil {
		return nil, err
	}
	return payload.Data, nil
}

type UpdateVariablesRequest struct {
	Variables map[string]string `json:"variables"`
}

func (req *UpdateVariablesRequest) ToPayload() *Payload {
	raw, _ := json.Marshal(req)
	payload := NewPayloadFromString(string(raw))
	payload.ContentType = "application/json"
	return payload
}

func (c *AtlasClient) UpdateVariables(env string, config *structs.TerraformRawConfig) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/environments/%s/variables", env)
	req := &UpdateVariablesRequest{
		Variables: config.GetVarsMap(),
	}

	if err := c.put(path, nil, req.ToPayload()); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("pushed %d variables to %s", len(req.Variables), env)), nil
}

type UpdateStateRequest struct {
	State []byte
}

func (req *UpdateStateRequest) ToPayload() *Payload {
	payload := NewPayloadFromBytes(req.State)
	payload.ContentType = "application/json"
	return payload
}

func (c *AtlasClient) UpdateState(env string, state []byte) ([]byte, error) {
	path := fmt.Sprintf("/api/v1/terraform/state/%s", env)
	req := &UpdateStateRequest{
		State: state,
	}

	if err := c.put(path, nil, req.ToPayload()); err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("pushed state to %s", env)), nil
}
