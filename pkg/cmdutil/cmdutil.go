package cmdutil

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func FatalError(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%v: %+v\n", msg, err)
	os.Exit(1)
}

func ValidateFlags() (string, string) {
	xmlFilePath := flag.String("xml", "", "path to XML input file")
	configFilePath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *xmlFilePath == "" || *configFilePath == "" {
		fmt.Fprintf(os.Stderr, "Both -xml and -config flags are required\n")
		flag.Usage()
		os.Exit(1)
	}

	absXMLPath, err := filepath.Abs(filepath.Clean(*xmlFilePath))
	if err != nil {
		FatalError("error resolving xml input path: %+v\n", err)
	}

	absConfigPath, err := filepath.Abs(filepath.Clean(*configFilePath))
	if err != nil {
		FatalError("error resolving config file path: %+v\n", err)
	}

	return absXMLPath, absConfigPath
}
