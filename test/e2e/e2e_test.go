package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

/*
 TO DO:
 - add additional tests using different configs and input

 TESTS TO ADD:
 - Test different config inputs
   - changes to input data
       - fields are added
           - new 1:1 mappings
           - new transformations
               - bool condition; if DateOfDeath exists --> deceased:true
               - calculation; "last visit" --> x months ago
               - counting: nested elements --> "number_of_allergies"
       - fields are renamed
       - nested structures are added (address)
       - format of data might change (DOB -- DD-MM-YYY)
       - root element changes
   - changes to output requirements
       - field names change
       - output becomes more nested(e.g contact info, medical info)
       - changing format -- name --> Last, First
       - field types change -- age --> string
*/

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name             string
		configPath       string
		inputXMLPath     string
		expectedJSONPath string
	}{
		{
			name:             "provided input and output",
			configPath:       "../testdata/basicpatient/config.json",
			inputXMLPath:     "../testdata/provided/input.xml",
			expectedJSONPath: "../testdata/provided/output.json",
		},
		{
			name:             "valid single patient",
			configPath:       "../testdata/basicpatient/config.json",
			inputXMLPath:     "../testdata/basicpatient/single_patient.xml",
			expectedJSONPath: "../testdata/basicpatient/single_patient.json",
		},
		{
			name:             "valid multiple patients",
			configPath:       "../testdata/basicpatient/config.json",
			inputXMLPath:     "../testdata/basicpatient/multiple_patients.xml",
			expectedJSONPath: "../testdata/basicpatient/multiple_patients.json",
		},
		{
			name:             "new fields are added - mappings",
			configPath:       "../testdata/inputchanges/new_fields_config.json",
			inputXMLPath:     "../testdata/inputchanges/new_fields.xml",
			expectedJSONPath: "../testdata/inputchanges/new_fields.json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpOutput, err := os.CreateTemp("", "output_tmp.json")
			require.NoError(t, err)

			defer os.Remove(tmpOutput.Name())

			cmd := exec.Command("go", "run", "../../cmd/main.go", "-xml", test.inputXMLPath, "-config", test.configPath, "-output", tmpOutput.Name())

			var stderr bytes.Buffer
			cmd.Stderr = &stderr

			// var stdout bytes.Buffer
			// cmd.Stdout = &stdout

			err = cmd.Run()
			require.NoError(t, err, "error running main.go: %v", stderr.String())

			// fmt.Println("=== Program Output ===")
			// fmt.Println(stdout.String())

			actualOutput, err := os.ReadFile(tmpOutput.Name())
			require.NoError(t, err)

			var actualJSON map[string]interface{}
			err = json.Unmarshal(actualOutput, &actualJSON)
			require.NoError(t, err)

			expectedOutput, err := os.ReadFile(test.expectedJSONPath)
			require.NoError(t, err)
			var expectedJSON map[string]interface{}
			err = json.Unmarshal(expectedOutput, &expectedJSON)
			require.NoError(t, err)

			require.Equal(t, expectedJSON, actualJSON)
		})
	}
}
