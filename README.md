

# havocai-assignment

## Assignment

Write a program in Go that transforms XML data into JSON. Your solution should allow for configurable field mappings and transformations without requiring changes to the program's source code

### Example Input
```xml
<?xml  version="1.0"  encoding="UTF-8"?>

<Patients>
  <Patient  ID="12345">
    <FirstName>John</FirstName>
    <LastName>Doe</LastName>
    <DateOfBirth>1985-07-15</DateOfBirth>
  </Patient>
  <Patient  ID="67890">
    <FirstName>Jane</FirstName>
    <LastName>Smith</LastName>
    <DateOfBirth>1992-03-22</DateOfBirth>
  </Patient>
</Patients>
```
### Example Output
```json
{
  "patients":  [
    {
      "id":  12345,
      "name":  "John  Doe",
      "age":  39
    },
    {
      "id":  67890,
      "name":  "Jane  Smith",
      "age":  32
    }
  ]
}
```
### Considerations:
-  The program should not require hardcoding of field mappings or transformations.
- Think about how you would handle future changes to the input XML structure or output requirements without modifying the core logic.
-  Ensure your solution is efficient and easy to extend.

## Solution
### Explanation 
lorem ipsum 

## Thought Process
- transforming data from XML to JSON is easy enough. The difficult thing here is to be able to do so without requiring changes to source code in the case that the input structure changes or output requirements change.
	- this requires the program to:
		-  forego use of predefined structs that define the shape of input and output data
		- offload transformation specifics from source code to some external mechanism

### Attempts
#### Simple "brute-force" Solution
To start, I simply ignored the requirements for extensibility and created a program that took pre-defined XML data and transformed it into expected JSON structure. 
	- this included use of predefined structs as well as transformation functions to apply transformations to fields like `age`.
	- on branch "[brute-force](https://github.com/charlottechalmer/havocai-assignment/tree/brute-force)"

#### Dynamic Parsing
Obviously, the simple solution would not be able to extend to meet all requirements. I needed some way to remove all domain-specific information from the source code including any structs that defined input and output structure, and any field-specific transformation functions.
- on branch ["dynamic-parsing"](https://github.com/charlottechalmer/havocai-assignment/tree/dynamic-parsing) (and eventually `main`).

##### Making things generic
###### xml.Unmarshal(data, &generic)
I began by removing any references to structs that defined the XML input and JSON output and started by adding logic to parse XML into a generic type like `map[string]interface{}`. 
Being fairly unfamiliar with the XML spec, I started by updating the `parseXML` function to simply unmarshal into a generic type:
```go
func ParseXML(input []byte) (map[string]interface{}, error) {
	var generic map[string]interface{} 
	err := xml.Unmarshal(data, &generic)
	if err != nil { 
		return nil, err
	}
	return genericXML, nil
}
```
which, as per the [xml.Unmarshal docs](https://pkg.go.dev/encoding/xml#Unmarshal) results in an error: `unknown type map[string]interface{}`.

###### parsing input XML into a generic type
I needed to be able to parse the input XML without defining a struct. 
This sent me down a Google rabbithole where various people were attempting to solve the same issue. I found  [This article](https://medium.com/@jhxjohn/dynamic-xml-parse-to-json-in-go-lang-e8196752e47f) that describes someone wanting to do something similar by using [this library](https://github.com/antchfx/xmlquery), however, I wasn't sure if using libraries outside of the standard Go library was alright, so I reached out to Matt for confirmation.

Not wanting to be in a holding pattern while I waited for his response, I started down the path of manually parsing the input XML into a generic type. [This Google conversation](https://groups.google.com/g/golang-nuts/c/BRpabwxVrOc) led me down the right path and I began to write the logic using `xml.Decoder()`, and then calling `decoder.Token()` within a loop to iterate over each of the elements within the XML document and handle `StartElement`, `EndElement`, and  `CharData` tokens.

I found myself getting a bit convoluted in the branches of logic in this method, so I attempted to offset that by adding comments to the `ParseXML` function to ensure that I was properly handing each scenario. See [here](https://github.com/charlottechalmer/havocai-assignment/blob/dynamic-parsing/parser/parser.go#L14). I also made sure to update tests as I went so I could confirm I was properly replacing portions of the brute-force solution.
  
##### offloading transformation logic
###### config structure
Next, I decided that I wanted to leverage a config file that could be loaded at runtime and would define field mappings as well as any transformations required on that field. 

Attempts:	
1. I began by creating a Config struct that would consist of a list of mappings for each field:
```go
type Config struct {
	Mappings []Mapping `json:"mappings"`
}

type Mapping struct {
	XMLPath string `json:"xml_path"`
	JSONField string `json:"json_field"`
	TransformRule string `json:"transform_rule"`
	TransformField []string `json:"transform_fields"`
}
```
  - I quickly found an issue with this approach because the XML paths do not have a one to one relationship with the JSON fields. For example, the XML input includes both `<FirstName>` and `<LastName>`, which results in `name` in the JSON output. This would mean that by using this structure for the config file, I would need to include a "mapping" for every field in both the XML and the JSON, resulting in a lot of unused or unneeded objects:
```json
{
  "xml_path": "FirstName",
  "json_field": "first_name",
  "transform_rule": "",
  "transform_fields": [] 
}, { 
  "xml_path": "LastName",
  "json_field": "last_name",
  "transform_rule": "",
  "transform_fields": []
}, { 
  "xml_path": "",
  "json_field": "name",
  "transform_rule": "concat",
  "transform_fields": ["FirstName", "LastName"]
},
```

2. I reevaluated the Config struct and decided a better approach would be to include separate entries for the mappings and the transformation rules:
```go
type Config struct {
	Mappings map[string]string `json:"mappings"`
	Transformations map[string]Transformation `json:"transformations"`
}

type Transformation struct {
	Type string `json:"type"`
	Fields []string `json:"fields"`
}
```
This then would allow for more extensibility by not requiring a 1:1 relationship between xml paths and json fields, nor between fields and transformations. Additionally, we will be able to process fields that have a 1:1 relationship between XML and JSON by looping over the `Mappings`, and then process transformations separately. 

This data structure worked great when attempting to define transformations for simple concating of strings, like `name`:
```json
{
	"transformations": {
		"name": {
			"type": "concat",
			"fields": [
				"FirstName",
				"LastName"
			]
		}
	}
}
```
```go
func concat(record map[string]interface{}, transformation models.Transformation) string {
	output := ""
	for _, field := range transformation.Fields {
		output += record[field].(string) + " "
	}
	return output
}
```
However, we are quite limited by the Transformation data structure:
	- we would never be able to concat with a separator other than a space
	- this assumes that any other transformation would only need field names in order to apply transformations successfully
 
 3. I modified the Transformation struct to allow for a `params` field of type `map[string]interface{}`, thus enabling further extensibility of transformations:
```go
type Config struct {
	Mappings        map[string]string         `json:"mappings"`
	Transformations map[string]Transformation `json:"transformations"`
}

type Transformation struct {
	Type   string            `json:"type"`
	Params map[string]interface{} `json:"params"`
}
```
```json
  "transformations": {
      "name": {
          "type": "concat",
          "params": {
              "fields": [
	              "FirstName",
	              "LastName"
	          ],
              "separator": " "
          }
      }
  }
```
we now can use this data structure to build out greater functionality of the transformations, e.g. allowing for different data types in params, adding a `separator` param for the `concat` transformation that dictates how to combine the strings in `fields`.
This was a required change especially when trying to make the transformation responsible for transforming `DateOfBirth` into `age`.

