# Docker Setup Guide

This guide explains how to build and run the Ingredient Recognition Backend using Docker.

## Prerequisites

- Docker (version 20.10+)
- Docker Compose (version 1.29+)
- AWS credentials configured (for production)

## Quick Start

### Local Development

**Build the Docker image:**
```bash
docker build -t ingredient-recognition-backend:latest .
```

This will start:
- The Go backend application on `http://localhost:8080`
- Local DynamoDB on `http://localhost:8000`
- Automatically create the Users table

## Production Deployment

### Build for Production

```bash
docker build \
  --build-arg ENV=production \
  -t ingredient-recognition-backend:v1.0.0 \
  -t ingredient-recognition-backend:latest \
  .
```

### Push to Container Registry

#### Docker Hub
```bash
docker tag ingredient-recognition-backend:latest your-username/ingredient-recognition-backend:latest
docker push your-username/ingredient-recognition-backend:latest
```

### Run in Production

```bash
docker run \
  --name ingredient-backend \
  -p 8080:8080 \
  -v ./config.json:/root/app/config.json:ro \
  -v ~/.aws:/root/.aws:ro \
  -e AWS_REGION=us-east-1 \
  --restart always \
  ingredient-recognition-backend:latest
```

## Health Check

The container includes a built-in health check that:
- Checks the `/health` endpoint every 30 seconds
- Starts after 40 seconds (grace period)
- Times out after 10 seconds
- Fails after 3 consecutive failures

Check health status:
```bash
docker ps --filter "name=ingredient-recognition-backend"
```

## Debugging

### View Logs
```bash
docker logs ingredient-recognition-backend
docker logs -f ingredient-recognition-backend  # Follow logs
docker logs --tail 100 ingredient-recognition-backend  # Last 100 lines
```

### Execute Commands in Container
```bash
docker exec -it ingredient-recognition-backend /bin/sh
```

### Check Container Stats
```bash
docker stats ingredient-recognition-backend
```

### Inspect Container
```bash
docker inspect ingredient-recognition-backend
```

## Docker Compose Commands

### Start services
```bash
docker-compose up
docker-compose up -d  # Run in background
```

### Stop services
```bash
docker-compose down
docker-compose down -v  # Also remove volumes
```

### View logs
```bash
docker-compose logs
docker-compose logs -f app  # Follow app logs
docker-compose logs -f dynamodb  # Follow DynamoDB logs
```

### Rebuild images
```bash
docker-compose build
docker-compose up --build
```

### Remove containers and volumes
```bash
docker-compose down -v
```

## Best Practices

1. **Use specific image tags** in production, not `latest`
2. **Scan images for vulnerabilities** before deployment
3. **Keep base image updated** for security patches
4. **Use environment variables** for configuration, not hardcoded values
5. **Monitor container health** with health checks
6. **Implement proper logging** for debugging
7. **Set resource limits** to prevent runaway containers
8. **Use volumes** for persistent data
9. **Regular backups** of configuration and logs
10. **Document environment-specific settings**

## Cleanup

### Remove stopped containers
```bash
docker container prune
```

### Remove unused images
```bash
docker image prune
```

### Remove unused volumes
```bash
docker volume prune
```

### Complete cleanup (careful!)
```bash
docker system prune -a --volumes
```
