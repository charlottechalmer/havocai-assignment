package parser

import (
	"encoding/xml"
	"fmt"
	"havocai-assignment/models"
)

func ParseXML(input []byte) error {
	var patients models.Patients
	err := xml.Unmarshal(input, &patients)
	if err != nil {
		return err
	}
	fmt.Printf("unmarshalled: %+v\n", patients)
	return nil
}
