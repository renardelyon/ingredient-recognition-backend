package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWSClient holds references to all AWS service clients
type AWSClient struct {
	Rekognition *RekognitionClient
	S3          *S3Client
	DynamoDB    *dynamodb.Client
}

// NewAWSClient initializes all AWS service clients
func NewAWSClient(ctx context.Context, region string, s3Bucket string) (*AWSClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Rekognition client
	rekognitionClient := rekognition.NewFromConfig(cfg)
	rekognitionSvc := NewRekognitionClient(rekognitionClient)

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)
	s3Svc := NewS3Client(s3Client, s3Bucket)

	// Create DynamoDB client
	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	return &AWSClient{
		Rekognition: rekognitionSvc,
		S3:          s3Svc,
		DynamoDB:    dynamoDBClient,
	}, nil
}
