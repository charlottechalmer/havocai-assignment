package e2e

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEndToEnd(t *testing.T) {
	tests := []struct {
		name             string
		configPath       string
		inputXMLPath     string
		expectedJSONPath string
	}{
		{
			name:             "provided input and output",
			configPath:       "../../config/config.json",
			inputXMLPath:     "../testdata/input.xml",
			expectedJSONPath: "../testdata/output.json",
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
