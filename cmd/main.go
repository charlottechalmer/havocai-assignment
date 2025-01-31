package main

import (
	"fmt"
	"havocai-assignment/config"
	"havocai-assignment/parser"
	"havocai-assignment/pkg/cmdutil"
	"havocai-assignment/pkg/fileutil"
	"os"
)

func main() {
	xmlPath, configPath, outputPath := cmdutil.ValidateFlags()

	config, err := config.LoadFile(configPath)
	if err != nil {
		cmdutil.FatalError("error loading config file: %+v\n", err)
	}

	input, err := os.ReadFile(xmlPath)
	if err != nil {
		cmdutil.FatalError("error reading xml input file: %+v\n", err)
	}

	xmlPatients, err := parser.ParseXML(input)
	if err != nil {
		cmdutil.FatalError("error parsing XML: %+v\n", err)
	}

	jsonPatients, err := parser.ConvertToJSON(xmlPatients, config)
	if err != nil {
		cmdutil.FatalError("error converting to JSON: %+v\n", err)
	}

	if outputPath == "" {
		outputPath, err = fileutil.GetOutputPath()
		if err != nil {
			cmdutil.FatalError("error determining output path: %+v\n", err)
		}
	}

	err = fileutil.WriteToFile(outputPath, jsonPatients)
	if err != nil {
		cmdutil.FatalError("error writing output to file: %+v\n", err)
	}

	fmt.Printf("Successfully converted XML data to JSON. Output written to: %v\n", outputPath)
}
