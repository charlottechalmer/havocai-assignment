package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"havocai-assignment/models"
	"strconv"
	"strings"
	"time"
)

func ParseXML(input []byte) ([]map[string]interface{}, error) {
	decoder := xml.NewDecoder(bytes.NewReader(input))
	var results []map[string]interface{}

	// to track current data within loop
	currentData := make(map[string]interface{})
	currentElementName := ""
	//track nesting to determine when to exit a grouping
	level := 0

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				// end of XML returns EOF, break from loop
				break
			}
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			// when we encounter a start element, want to increment level (noting we are within an element) and track the current element name
			level++
			currentElementName = t.Name.Local

			// need to handle scenario where a start element has attributes
			for _, attr := range t.Attr {
				currentData[attr.Name.Local] = parseValue(attr.Value)
			}
		case xml.EndElement:
			// when we encounter an end element, want to decrease level (exiting an element) and then append data to results and reset currentData and currentElementName
			level--

			// if level == 1, we are at the end of a grouping
			if level == 1 {
				results = append(results, currentData)
				currentData = make(map[string]interface{})
			}

			currentElementName = ""
		case xml.CharData:
			// when we encounter CharData, store the character data at the current element
			content := strings.TrimSpace(string(t))
			if currentElementName != "" && content != "" {
				currentData[currentElementName] = parseValue(content)
			}
		}
	}
	return results, nil
}

func parseValue(val string) interface{} {
	// try to parse as int
	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	}

	// try to parse as float
	if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		return floatVal
	}

	//try to parse as bool
	if boolVal, err := strconv.ParseBool(val); err == nil {
		return boolVal
	}

	return val
}

func ConvertToJSON(input []map[string]interface{}, cfg *models.Config) ([]byte, error) {
	transformedInput, err := applyTransformations(input, cfg)
	if err != nil {
		return nil, err
	}

	jsonOutput, err := json.MarshalIndent(transformedInput, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonOutput, nil
}

func applyTransformations(input []map[string]interface{}, cfg *models.Config) ([]map[string]interface{}, error) {
	var output []map[string]interface{}

	for _, record := range input {
		transformed := make(map[string]interface{})
		// apply mappings based on 1:1 mapping definition
		for xmlField, jsonField := range cfg.Mappings {
			if val, ok := record[xmlField]; ok {
				transformed[jsonField] = val
			}
		}

		for jsonField, transformation := range cfg.Transformations {
			switch transformation.Type {
			case "concat":
				val, err := concatTransformation(record, transformation)
				if err != nil {
					return nil, err
				}
				transformed[jsonField] = val
			case "calculate":
				val, err := calculateTransformation(record, transformation)
				if err != nil {
					return nil, err
				}
				transformed[jsonField] = val
			}

		}

		output = append(output, transformed)
	}
	return output, nil
}

// what happpens when concating string and int
// e.g. month = "july" and day = 6
// would need to cast to a string
// need to think about if this compromises the integrity of the data
func concatTransformation(record map[string]interface{}, transformation models.Transformation) (string, error) {
	fields := transformation.Params.Fields

	fieldValues := []string{}
	for _, field := range fields {
		// get value from input
		value, ok := record[field]
		if !ok {
			return "", fmt.Errorf("field %v not found in input", field)
		}

		strVal, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("field %v is not a string, cannot concat", field)
		}

		fieldValues = append(fieldValues, strVal)
	}

	separator := ""
	if separatorIface, ok := transformation.Params.Extras["separator"]; ok {
		separator, _ = separatorIface.(string)
	}
	return strings.Join(fieldValues, separator), nil
}

func calculateTransformation(record map[string]interface{}, transformation models.Transformation) (interface{}, error) {
	// get operation from extras
	extras := transformation.Params.Extras
	// ensure operation is a string
	operation, ok := extras["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid operation")
	}

	//parse input fields
	fields := transformation.Params.Fields

	format := "2006-01-02"
	if formatIface, ok := extras["format"]; ok {
		format, _ = formatIface.(string)
	}

	if operation == "time_difference" {
		if len(fields) != 2 {
			return nil, fmt.Errorf("time_difference requires two values")
		}
		startField := fields[0]
		startDateVal, found := getFieldValue(startField, record, extras)
		if !found {
			return nil, fmt.Errorf("field %v not found in xml or in extras", startField)
		}

		startDateStr, ok := startDateVal.(string)
		if !ok {
			return nil, fmt.Errorf("field %v is not a string", startField)
		}

		startDate, err := time.Parse(format, startDateStr)
		if err != nil {
			return nil, fmt.Errorf("unable to parse start date (%v), err %v", startDateStr, err)
		}

		endField := fields[1]
		var endDate time.Time

		if endField == "CurrentTime" {
			endDate = time.Now()
		} else {
			endDateVal, found := getFieldValue(endField, record, extras)
			if !found {
				return nil, fmt.Errorf("field %v not found in xml or in extras", startField)
			}

			endDateStr, ok := endDateVal.(string)
			if !ok {
				return nil, fmt.Errorf("field %v is not a string", endField)
			}

			endDate, err = time.Parse(format, endDateStr)
			if err != nil {
				return nil, fmt.Errorf("unable to parse end date (%v), err %v", endDateStr, err)
			}
		}
		duration := endDate.Sub(startDate)

		unit := "seconds"
		if unitIface, ok := extras["unit"]; ok {
			unit, _ = unitIface.(string)
		}

		switch unit {
		case "years":
			years := float64(duration.Hours()) / (365.25 * 24)
			return years, nil
		case "months":
			months := float64(duration.Hours()) / (30.44 * 24)
			return months, nil
		case "weeks":
			weeks := float64(duration.Hours()) / (7 * 24)
			return weeks, nil
		case "days":
			days := duration.Hours() / 24
			return days, nil
		case "hours":
			return duration.Hours(), nil
		case "minutes":
			return duration.Minutes(), nil
		case "seconds":
			return duration.Seconds(), nil
		default:
			return nil, fmt.Errorf("unsupported time unit: %v", unit)
		}

	}

	values := []float64{}
	for _, field := range fields {
		val, found := getFieldValue(field, record, extras)
		if !found {
			return nil, fmt.Errorf("field %v not found in XML or extras", field)
		}

		floatVal, err := toFloat64(val, format)
		if err != nil {
			return nil, err
		}
		values = append(values, floatVal)
	}
	// switch on operation
	switch operation {
	case "add":
		sum := 0.0
		for _, val := range values {
			sum += val
		}
		return sum, nil
	case "subtract":
		result := values[0]
		for i := 1; i < len(values); i++ {
			result -= values[i]
		}
		return result, nil
	case "multiply":
		product := 1.0
		for _, val := range values {
			product *= val
		}
		return product, nil
	case "divide":
		result := values[0]
		for i := 1; i < len(values); i++ {
			if values[i] == 0 {
				return nil, fmt.Errorf("attempting to divide by 0")
			}
			result /= values[i]
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported operation: %v", operation)
	}
}
