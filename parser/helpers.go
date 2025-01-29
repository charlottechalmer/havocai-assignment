package parser

import (
	"fmt"
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
