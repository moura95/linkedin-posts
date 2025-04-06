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

# Comandos Terraform
tf-prepare:
	./scripts/prepare-terraform.sh

tf-init:
	cd terraform && terraform init

tf-plan:
	cd terraform && terraform plan

tf-apply:
	cd terraform && terraform apply -auto-approve

tf-destroy:
	cd terraform && terraform destroy -auto-approve

deploy-tf: build tf-apply
	@echo "Aplicação implantada com Terraform!"

start-tf: start-local tf-prepare build tf-apply
	@echo "Ambiente local com Terraform totalmente configurado e implantado!"

.PHONY: build clean start-local deploy-local setup-localstack invoke-local start down tf-init tf-plan tf-apply tf-destroy tf-prepare