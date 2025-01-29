package parser

import (
	"havocai-assignment/models"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
 TESTS TO ADD:
 - TestCalculateTransformation
 - TestConvertToJSON
 - Test different config inputs
   - changes to input data
   - changes to output requirements
   - different transformations
     - deceased: bool (if DateOfDeath exists)
 - calculations
    - counting
    - "last visit" -- return (x months || x days ago)
*/

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
			expected:    "",
			expectedErr: true,
		},
		{
			name: "empty fields",
			record: map[string]interface{}{
				"Field1": "foobar",
			},
			transformation: models.Transformation{
				Type: "concat",
				Params: models.Params{
					Fields: []string{""},
					Extras: map[string]interface{}{
						"separator": " ",
					},
				},
			},
			expected:    "",
			expectedErr: true,
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
	// tests := []struct {
	// }{}
}

// func TestTranslateAge(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		dateOfBirth string
// 		expected    int
// 		expectedErr bool
// 	}{
// 		{
// 			name:        "valid date of birth, birthday not yet passed this year",
// 			dateOfBirth: "1993-07-06",
// 			expected:    time.Now().Year() - 1993 - 1,
// 			expectedErr: false,
// 		},
// 		{
// 			name:        "valid date of birth, birthday has passed this year",
// 			dateOfBirth: "1993-01-20",
// 			expected:    time.Now().Year() - 1993,
// 			expectedErr: false,
// 		},
// 		{
// 			name:        "valid date of birth, birthday is today",
// 			dateOfBirth: time.Now().Format("2006-01-02"),
// 			expected:    0,
// 			expectedErr: false,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			actual, err := translateAge(test.dateOfBirth)
// 			if test.expectedErr {
// 				require.Error(t, err)
// 				require.Equal(t, 0, actual)
// 			} else {
// 				require.NoError(t, err)
// 				require.Equal(t, test.expected, actual)
// 			}
// 		})
// 	}
// }

func TestParseXML(t *testing.T) {
	tests := []struct {
		name          string
		inputFilePath string
		expected      []map[string]interface{}
		expectedErr   bool
	}{
		{
			name:          "valid single patient",
			inputFilePath: "../test/testdata/single_patient.xml",
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
			inputFilePath: "../test/testdata/multiple_patients.xml",
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
			inputFilePath: "../test/testdata/invalid_xml.xml",
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

// func TestConvertToJSON(t *testing.T) {
// 	// TODO: add config
// 	tests := []struct {
// 		name                       string
// 		inputXMLFilePath           string
// 		expectedJSONOutputFilePath string
// 		expectedErr                bool
// 	}{
// 		{
// 			name:                       "provided input and output",
// 			inputXMLFilePath:           "../test/testdata/input.xml",
// 			expectedJSONOutputFilePath: "../test/testdata/output.json",
// 			expectedErr:                false,
// 		},
// 		{
// 			name:                       "valid single patient",
// 			inputXMLFilePath:           "../test/testdata/single_patient.xml",
// 			expectedJSONOutputFilePath: "../test/testdata/single_patient.json",
// 			expectedErr:                false,
// 		},
// 		{
// 			name:                       "valid multiple patients",
// 			inputXMLFilePath:           "../test/testdata/multiple_patients.xml",
// 			expectedJSONOutputFilePath: "../test/testdata/multiple_patients.json",
// 			expectedErr:                false,
// 		},
// 		{
// 			name:                       "invalid date",
// 			inputXMLFilePath:           "../test/testdata/invalid_date.xml",
// 			expectedJSONOutputFilePath: "../test/testdata/error.json",
// 			expectedErr:                true,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			xmlData, err := os.ReadFile(test.inputXMLFilePath)
// 			require.NoError(t, err)

// 			xmlPatients, err := ParseXML(xmlData)
// 			require.NoError(t, err)

// 			actual, err := ConvertToJSON(xmlPatients)
// 			if test.expectedErr {
// 				require.Error(t, err)
// 				require.Nil(t, actual)
// 			} else {
// 				require.NoError(t, err)
// 				expectedJSON, err := os.ReadFile(test.expectedJSONOutputFilePath)
// 				require.NoError(t, err)
// 				require.JSONEq(t, string(expectedJSON), string(actual))
// 			}
// 		})
// 	}
// }
