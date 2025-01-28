package models

type Config struct {
	Mappings        map[string]string         `json:"mappings"`
	Transformations map[string]Transformation `json:"transformations"`
}

type Transformation struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params"`
}

/*
 examples of transformations & info needed
 - concat
   - need: list of elements to concat
 - calculate
   - need
     - operator: +, -, \, *, %
     - input:

*/
