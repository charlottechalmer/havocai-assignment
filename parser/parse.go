package parser

import (
	"encoding/xml"
	"fmt"
	"havocai-assignment/models"
)

func ParseXML(input []byte) (*models.XMLPatients, error) {
	var patients *models.XMLPatients
	err := xml.Unmarshal(input, &patients)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func ConvertToJSON(input *models.XMLPatients) (models.JSONPatients, error) {
	var jsonPatients models.JSONPatients
	for _, patient := range input.Patients {
		name := fmt.Sprintf("%v %v", patient.FirstName, patient.LastName)
		jsonPatients.Patients = append(jsonPatients.Patients, models.JSONPatient{
			ID:   patient.ID,
			Name: name,
			Age:  3,
		})
	}
	return jsonPatients, nil
}
