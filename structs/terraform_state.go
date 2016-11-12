package structs

type TerraformState struct {
	Version          int                     `json:"version,omitempty"`
	TerraformVersion string                  `json:"terraform_version,omitempty"`
	Serial           int                     `json:"serial,omitempty"`
	Lineage          string                  `json:"lineage,omitempty"`
	Remote           *TerraformStateRemote   `json:"remote,omitempty"`
	Modules          []*TerraformStateModule `json:"modules"`
}

type TerraformStateRemote struct {
	Type   string `json:"type,omitempty"`
	Config struct {
		Name string `json:"name,omitempty"`
	} `json:"config,omitempty"`
}

type TerraformStateModule struct {
	Path      []string                         `json:"path,omitempty"`
	Outputs   map[string]*TerraformStateOutput `json:"outputs,omitempty"`
	Resources map[string]interface{}           `json:"resources,omitempty"`
	DependsOn []interface{}                    `json:"depends_on"`
}

type TerraformStateOutput struct {
	Sensitive bool        `json:"sensitive,omitempty"`
	Type      string      `json:"type,omitempty"`
	Value     interface{} `json:"value,omitempty"`
}

func (s *TerraformState) GetRootModule() *TerraformStateModule {
	for _, module := range s.Modules {
		if len(module.Path) == 1 && module.Path[0] == "root" {
			return module
		}
	}
	return nil
}
