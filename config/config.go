package config

import (
	"encoding/json"
	"havocai-assignment/models"
	"os"
)

func LoadFile(filePath string) (*models.Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config *models.Config
	err = json.Unmarshal(data, &config)
	return config, err
}
