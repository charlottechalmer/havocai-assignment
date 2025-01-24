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

	err = parser.ParseXML(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error deserializing: %+v\n", err)
	}
}
