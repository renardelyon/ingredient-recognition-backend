package aws

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Client wraps the AWS S3 service
type S3Client struct {
	client *s3.Client
	bucket string
}

// NewS3Client creates a new S3 client
func NewS3Client(client *s3.Client, bucket string) *S3Client {
	return &S3Client{
		client: client,
		bucket: bucket,
	}
}

// UploadImage uploads an image to S3
func (sc *S3Client) UploadImage(ctx context.Context, key string, imageData []byte) (string, error) {
	input := &s3.PutObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageData),
	}

	_, err := sc.client.PutObject(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload image to S3: %w", err)
	}

	objectURL := fmt.Sprintf("s3://%s/%s", sc.bucket, key)
	return objectURL, nil
}

// DownloadImage downloads an image from S3
func (sc *S3Client) DownloadImage(ctx context.Context, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(key),
	}

	output, err := sc.client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download image from S3: %w", err)
	}
	defer output.Body.Close()

	imageData, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data from S3: %w", err)
	}

	return imageData, nil
}

// DeleteImage deletes an image from S3
func (sc *S3Client) DeleteImage(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(key),
	}

	_, err := sc.client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete image from S3: %w", err)
	}

	return nil
}

// ListImages lists all images in the bucket
func (sc *S3Client) ListImages(ctx context.Context) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(sc.bucket),
	}

	output, err := sc.client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list images from S3: %w", err)
	}

	var keys []string
	for _, obj := range output.Contents {
		keys = append(keys, aws.ToString(obj.Key))
	}

	return keys, nil
}
