package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"havocai-assignment/models"
	"math"
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

	wrappedOutput := map[string]interface{}{
		cfg.RootName: transformedInput,
	}

	jsonOutput, err := json.MarshalIndent(wrappedOutput, "", "  ")
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

func concatTransformation(record map[string]interface{}, transformation models.Transformation) (string, error) {
	fields := transformation.Params.Fields

	fieldValues := []string{}
	for _, field := range fields {
		value, ok := record[field]
		if !ok {
			return "", fmt.Errorf("field %v not found in input", field)
		}

		strVal := fmt.Sprintf("%v", value)

		fieldValues = append(fieldValues, strVal)
	}

	separator := ""
	if separatorIface, ok := transformation.Params.Extras["separator"]; ok {
		separator, _ = separatorIface.(string)
	}
	return strings.Join(fieldValues, separator), nil
}

func calculateTransformation(record map[string]interface{}, transformation models.Transformation) (interface{}, error) {
	extras := transformation.Params.Extras
	operation, ok := extras["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid operation")
	}

	fields := transformation.Params.Fields

	format := time.RFC3339
	if formatIface, ok := extras["format"]; ok {
		format, _ = formatIface.(string)
	}

	if operation == "time_difference" {
		return calculateTimeDifference(fields, record, extras, format)
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
		return addValues(values), nil
	case "subtract":
		return subtractValues(values), nil
	case "multiply":
		return multiplyValues(values), nil
	case "divide":
		return divideValues(values)
	case "modulo":
		return modValues(values)
	default:
		return nil, fmt.Errorf("unsupported operation: %v", operation)
	}
}

func parseDate(field string, record map[string]interface{}, extras map[string]interface{}, format string) (time.Time, error) {
	val, found := getFieldValue(field, record, extras)
	if !found {
		return time.Time{}, fmt.Errorf("field %v not found in xml or in extras", field)
	}

	dateStr, ok := val.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("field %v is not a string", field)
	}

	return time.Parse(format, dateStr)
}

func calculateTimeDifference(fields []string, record map[string]interface{}, extras map[string]interface{}, format string) (interface{}, error) {
	if len(fields) != 2 {
		return nil, fmt.Errorf("time_difference requires two values")
	}
	startField := fields[0]
	startDate, err := parseDate(startField, record, extras, format)
	if err != nil {
		fmt.Println("err here")
		return nil, err
	}

	endField := fields[1]
	var endDate time.Time

	if endField == "CurrentTime" {
		endDate = time.Now()
	} else {
		endDate, err = parseDate(endField, record, extras, format)
		if err != nil {
			return nil, err
		}
	}

	unit := "seconds" //default unit
	if u, ok := extras["unit"]; ok {
		unit, _ = u.(string)
	}

	return calculateDuration(startDate, endDate, unit, extras)
}

func calculateDuration(startDate time.Time, endDate time.Time, unit string, extras map[string]interface{}) (interface{}, error) {
	var result float64
	duration := endDate.Sub(startDate)

	switch unit {
	case "years":
		result = float64(endDate.Year() - startDate.Year())

		//check if need to adjust if day in question has not passed yet this year
		if adjust, ok := extras["adjust_if_day_not_passed"].(bool); ok && adjust {
			if endDate.YearDay() < startDate.YearDay() {
				result--
			}
		}
		return result, nil
	case "months":
		result = duration.Hours() / (30.44 * 24)
	case "weeks":
		result = duration.Hours() / (7 * 24)
	case "days":
		result = duration.Hours() / 24
	case "hours":
		result = duration.Hours()
	case "minutes":
		result = duration.Minutes()
	case "seconds":
		result = duration.Seconds()
	case "milliseconds":
		result = float64(duration.Milliseconds())
	case "microseconds":
		result = float64(duration.Microseconds())
	case "nanoseconds":
		result = float64(duration.Nanoseconds())
	default:
		return nil, fmt.Errorf("unsupported time unit: %v", unit)
	}

	if precision, ok := extras["decimal_precision"].(int); ok {
		multiplier := math.Pow(10, float64(precision))
		return math.Round(result*multiplier) / multiplier, nil
	}

	if round, ok := extras["round_to_int"].(bool); ok && round {
		return int(math.Round(result)), nil
	}

	return result, nil
}
