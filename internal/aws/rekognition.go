package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
)

// RekognitionClient wraps the AWS Rekognition service
type RekognitionClient struct {
	client *rekognition.Client
}

// NewRekognitionClient creates a new Rekognition client
func NewRekognitionClient(client *rekognition.Client) *RekognitionClient {
	return &RekognitionClient{client: client}
}

// DetectLabels detects labels (objects, scenes, concepts) in an image
func (rc *RekognitionClient) DetectLabels(ctx context.Context, imageData []byte) ([]string, error) {
	input := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: imageData,
		},
		MaxLabels:     aws.Int32(100),
		MinConfidence: aws.Float32(50),
	}

	output, err := rc.client.DetectLabels(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to detect labels: %w", err)
	}

	var labels []string
	for _, label := range output.Labels {
		labels = append(labels, aws.ToString(label.Name))
	}

	return labels, nil
}

// DetectLabelsFromS3 detects labels in an image stored in S3
func (rc *RekognitionClient) DetectLabelsFromS3(ctx context.Context, bucket, key string) ([]string, error) {
	input := &rekognition.DetectLabelsInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
		MaxLabels:     aws.Int32(100),
		MinConfidence: aws.Float32(50),
	}

	output, err := rc.client.DetectLabels(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to detect labels from S3: %w", err)
	}

	var labels []string
	for _, label := range output.Labels {
		labels = append(labels, aws.ToString(label.Name))
	}

	return labels, nil
}

// DetectCustomLabels detects custom labels in an image using a trained model
func (rc *RekognitionClient) DetectCustomLabels(ctx context.Context, imageData []byte, projectARN, modelVersion string, minConfidence float32) (map[string]float32, error) {
	input := &rekognition.DetectCustomLabelsInput{
		Image: &types.Image{
			Bytes: imageData,
		},
		ProjectVersionArn: aws.String(projectARN),
		MinConfidence:     aws.Float32(minConfidence),
	}

	output, err := rc.client.DetectCustomLabels(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to detect custom labels: %w", err)
	}

	// Map label names to confidence scores
	labels := make(map[string]float32)
	for _, label := range output.CustomLabels {
		labels[aws.ToString(label.Name)] = aws.ToFloat32(label.Confidence)
	}

	return labels, nil
}

// DetectCustomLabelsFromS3 detects custom labels in an S3 image
func (rc *RekognitionClient) DetectCustomLabelsFromS3(ctx context.Context, bucket, key, projectARN, modelVersion string, minConfidence float32) (map[string]float32, error) {
	input := &rekognition.DetectCustomLabelsInput{
		Image: &types.Image{
			S3Object: &types.S3Object{
				Bucket: aws.String(bucket),
				Name:   aws.String(key),
			},
		},
		ProjectVersionArn: aws.String(projectARN),
		MinConfidence:     aws.Float32(minConfidence),
	}

	output, err := rc.client.DetectCustomLabels(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to detect custom labels from S3: %w", err)
	}

	// Map label names to confidence scores
	labels := make(map[string]float32)
	for _, label := range output.CustomLabels {
		labels[aws.ToString(label.Name)] = aws.ToFloat32(label.Confidence)
	}

	return labels, nil
}
