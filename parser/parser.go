package parser

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

func ParseXML(input []byte) ([]map[string]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(input))
	var results []map[string]string

	// to track current data within loop
	currentData := make(map[string]string)
	currentElementName := ""
	//track nesting to determine when to exit a grouping
	level := 0

	for {
		token, err := decoder.Token()
		if err != nil {
			if err.Error() == "EOF" {
				// end of XML returns EOF, break from loop
				break
			}
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			// when we encounter a start element, want to increment level (noting we are within an element) and track the current element name
			level++
			currentElementName = t.Name.Local
		case xml.EndElement:
			// when we encounter an end element, want to decrease level (exiting an element) and then append data to results and reset currentData and currentElementName
			level--

			// if level == 1, we are at the end of a grouping
			if level == 1 {
				results = append(results, currentData)
				currentData = make(map[string]string)
			}

			currentElementName = ""
		case xml.CharData:
			// when we encounter CharData, store the character data at the current element
			content := strings.TrimSpace(string(t))
			if currentElementName != "" && content != "" {
				currentData[currentElementName] = content
			}
		}
	}
	return results, nil
}

func translateName(firstName string, lastName string) string {
	sanitizedFirstName := strings.TrimSpace(firstName)
	sanitizedLastName := strings.TrimSpace(lastName)
	return strings.TrimSpace(fmt.Sprintf("%v %v", sanitizedFirstName, sanitizedLastName))
}

func translateAge(dateOfBirth string) (int, error) {
	birthDate, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return 0, err
	}

	curr := time.Now()
	age := curr.Year() - birthDate.Year()
	if curr.YearDay() < birthDate.YearDay() {
		age--
	}

	return age, nil
}

func ConvertToJSON(input []map[string]string) ([]byte, error) {
	jsonOutput, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonOutput, nil
}

// func ConvertToJSON(input *models.XMLPatients) ([]byte, error) {
// 	jsonPatients := &models.JSONPatients{}
// 	for _, patient := range input.Patients {
// 		name := translateName(patient.FirstName, patient.LastName)
// 		age, err := translateAge(patient.DateOfBirth)
// 		if err != nil {
// 			return nil, err
// 		}

// 		jsonPatients.Patients = append(jsonPatients.Patients, models.JSONPatient{
// 			ID:   patient.ID,
// 			Name: name,
// 			Age:  age,
// 		})
// 	}

// 	jsonOutput, err := json.MarshalIndent(jsonPatients, "", "  ")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jsonOutput, nil
// }
