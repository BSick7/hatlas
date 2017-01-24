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
	for _, tfvar := range c.Version.TfVars {
		buf.WriteString(tfvar.Dump())
	}
	return buf.String()
}

func (c *TerraformRawConfig) DumpKey(key string) string {
	return fmt.Sprintf("%s = %s\n", key, c.Version.Variables[key])
}

func (v *TerraformRawConfigTfVar) Dump() string {
	if v.Hcl {
		return fmt.Sprintf("%s = %s\n", v.Key, v.Value)
	} else {
		return fmt.Sprintf("%s = %q\n", v.Key, v.Value)
	}
}
