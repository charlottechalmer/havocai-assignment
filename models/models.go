package models

type Config struct {
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

// type Transformation struct {
// 	Type   string                 `json:"type"`
// 	Params map[string]interface{} `json:"params"`
// }

/*
 examples of transformations & info needed
 - concat
   - need:
     - fields: list of elements to concat
 - calculate
   - need
     - fields: list of elemetns for calculation
     - operation: subtract, add, multiply, divide, time_difference,

*/
