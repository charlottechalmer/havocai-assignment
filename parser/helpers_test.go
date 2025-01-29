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
