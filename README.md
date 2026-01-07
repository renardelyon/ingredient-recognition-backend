# Ingredient Recognition Backend

## Overview
The Ingredient Recognition Backend is a Go-based REST API that allows users to detect ingredients from images using AWS Rekognition, get AI-powered recipe recommendations using AWS Bedrock (Claude), and save their favorite recipes. The application features user authentication with JWT, structured logging, and DynamoDB for data persistence.

## Features
- **Ingredient Detection**: Upload images to detect ingredients using AWS Rekognition (with custom labels support)
- **Recipe Recommendations**: Get AI-powered recipe suggestions based on detected ingredients using AWS Bedrock (Claude)
- **Recipe Management**: Save, retrieve, and delete favorite recipes
- **User Authentication**: Secure JWT-based authentication system
- **Structured Logging**: Comprehensive logging with request tracking
- **AWS Integration**: S3 for image storage, DynamoDB for data persistence, Rekognition for image analysis, Bedrock for AI recommendations

## Getting Started

### Prerequisites
- Go 1.21 or later
- AWS Account with:
  - S3 bucket for image storage
  - DynamoDB tables: `Users` and `SavedRecipes`
  - Rekognition access (optionally with custom labels)
  - Bedrock access with Claude model
- AWS credentials configured

### DynamoDB Table Setup

#### Users Table
- **Partition Key**: `id` (String)
- **Global Secondary Index**: `EmailIndex`
  - Partition Key: `email` (String)

#### SavedRecipes Table
- **Partition Key**: `id` (String)
- **Sort Key**: `created_at` (String)
- **Global Secondary Index**: `UserIdIndex`
  - Partition Key: `user_id` (String)

### Configuration
Create a `config.json` file in the root directory:
```json
{
  "aws_region": "us-east-1",
  "aws_bucket": "your-s3-bucket-name",
  "bedrock_model_id": "anthropic.claude-3-5-sonnet-20240620-v1:0",
  "rekognition_project_arn": "arn:aws:rekognition:...",
  "rekognition_model_arn": "arn:aws:rekognition:...",
  "rekognition_model_version": "1",
  "rekognition_min_confidence": 70,
  "jwt_secret": "your-secure-secret-key",
  "jwt_expiry": 24,
  "server_address": ":8080"
}
```

### Installation
1. Clone the repository:
   ```
   git clone <repository-url>
   cd ingredient-detector
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Installation
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd ingredient-recognition-backend
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up your AWS credentials:
   ```bash
   export AWS_ACCESS_KEY_ID=your_access_key
   export AWS_SECRET_ACCESS_KEY=your_secret_key
   # Or use AWS CLI: aws configure
   ```

### Running the Application

Using Make:
```bash
make run
```

Or directly with Go:
```bash
go run cmd/main.go
```

Using Docker:
```bash
docker build -t ingredient-recognition-backend .
docker run -p 8080:8080 ingredient-recognition-backend
```

The server will start on `http://localhost:8080`

### Logs
Application logs are stored in `logs/app.log` with structured JSON format.

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.