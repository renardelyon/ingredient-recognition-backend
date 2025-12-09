package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ServerPort               string  `json:"server_port"`
	ImageServiceURL          string  `json:"image_service_url"`
	ServerAddress            string  `json:"server_address"`
	AWSRegion                string  `json:"aws_region"`
	AWSBucket                string  `json:"aws_bucket"`
	RekognitionProjectARN    string  `json:"rekognition_project_arn"`
	RekognitionModelVersion  string  `json:"rekognition_model_version"`
	RekognitionMinConfidence float32 `json:"rekognition_min_confidence"`
	DynamoDBTable            string  `json:"dynamodb_table"`
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
