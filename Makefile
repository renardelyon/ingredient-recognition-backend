## Golang Stuff
GOCMD=go
GORUN=$(GOCMD) run

SERVICE=ingredient-recognition-backend

# Swagger API docs
SWAGGER_PORT=51234

init:
	$(GOCMD) mod init $(SERVICE)

tidy:
	$(GOCMD) mod tidy

build:
	$(GOCMD) build -o bin/ingredient-detector ./cmd

run:
	echo "for local development, please run: make run ENV=local"
	$(GORUN) cmd/main.go

	