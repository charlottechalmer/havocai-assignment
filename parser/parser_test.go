package parser

import (
	"encoding/json"
	"havocai-assignment/models"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcatTransformation(t *testing.T) {
	tests := []struct {
		name           string
		record         map[string]interface{}
		transformation models.Transformation
		expected       string
		expectedErr    bool
	}{
		{
			name: "valid concat of 2 strings with space separator",
			record: map[string]interface{}{
				"FirstName": "Charlotte",
				"LastName":  "Taylor",
			},
			transformation: models.Transformation{
				Type: "concat",
				Params: models.Params{
					Fields: []string{
						"FirstName",
						"LastName",
					},
					Extras: map[string]interface{}{
						"separator": " ",
					},
				},
			},
			expected:    "Charlotte Taylor",
			expectedErr: false,
		},
		{
			name: "valid concat of multiple fields of different types",
			record: map[string]interface{}{
				"BirthMonth": "July",
				"BirthDay":   6,
				"BirthYear":  1993,
			},
			transformation: models.Transformation{
				Type: "concat",
				Params: models.Params{
					Fields: []string{
						"BirthMonth",
						"BirthDay",
						"BirthYear",
					},
					Extras: map[string]interface{}{
						"separator": " ",
					},
				},
			},
			expected:    "July 6 1993",
			expectedErr: false,
		},
		{
			name: "valid concat of multiple strings with newline",
			record: map[string]interface{}{
				"Street":    "1234 Foo Ave",
				"CityState": "FooBar, WI",
				"ZipCode":   "12345",
			},
			transformation: models.Transformation{
				Type: "concat",
				Params: models.Params{
					Fields: []string{
						"Street",
						"CityState",
						"ZipCode",
					},
					Extras: map[string]interface{}{
						"separator": "\n",
					},
				},
			},
			expected: `1234 Foo Ave
FooBar, WI
12345`,
			expectedErr: false,
		},
		{
			name: "field not in record",
			record: map[string]interface{}{
				"FirstName": "Charlotte",
			},
			transformation: models.Transformation{
				Type: "concat",
				Params: models.Params{
					Fields: []string{
						"FirstName",
						"LastName",
					},
					Extras: map[string]interface{}{
						"separator": " ",
					},
				},
			},
			expected:    "Charlotte",
			expectedErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := concatTransformation(test.record, test.transformation)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestCalculateTransformation(t *testing.T) {
	tests := []struct {
		name           string
		record         map[string]interface{}
		transformation models.Transformation
		expected       interface{}
		expectedErr    bool
	}{
		{
			name: "Addition",
			record: map[string]interface{}{
				"Field1": 10,
				"Field2": 20,
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Field1",
						"Field2",
					},
					Extras: map[string]interface{}{
						"operation": "add",
					},
				},
			},
			expected:    30.0,
			expectedErr: false,
		},
		{
			name: "Subtraction",
			record: map[string]interface{}{
				"Field1": 50,
				"Field2": 20,
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Field1",
						"Field2",
					},
					Extras: map[string]interface{}{
						"operation": "subtract",
					},
				},
			},
			expected:    30.0,
			expectedErr: false,
		},
		{
			name: "Multiplication",
			record: map[string]interface{}{
				"Field1": 5,
				"Field2": 4,
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Field1",
						"Field2",
					},
					Extras: map[string]interface{}{
						"operation": "multiply",
					},
				},
			},
			expected:    20.0,
			expectedErr: false,
		},
		{
			name: "Division",
			record: map[string]interface{}{
				"Field1": 20,
				"Field2": 4,
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Field1",
						"Field2",
					},
					Extras: map[string]interface{}{
						"operation": "divide",
					},
				},
			},
			expected:    5.0,
			expectedErr: false,
		},
		{
			name: "Modulo",
			record: map[string]interface{}{
				"Field1": 20,
				"Field2": 3,
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Field1",
						"Field2",
					},
					Extras: map[string]interface{}{
						"operation": "modulo",
					},
				},
			},
			expected:    2.0,
			expectedErr: false,
		},
		{
			name: "Time difference in days",
			record: map[string]interface{}{
				"Start": "2025-01-05",
				"End":   "2025-01-29",
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Start",
						"End",
					},
					Extras: map[string]interface{}{
						"operation": "time_difference",
						"format":    "2006-01-02",
						"unit":      "days",
					},
				},
			},
			expected:    24.0,
			expectedErr: false,
		},
		{
			name: "Time difference in days, no format specified",
			record: map[string]interface{}{
				"Start": "2025-01-05T00:00:00Z",
				"End":   "2025-02-05T00:00:00Z",
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Start",
						"End",
					},
					Extras: map[string]interface{}{
						"operation": "time_difference",
						"unit":      "days",
					},
				},
			},
			expected:    31.0,
			expectedErr: false,
		},
		{
			name: "Time difference in seconds",
			record: map[string]interface{}{
				"Start": "2025-01-05T00:00:00",
				"End":   "2025-01-05T00:01:00",
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Start",
						"End",
					},
					Extras: map[string]interface{}{
						"operation": "time_difference",
						"format":    "2006-01-02T15:04:05",
						"unit":      "seconds",
					},
				},
			},
			expected:    60.0,
			expectedErr: false,
		},
		{
			name: "Time difference in weeks, rounded to int",
			record: map[string]interface{}{
				"Start": "2025-01-05",
				"End":   "2025-03-05",
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Start",
						"End",
					},
					Extras: map[string]interface{}{
						"operation":    "time_difference",
						"format":       "2006-01-02",
						"unit":         "weeks",
						"round_to_int": true,
					},
				},
			},
			expected:    8,
			expectedErr: false,
		},
		{
			name: "Time difference in months, with decimal precision",
			record: map[string]interface{}{
				"Start": "2025-01-05",
				"End":   "2025-03-05",
			},
			transformation: models.Transformation{
				Type: "calculate",
				Params: models.Params{
					Fields: []string{
						"Start",
						"End",
					},
					Extras: map[string]interface{}{
						"operation":         "time_difference",
						"format":            "2006-01-02",
						"unit":              "months",
						"decimal_precision": 3,
					},
				},
			},
			expected:    1.938,
			expectedErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := calculateTransformation(test.record, test.transformation)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}
}

func TestParseXML(t *testing.T) {
	tests := []struct {
		name          string
		inputFilePath string
		expected      []map[string]interface{}
		expectedErr   bool
	}{
		{
			name:          "valid single patient",
			inputFilePath: "../test/testdata/basicpatient/single_patient.xml",
			expected: []map[string]interface{}{
				{
					"ID":          12345,
					"DateOfBirth": "1993-07-06",
					"FirstName":   "Charlotte",
					"LastName":    "Taylor",
				},
			},
			expectedErr: false,
		},
		{
			name:          "valid multiple patient",
			inputFilePath: "../test/testdata/basicpatient/multiple_patients.xml",
			expected: []map[string]interface{}{
				{
					"ID":          12345,
					"DateOfBirth": "1993-07-06",
					"FirstName":   "Charlotte",
					"LastName":    "Taylor",
				},
				{
					"ID":          53425,
					"DateOfBirth": "1920-11-25",
					"FirstName":   "Jane",
					"LastName":    "Doe",
				},
			},
			expectedErr: false,
		},
		{
			name:          "invalid xml",
			inputFilePath: "../test/testdata/basicpatient/invalid_xml.xml",
			expected:      nil,
			expectedErr:   true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			xmlData, err := os.ReadFile(test.inputFilePath)
			require.NoError(t, err)

			actual, err := ParseXML(xmlData)

			if test.expectedErr {
				require.Error(t, err)
				require.Nil(t, actual)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expected, actual)
			}
		})
	}

}

func TestConvertToJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       []map[string]interface{}
		config      *models.Config
		expected    string
		expectedErr bool
	}{
		{
			name: "Basic mapping",
			input: []map[string]interface{}{
				{"old_field": "value1"},
			},
			config: &models.Config{
				RootName: "data",
				Mappings: map[string]string{
					"old_field": "new_field",
				},
			},
			expected: `{
				"data": [{
					"new_field": "value1"
				}]
			}`,
			expectedErr: false,
		},
		{
			name: "Concat transformation",
			input: []map[string]interface{}{
				{"first": "John", "last": "Doe"},
			},
			config: &models.Config{
				RootName: "users",
				Transformations: map[string]models.Transformation{
					"full_name": {
						Type: "concat",
						Params: models.Params{
							Fields: []string{"first", "last"},
							Extras: map[string]interface{}{"separator": " "},
						},
					},
				},
			},
			expected: `{
				"users": [{
					"full_name": "John Doe"
				}]
			}`,
			expectedErr: false,
		},
		{
			name: "Transformation with missing field should skip missing field",
			input: []map[string]interface{}{
				{"first": "John"},
			},
			config: &models.Config{
				RootName: "users",
				Transformations: map[string]models.Transformation{
					"full_name": {
						Type: "concat",
						Params: models.Params{
							Fields: []string{"first", "last"},
							Extras: map[string]interface{}{"separator": " "},
						},
					},
				},
			},
			expected: `{
				"users": [{
					"full_name": "John"
				}]
			}`,
			expectedErr: false,
		},
		{
			name: "Calculation transformation",
			input: []map[string]interface{}{
				{"a": 5, "b": 3},
			},
			config: &models.Config{
				RootName: "calculations",
				Transformations: map[string]models.Transformation{
					"sum": {
						Type: "calculate",
						Params: models.Params{
							Fields: []string{"a", "b"},
							Extras: map[string]interface{}{"operation": "add"},
						},
					},
				},
			},
			expected: `{
				"calculations": [{
					"sum": 8
				}]
			}`,
			expectedErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := ConvertToJSON(test.input, test.config)

			if test.expectedErr {
				require.Error(t, err)
				require.Nil(t, actual)
			} else {
				require.NoError(t, err)
				var unmarshalledActual map[string]interface{}
				var unmarshalledExpected map[string]interface{}
				require.NoError(t, json.Unmarshal(actual, &unmarshalledActual))
				require.NoError(t, json.Unmarshal([]byte(test.expected), &unmarshalledExpected))
				require.Equal(t, unmarshalledExpected, unmarshalledActual)
			}
		})
	}
}
