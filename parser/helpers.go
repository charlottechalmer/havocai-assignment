package parser

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

func toFloat64(input interface{}, format string) (float64, error) {
	v := reflect.ValueOf(input)
	v = reflect.Indirect(v)
	floatType := reflect.TypeOf(float64(0))

	if v.Kind() == reflect.String {
		dateStr := v.String()
		if parsedDate, err := parseDateString(dateStr, format); err == nil {
			return parsedDate, nil
		}

		if floatVal, err := strconv.ParseFloat(dateStr, 64); err == nil {
			return floatVal, nil
		}
		return 0, fmt.Errorf("cannot conver string %v to float64", dateStr)
	}

	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	floatVal := v.Convert(floatType)
	return floatVal.Float(), nil
}

func parseDateString(dateStr string, format string) (float64, error) {
	date, err := time.Parse(format, dateStr)
	if err != nil {
		return 0, fmt.Errorf("error converting datestring to float64: %v", err)
	}
	return float64(date.Unix()), nil
}

func getFieldValue(field string, record map[string]interface{}, extras map[string]interface{}) (interface{}, bool) {
	if val, found := record[field]; found {
		return val, true
	}
	if val, found := extras[field]; found {
		return val, true
	}
	return nil, false
}

func addValues(values []float64) float64 {
	sum := 0.0
	for _, val := range values {
		sum += val
	}
	return sum
}

func subtractValues(values []float64) float64 {
	result := values[0]
	for i := 1; i < len(values); i++ {
		result -= values[i]
	}
	return result
}

func multiplyValues(values []float64) float64 {
	product := 1.0
	for _, val := range values {
		product *= val
	}
	return product
}

func divideValues(values []float64) (float64, error) {
	result := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] == 0 {
			return 0.0, fmt.Errorf("attempting to divide by 0")
		}
		result /= values[i]
	}
	return result, nil
}

func modValues(values []float64) (float64, error) {
	if len(values) != 2 {
		return 0.0, fmt.Errorf("modulo requires exactly 2 inputs")
	}
	if values[1] == 0 {
		return 0.0, fmt.Errorf("modulo by zero is undefined")
	}
	return math.Mod(values[0], values[1]), nil
}
