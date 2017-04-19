package structs

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type TerraformRawConfig struct {
	Version TerraformRawConfigVersion `json:"version"`
}

type TerraformRawConfigVersion struct {
	Version   int                      `json:"version"`
	Metadata  map[string]interface{}   `json:"metadata"`
	TfVars    TerraformRawConfigTfVars `json:"tf_vars"`
	Variables map[string]string        `json:"variables"`
}

type TerraformRawConfigTfVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Hcl   bool   `json:"hcl"`
}

type TerraformRawConfigTfVars []TerraformRawConfigTfVar

func (vars TerraformRawConfigTfVars) CreateMap() map[string]interface{} {
	data := map[string]interface{}{}
	for _, v := range vars {
		data[v.Key] = data[v.Value]
	}
	return data
}

func NewTerraformRawConfigFromJson(raw []byte) (*TerraformRawConfig, error) {
	trc := &TerraformRawConfig{}
	err := json.Unmarshal(raw, trc)
	return trc, err
}

func (c *TerraformRawConfig) Dump() string {
	buf := bytes.NewBufferString("")
	for _, tfvar := range c.Version.TfVars {
		buf.WriteString(tfvar.Dump())
	}
	return buf.String()
}

func (c *TerraformRawConfig) DumpRaw() string {
	raw, _ := json.Marshal(c)
	out := bytes.NewBufferString("")
	if err := json.Indent(out, raw, "", "  "); err != nil {
		return string(raw)
	}
	return out.String()
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
