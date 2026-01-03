package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort               string  `mapstructure:"server_port"`
	ServerAddress            string  `mapstructure:"server_address"`
	AWSRegion                string  `mapstructure:"aws_region"`
	AWSBucket                string  `mapstructure:"aws_bucket"`
	RekognitionProjectARN    string  `mapstructure:"rekognition_project_arn"`
	RekognitionModelARN      string  `mapstructure:"rekognition_model_arn"`
	RekognitionModelVersion  string  `mapstructure:"rekognition_model_version"`
	RekognitionMinConfidence float32 `mapstructure:"rekognition_min_confidence"`
	JWTSecret                string  `mapstructure:"jwt_secret"`
	JWTExpiry                int     `mapstructure:"jwt_expiry_hours"`
	BedrockModelID           string  `mapstructure:"bedrock_model_id"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	// Set config file settings
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/ingredient-detector/")
	v.AddConfigPath("$HOME/.ingredient-detector")

	// Enable environment variable reading
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables to config keys
	v.BindEnv("server_port", "SERVER_PORT")
	v.BindEnv("server_address", "SERVER_ADDRESS")
	v.BindEnv("aws_region", "AWS_REGION")
	v.BindEnv("aws_bucket", "AWS_BUCKET")
	v.BindEnv("rekognition_project_arn", "REKOGNITION_PROJECT_ARN")
	v.BindEnv("rekognition_model_arn", "REKOGNITION_MODEL_ARN")
	v.BindEnv("rekognition_model_version", "REKOGNITION_MODEL_VERSION")
	v.BindEnv("rekognition_min_confidence", "REKOGNITION_MIN_CONFIDENCE")
	v.BindEnv("jwt_secret", "JWT_SECRET")
	v.BindEnv("jwt_expiry_hours", "JWT_EXPIRY_HOURS")
	v.BindEnv("bedrock_model_id", "BEDROCK_MODEL_ID")

	// Try to read config file (ignore error if not found - will use env vars)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Return error only if it's not a "config file not found" error
			return nil, err
		}
	}

	config := Config{}
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
