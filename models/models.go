package models

import "encoding/xml"

type Patients struct {
	XMLName  xml.Name  `xml:Patients`
	Patients []Patient `xml:"Patient" json:"patients"`
}

type Patient struct {
	ID          int    `xml:"ID,attr" json:"id"`
	FirstName   string `xml:"FirstName" json:"name"`
	LastName    string `xml:"LastName"`
	DateOfBirth string `xml:"DateOfBirth"`
	// Age  int    `xml:"DateOfBirth" json:"age"`
}
