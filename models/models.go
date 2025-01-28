package models

type Config struct {
	Mappings        map[string]string         `json:"mappings"`
	Transformations map[string]Transformation `json:"transformations"`
}

type Transformation struct {
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

// type XMLPatients struct {
// 	XMLName  xml.Name     `xml:"Patients"`
// 	Patients []XMLPatient `xml:"Patient"`
// }

// type XMLPatient struct {
// 	ID          int    `xml:"ID,attr"`
// 	FirstName   string `xml:"FirstName"`
// 	LastName    string `xml:"LastName"`
// 	DateOfBirth string `xml:"DateOfBirth"`
// }

// type JSONPatients struct {
// 	Patients []JSONPatient `json:"patients"`
// }

// type JSONPatient struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// 	Age  int    `json:"age"`
// }
