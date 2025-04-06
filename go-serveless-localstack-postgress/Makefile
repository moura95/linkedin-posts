# Envs
BINARY_NAME=bootstrap
LAMBDA_HANDLER=main

# Build
build:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) ./cmd/api/main.go

# Clean
clean:
	rm -f $(BINARY_NAME)
	rm -rf ./bin
	rm -rf function.zip

# Start Containers local
start-local:
	docker compose up -d

# Stop Containers
down:
		docker compose down --volumes && docker volume prune -f

# Start LocalStack and deploy
deploy-local: build
	chmod +x ./scripts/deploy-local.sh
	./scripts/deploy-local.sh

# Start and configure LocalStack
setup-localstack:
	chmod +x ./scripts/setup-localstack.sh
	./scripts/setup-localstack.sh

start: setup-localstack start-local deploy-local
	@echo "Local environment fully configured and deployed!"

.PHONY: build clean start-local deploy-local setup-localstack invoke-local start