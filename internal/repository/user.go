package repository

import (
	"context"
	"fmt"
	"ingredient-recognition-backend/internal/domain"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// UserRepository is a DynamoDB implementation of UserRepository
type UserRepository struct {
	client    *dynamodb.Client
	tableName string
}

// UserRepository creates a new DynamoDB user repository
func NewUserRepository(client *dynamodb.Client, tableName string) *UserRepository {
	return &UserRepository{
		client:    client,
		tableName: tableName,
	}
}

// Create creates a new user in DynamoDB
func (r UserRepository) Create(ctx context.Context, user *domain.User) error {
	// Check if user already exists
	_, err := r.GetByEmail(ctx, user.Email)
	if err == nil && err != domain.ErrUserNotFound {
		return domain.ErrUserAlreadyExists
	}
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}

	// Marshal user to DynamoDB item
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Add email as a GSI key for faster lookups
	item["email"] = &types.AttributeValueMemberS{Value: user.Email}

	// Put item in DynamoDB
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByEmail retrieves a user by email from DynamoDB
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Query using GSI on email
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("EmailIndex"),
		KeyConditionExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query user by email: %w", err)
	}

	if result.Count == 0 {
		return nil, domain.ErrUserNotFound
	}

	var user domain.User
	err = attributevalue.UnmarshalMap(result.Items[0], &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// GetByID retrieves a user by ID from DynamoDB
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if result.Item == nil {
		return nil, domain.ErrUserNotFound
	}

	var user domain.User
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal user: %w", err)
	}

	return &user, nil
}

// Update updates an existing user in DynamoDB
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	// Check if user exists
	_, err := r.GetByID(ctx, user.Id)
	if err != nil {
		return err
	}

	// Marshal user to DynamoDB item
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	// Put item in DynamoDB (this will update if exists)
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user from DynamoDB
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	// Check if user exists
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete item from DynamoDB
	_, err = r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
