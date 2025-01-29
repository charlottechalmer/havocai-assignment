package models

type Config struct {
	RootName        string                    `json:"root"`
	Mappings        map[string]string         `json:"mappings"`
	Transformations map[string]Transformation `json:"transformations"`
}

type Transformation struct {
	Type   string `json:"type"`
	Params Params `json:"params"`
}

type Params struct {
	Fields []string               `json:"fields"`
	Extras map[string]interface{} `json:"extras"`
}
