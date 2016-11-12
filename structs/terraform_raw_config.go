package structs

import (
	"bytes"
	"fmt"
)

type TerraformRawConfig struct {
	Version TerraformRawConfigVersion `json:"version"`
}

type TerraformRawConfigVersion struct {
	Version   int                       `json:"version"`
	Metadata  map[string]interface{}    `json:"metadata"`
	TfVars    []TerraformRawConfigTfVar `json:"tf_vars"`
	Variables map[string]string         `json:"variables"`
}

type TerraformRawConfigTfVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Hcl   bool   `json:"hcl"`
}

func (c *TerraformRawConfig) Dump() string {
	buf := bytes.NewBufferString("")
	for k, v := range c.Version.Variables {
		buf.WriteString(fmt.Sprintf("%s = %s\n", k, v))
	}
	return buf.String()
}
