package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToFloat(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		format      string
		expected    float64
		expectedErr bool
	}{
		{
			name:        "int to float64",
			input:       123,
			format:      "",
			expected:    123.0,
			expectedErr: false,
		},
		{
			name:        "string float64 to float64",
			input:       "123.45",
			format:      "",
			expected:    123.45,
			expectedErr: false,
		},
		{
			name:        "datetime to float64",
			input:       "2023-01-01",
			format:      "2006-01-02",
			expected:    1672531200,
			expectedErr: false,
		},
		{
			name:        "invalid string input",
			input:       "invalid",
			format:      "",
			expected:    0,
			expectedErr: true,
		},
		{
			name:        "invalid bool input",
			input:       true,
			format:      "",
			expected:    0,
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := toFloat64(test.input, test.format)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestParseDateString(t *testing.T) {
	tests := []struct {
		name        string
		dateStr     string
		format      string
		expected    float64
		expectedErr bool
	}{
		{
			name:        "valid date",
			dateStr:     "2025-01-29",
			format:      "2006-01-02",
			expected:    1738108800,
			expectedErr: false,
		},
		{
			name:        "invalid date",
			dateStr:     "invalid",
			format:      "2006-01-02",
			expected:    0,
			expectedErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := parseDateString(test.dateStr, test.format)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		record     map[string]interface{}
		extras     map[string]interface{}
		expected   interface{}
		expectedOk bool
	}{
		{
			name:  "field exists in record",
			field: "FirstName",
			record: map[string]interface{}{
				"FirstName": "John",
			},
			extras:     map[string]interface{}{},
			expected:   "John",
			expectedOk: true,
		},
		{
			name:   "field exists in extras",
			field:  "age",
			record: map[string]interface{}{},
			extras: map[string]interface{}{
				"age": 30,
			},
			expected:   30,
			expectedOk: true,
		},
		{
			name:       "field does not exist",
			field:      "imaginary",
			record:     map[string]interface{}{},
			extras:     map[string]interface{}{},
			expected:   nil,
			expectedOk: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, ok := getFieldValue(test.field, test.record, test.extras)
			if actual != test.expected || ok != test.expectedOk {
				t.Errorf("expected (%v %v), got (%v %v)", test.expected, test.expectedOk, actual, ok)
			}
		})
	}
}

func TestMathOperations(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		values      []float64
		expected    float64
		expectedErr bool
	}{
		{
			name:        "addition",
			operation:   "add",
			values:      []float64{1.5, 2.5, 3.0},
			expected:    7.0,
			expectedErr: false,
		},
		{
			name:        "subtraction",
			operation:   "subtract",
			values:      []float64{10.0, 3.0, 2.0},
			expected:    5.0,
			expectedErr: false,
		},
		{
			name:        "multiplication",
			operation:   "multiply",
			values:      []float64{2.0, 3.0, 4.0},
			expected:    24.0,
			expectedErr: false,
		},
		{
			name:        "valid division",
			operation:   "divide",
			values:      []float64{10.0, 2.0, 5.0},
			expected:    1.0,
			expectedErr: false,
		},
		{
			name:        "division by zero",
			operation:   "divide",
			values:      []float64{10.0, 0.0},
			expected:    0,
			expectedErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var actual float64
			var err error

			switch test.operation {
			case "add":
				actual = addValues(test.values)
			case "subtract":
				actual = subtractValues(test.values)
			case "multiply":
				actual = multiplyValues(test.values)
			case "divide":
				actual, err = divideValues(test.values)
			}
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}
