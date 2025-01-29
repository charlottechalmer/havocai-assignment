package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetOutputPath() (string, error) {
	outputDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error getting user home directory: %w", err)
	}

	outputDir = filepath.Join(outputDir, "Documents", "xml-to-json-output")

	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating output directory: %w", err)
	}

	outputFileName := fmt.Sprintf("output_%d.json", time.Now().Unix())
	outputFilePath := filepath.Join(outputDir, outputFileName)

	return outputFilePath, nil
}

func WriteToFile(filepath string, data []byte) error {
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}

	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("error writing data to file: %w", err)
	}
	return nil
}
