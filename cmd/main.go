package main

import (
	"fmt"
	"havocai-assignment/config"
	"havocai-assignment/parser"
	"os"
)

func main() {
	config, err := config.LoadFile("config/config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config file: %+v\n", err)
	}
	fmt.Printf("config: %+v\n", config)

	// TODO: change to be passed in via cmdln
	input, err := os.ReadFile("test/testdata/input.xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input file: %+v\n", err)
	}

	xmlPatients, err := parser.ParseXML(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing XML: %+v\n", err)
	}

	fmt.Printf("xml patients: %+v\n", xmlPatients)

	jsonPatients, err := parser.ConvertToJSON(xmlPatients, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting to JSON: %+v\n", err)
	}

	// TODO: write to file
	fmt.Printf("%+v\n", string(jsonPatients))
}
