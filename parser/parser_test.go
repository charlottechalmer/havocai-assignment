package parser

import (
	"encoding/xml"
	"havocai-assignment/models"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTranslateName(t *testing.T) {
	tests := []struct {
		name      string
		firstName string
		lastName  string
		expected  string
	}{
		{
			name:      "valid first and lastname",
			firstName: "Charlotte",
			lastName:  "Taylor",
			expected:  "Charlotte Taylor",
		},
		{
			name:      "empty first name",
			firstName: "",
			lastName:  "Taylor",
			expected:  "Taylor",
		},
		{
			name:      "empty last name",
			firstName: "Charlotte",
			lastName:  "",
			expected:  "Charlotte",
		},
		{
			name:      "both empty",
			firstName: "",
			lastName:  "",
			expected:  "",
		},
		{
			name:      "spaces between names",
			firstName: "  Charlotte ",
			lastName:  " Taylor   ",
			expected:  "Charlotte Taylor",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := translateName(test.firstName, test.lastName)
			require.Equal(t, test.expected, actual)
		})
	}
}

func TestTranslateAge(t *testing.T) {
	tests := []struct {
		name        string
		dateOfBirth string
		expected    int
		expectedErr bool
	}{
		{
			name:        "valid date of birth, birthday not yet passed this year",
			dateOfBirth: "1993-07-06",
			expected:    time.Now().Year() - 1993 - 1,
			expectedErr: false,
		},
		{
			name:        "valid date of birth, birthday has passed this year",
			dateOfBirth: "1993-01-20",
			expected:    time.Now().Year() - 1993,
			expectedErr: false,
		},
		{
			name:        "valid date of birth, birthday is today",
			dateOfBirth: time.Now().Format("2006-01-02"),
			expected:    0,
			expectedErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := translateAge(test.dateOfBirth)
			if test.expectedErr {
				require.Error(t, err)
				require.Equal(t, 0, actual)
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
		expected      *models.XMLPatients
		expectedErr   bool
	}{
		{
			name:          "valid single patient",
			inputFilePath: "../test/testdata/single_patient.xml",
			expected: &models.XMLPatients{
				XMLName: xml.Name{Local: "Patients"},
				Patients: []models.XMLPatient{
					{
						ID:          12345,
						FirstName:   "Charlotte",
						LastName:    "Taylor",
						DateOfBirth: "1993-07-06",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:          "valid multiple patient",
			inputFilePath: "../test/testdata/multiple_patients.xml",
			expected: &models.XMLPatients{
				XMLName: xml.Name{Local: "Patients"},
				Patients: []models.XMLPatient{
					{
						ID:          12345,
						FirstName:   "Charlotte",
						LastName:    "Taylor",
						DateOfBirth: "1993-07-06",
					},
					{
						ID:          53425,
						FirstName:   "Jane",
						LastName:    "Doe",
						DateOfBirth: "1920-11-25",
					},
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
		{
			name:          "empty",
			inputFilePath: "../test/testdata/empty.xml",
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
		name                       string
		inputXMLFilePath           string
		expectedJSONOutputFilePath string
		expectedErr                bool
	}{
		{
			name:                       "provided input and output",
			inputXMLFilePath:           "../test/testdata/input.xml",
			expectedJSONOutputFilePath: "../test/testdata/output.json",
			expectedErr:                false,
		},
		{
			name:                       "valid single patient",
			inputXMLFilePath:           "../test/testdata/single_patient.xml",
			expectedJSONOutputFilePath: "../test/testdata/single_patient.json",
			expectedErr:                false,
		},
		{
			name:                       "valid multiple patients",
			inputXMLFilePath:           "../test/testdata/multiple_patients.xml",
			expectedJSONOutputFilePath: "../test/testdata/multiple_patients.json",
			expectedErr:                false,
		},
		{
			name:                       "invalid date",
			inputXMLFilePath:           "../test/testdata/invalid_date.xml",
			expectedJSONOutputFilePath: "../test/testdata/error.json",
			expectedErr:                true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			xmlData, err := os.ReadFile(test.inputXMLFilePath)
			require.NoError(t, err)

			xmlPatients, err := ParseXML(xmlData)
			require.NoError(t, err)

			actual, err := ConvertToJSON(xmlPatients)
			if test.expectedErr {
				require.Error(t, err)
				require.Nil(t, actual)
			} else {
				require.NoError(t, err)
				expectedJSON, err := os.ReadFile(test.expectedJSONOutputFilePath)
				require.NoError(t, err)
				require.JSONEq(t, string(expectedJSON), string(actual))
			}
		})
	}
}
