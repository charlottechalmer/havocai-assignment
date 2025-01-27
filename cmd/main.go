package main

import (
	"fmt"
	"havocai-assignment/parser"
	"os"
)

func main() {
	input, err := os.ReadFile("test/testdata/input.xml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input file: %+v\n", err)
	}

	xmlPatients, err := parser.ParseXML(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing XML: %+v\n", err)
	}

	fmt.Printf("xml patients: %+v\n", xmlPatients)

	jsonPatients, err := parser.ConvertToJSON(xmlPatients)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting to JSON: %+v\n", err)
	}

	fmt.Printf("patients: %+v\n", string(jsonPatients))
}
