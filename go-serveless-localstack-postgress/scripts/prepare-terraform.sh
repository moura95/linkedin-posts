#!/bin/bash
set -e

echo "Preparando aplicação para implantação com Terraform..."

# Compilar a aplicação
GOOS=linux GOARCH=amd64 go build -o bootstrap ./cmd/api/main.go

# Criar arquivo ZIP para o Lambda
zip -j function.zip bootstrap

echo "Aplicação pronta para implantação com Terraform!"