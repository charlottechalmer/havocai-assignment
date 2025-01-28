package parser

import (
	"encoding/json"
	"fmt"
	"havocai-assignment/models"
	"strings"
	"time"
)

func ParseXML(input []byte) ([]map[string]string, error) {
	//use xml.Decoder to create a decoder from a Reader, and then call decoder.Token() in a loop to iterate through all the elements of the xml document. You will need a type switch on the token to do anything useful (e.g. on StartToken, EndToken,etc.).
}

// func ParseXML(input []byte) (*models.XMLPatients, error) {
// 	var patients *models.XMLPatients
// 	err := xml.Unmarshal(input, &patients)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return patients, nil
// }

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

func ConvertToJSON(input *models.XMLPatients) ([]byte, error) {
	jsonPatients := &models.JSONPatients{}
	for _, patient := range input.Patients {
		name := translateName(patient.FirstName, patient.LastName)
		age, err := translateAge(patient.DateOfBirth)
		if err != nil {
			return nil, err
		}

		jsonPatients.Patients = append(jsonPatients.Patients, models.JSONPatient{
			ID:   patient.ID,
			Name: name,
			Age:  age,
		})
	}

	jsonOutput, err := json.MarshalIndent(jsonPatients, "", "  ")
	if err != nil {
		return nil, err
	}

	return jsonOutput, nil
}
