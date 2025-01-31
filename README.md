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
The current solution leverages dynamic XML parsing and user-provided config.json file that helps offload any domain-specific knowledge to external sources. The XML is parsed in such a way that does not require use of structs defining the shape of the data beforehand. The config file includes entries for direct 1:1 mappings as well as definitions for transformation of input data to produce expected output results (for example, given the input `DateOfBirth`, calculate age). 
The existing transformations were written in a way that enables them to be flexible to future use cases and changes to data and the existing pattern laid out in source code makes it easy to potentially add new transformation types in the future. 

### Running the program
```
go run cmd/main.go -xml test/testdata/provided/input.xml -config test/testdata/basicpatient/config.json
```
The above command will run the program providing the example input XML. The output will be written to `~/Documents/xml-to-json-output/output_{timestamp}.json` unless an optional `-output` flag is provided.
### cmdline flags:
- `-xml` specifies the path to the input xml file
- `-config` specifies the path to the user-created config.json file
- `-output` specifies the path to which the program will write the output json file

### Supported config.json fields
- the config.json file consists of three elements:
	1. the name of the root element to output in json (e.g. `patients`)
	2.  `mappings` that holds 1:1 direct mappings of xml input fields to json output fields
	3.  `transformations` that includes a map of transformations, consisting of the output field name, followed by the transformation definition

#### Transformations
Transformations are maps of `output field name` to a transformation. A transformation consists of the `type` which is used in a switch statement in the `applyTransformations` function that determines how to apply the transformation, and `params`. 
`Params` consists of a list of fields that the transformation works on. These may be fields in the input (for example, `FirstName`, `LastName`, `DateOfBirth`) or fields required to perform the transformation (for example, fields required to be dynamically generated like `CurrentTime`).
`Params` also includes a map of `extras`. `Extras` is intended to store any of the transformation-specific fields required. For example, given a transformation of `"type": "calculate"`, `extras` includes a key value pair that specifies the `operation` of the calculation. Below I have outlined the currently supported transformations and their `extra` params:
 -  type:`concat`
	 - Supported params:
		 - `separator` - defines the separator that the concat transformation should use when combining elements
- type: `calculate`
	- supported params:
		- `operation` - defines the operation of the calculation. Currently supported values for `operation` are: `add`, `subtract`, `multiply`, `divide`, `modulo`, and `time_difference`
			- params specific to `time_difference`:
				- `format` - defines the format of the input dates, required for `time.Parse`. If no `format` is specified, `RFC.33391` is used. If input is not in the specified format, it will throw an error.
				- `unit` - defines the unit of time to output. Supported values are `years`, `months`, `weeks`, `days`, `hours`, `minutes`, `seconds`, `milliseconds`, `microseconds`, and `nanoseconds`. If no unit is specified, the output defaults to `seconds`.
				- `adjust_if_day_not_passed` - boolean value used specifically for age calculation to adjust for the case if the person's birthday has not passed yet this year
				- `round_to_int` - boolean value used to round float64 value to nearest int value
				- `decimal_precision` - int value specifying the number of decimal places to round to

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
Obviously, the simple solution would not be able to extend to meet all requirements. I needed some way to remove all domain-specific information from the source code including any structs that defined input and output structure, and any field-specific transformation functions. I achieved this by parsing the XML into a `[]map[string]interface{}` and maintaining a `level` variable that handled the nesting present in the sample input data. I then leveraged a config file that off-loaded any of the transformation logic to keep source code extensible.
- on branch ["dynamic-parsing"](https://github.com/charlottechalmer/havocai-assignment/tree/dynamic-parsing) (and eventually `main`).

#### Supporting Nested Stuctures
While the logic on branch `dynamic-parsing` works for the current input, I began to think about how the input data could possibly change. One of the scenarios I imagined was support for more highly nested XML structures. While I initially thought supporting this wouldn't be too bad, it turned out to be quite a doozy. I spent over 15 hours attempting to wrangle my code to allow for this. I decided that the best way to support nested data and maintain the order while looping through the xml tokens was to leverage a list of maps and treat them like a stack, ensuring FIFO. However, this did not account for the fact that as we are processing the XML elements, when we encounter an EndElement, there is no easy way for us to know which element is closing. This resulted in strangely nested outputs or outputs where each element was its own entry in the map. 
I attempted to offset this by maintaining two different data structures -- one being the same `[]map[string]interface{}` that stored the content of each XML node as before, and one `[]string` that is also treated as a stack that simply maintained the name of the field we were processing. This then would enable me to process both of the stacks at the same time, allowing me to correlate the node content with the element name upon encountering an EndElement. 
After many hours of attempting to make this work, I decided that, because my initial solution worked for the sample data and because this is more to support data extensibility than core functionality at the moment, it should be sufficient to push my commit with the stack implementation to a new branch, in the case the input data structure changed in the future.
- in the future, rather than bang my head against the wall, I would likely consider leveraging external packages that handle this for me. For example [mxj](https://pkg.go.dev/github.com/clbanning/mxj)
- on branch "[support-nested-structure-attempt](https://github.com/charlottechalmer/havocai-assignment/tree/support-nested-structure-attempt)".

### Making things generic
#### Parsing XML
##### xml.Unmarshal(data, &generic)
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

##### parsing input XML into a generic type
I needed to be able to parse the input XML without defining a struct. 
This sent me down a Google rabbithole where various people were attempting to solve the same issue. I found  [This article](https://medium.com/@jhxjohn/dynamic-xml-parse-to-json-in-go-lang-e8196752e47f) that describes someone wanting to do something similar by using [this library](https://github.com/antchfx/xmlquery), however, I wasn't sure if using libraries outside of the standard Go library was alright, so I reached out to Matt for confirmation.

Not wanting to be in a holding pattern while I waited for his response, I started down the path of manually parsing the input XML into a generic type. [This Google conversation](https://groups.google.com/g/golang-nuts/c/BRpabwxVrOc) led me down the right path and I began to write the logic using `xml.Decoder()`, and then calling `decoder.Token()` within a loop to iterate over each of the elements within the XML document and handle `StartElement`, `EndElement`, and  `CharData` tokens.

I found myself getting a bit convoluted in the branches of logic in this method, so I attempted to offset that by adding comments to the `ParseXML` function to ensure that I was properly handing each scenario. See [here](https://github.com/charlottechalmer/havocai-assignment/blob/dynamic-parsing/parser/parser.go#L14). I also made sure to update tests as I went so I could confirm I was properly replacing portions of the brute-force solution.
  
#### offloading transformation logic
##### config structure
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
 
 3. I modified the Transformation struct to allow for a `params` field of type `Params`, thus enabling further extensibility of transformations:
```go
type Config struct {
	Mappings        map[string]string         `json:"mappings"`
	Transformations map[string]Transformation `json:"transformations"`
}

type Transformation struct {
	Type   string            `json:"type"`
	Params Params `json:"params"`
}

type Params struct {
	Fields []string `json:"fields"`
	Extras map[string]interface{} `json:"extras"`
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
              "extras": {
	              "separator": " "
	          }
          }
      }
  }
```
We now can use this data structure to build out greater functionality of the transformations, e.g. allowing for different data types in params, adding a `separator` param for the `concat` transformation that dictates how to combine the strings in `fields`. While the code now requires many more type assertions, making it a bit less readable, this is a necessary trade-off in order to ensure extensibility of our config file.
This was a required change especially when trying to make the transformation responsible for transforming `DateOfBirth` into `age` that required a number of "extra" parameters in order to properly handle the transformation.

### Future Considerations
- the current implementation should be generic enough that changes to input structure or output requirements should require changes to the existing configs and some minor changes to add new types of transformations. If nested structures are needed in the future, a new data structure in `ParseXML()` will be needed (see note about that above). While I tried to make this as generic as possible, it was near impossible to address every possible change to data. I have outlined some of the changes that I would imagine could be possible in the future and how they may be handled:
	- changes to input data:
		- new fields are added or fields are renamed
			- should only require changes to the config.json file
		- new transformation types
			- formatting changes
				-  for example:
					- date format: DateOfBirth YYYY-MM-DD --> DD-MM-YYYY
					- phone numbers: 888-555-1234 --> (888) 555-1234
					- padding numbers: 23 --> 00023
					- numbers to money: 1000.5 --> $1000.50
			  - To support formatting changes, we would need a new type of transformation called `format`. This would require quite a few additional params in `extras` to support each of these cased. For example, to support changes to the date format, we would need some key to tell the program that we are handing `time.Time` as well as a key that specifies the new layout. Then, we could include logic that parses that value and then calls `time.Format(layout)`, passing in the layout.
			  - The other transformations here could be done by processing the number input and iterating over each element, either performing some transformation by inserting characters, ensuring to handle edge cases/
			- boolean conditions based on existence of other field
				- for example:
					- output `deceased: true` if `DateOfDeath` is present in the input
				- this could be handled with a new transformation type that takes an optional field and returns a boolean depending on if it is present or not
			- counting
				- for example
					- if input data specifies a list of allergies, output returns `allergy_count: int`
				- this transformation would require support of nested structures and parsing input into an array. Once that is complete, we could simply return `len(inputField)` for the specified field.
	- changes to output data requirements
		- field names change: should only require changes to config.json
		- changing concat format (name --> Last, First): should require changing order of input in config.json and updating separator
		- field types change (age(int) --> string): update `extras` to include `output_type` and return that output type from transformation 
