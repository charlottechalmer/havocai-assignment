
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
	- this essentially requires the program to forego use of predefined structs that define the shape of input and output data

### Attempts
1. To start, I simply ignored the requirements for extensibility and created a program that took pre-defined XML data and transformed it into expected JSON structure. 
	- this included use of predefined structs as well as transformation functions to apply transformations to fields like `age`.
	- on branch "[brute-force](https://github.com/charlottechalmer/havocai-assignment/tree/brute-force)"
2. Next, I decided that I wanted to leverage a config file that could be loaded at runtime and would define field mappings as well as any transformations required on that field. 
  - I began by creating a Config struct that would consist of a list of mappings for each field:
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
 3. I reevaluated the Config struct and decided a better approach would be to include separate entries for the mappings as well as the transformation rules:
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
This then would allow for more extensibility, 
