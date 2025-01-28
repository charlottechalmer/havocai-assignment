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
       - TODO: think about how to handle if required field is not in record
     - operation: subtract, add, multiply, divide, time_difference,
     - unit (?)
       - may need unit for time_difference to handle if birthday hasnt passed yet this year when calculating age
*/
