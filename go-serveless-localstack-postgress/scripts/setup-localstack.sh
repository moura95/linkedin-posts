#!/bin/bash
set -e
set -o pipefail

# Configurações silenciosas
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export AWS_DEFAULT_REGION=us-east-1
export AWS_ENDPOINT_URL=http://localhost:4566

# Função para log silencioso
log() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $*" >&2
}

# Configuração de função Lambda
setup_lambda_role() {
    log "Configurando IAM Role para Lambda..."
    aws --endpoint-url=http://localhost:4566 iam create-role \
        --role-name lambda-role \
        --assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}' \
        > /dev/null 2>&1 || true

    aws --endpoint-url=http://localhost:4566 iam attach-role-policy \
        --role-name lambda-role \
        --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole \
        > /dev/null 2>&1 || true
}

# Verificação de PostgreSQL
wait_for_postgres() {
    local max_attempts=10
    local attempt=1

    log "Verificando conexão PostgreSQL..."
    while [ $attempt -le $max_attempts ]; do
        if PGPASSWORD=postgres psql -h localhost -U postgres -d tickets_db -c "SELECT 1" > /dev/null 2>&1; then
            log "PostgreSQL está pronto"
            return 0
        fi

        log "Tentativa $attempt de $max_attempts..."
        sleep 3
        ((attempt++))
    done

    log "PostgreSQL não responde. Criando banco de dados..."
    PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE tickets_db;" > /dev/null 2>&1 || true
}

# Função principal
main() {
    log "Iniciando configuração do ambiente LocalStack..."

    wait_for_postgres
    setup_lambda_role

    log "Ambiente configurado com sucesso!"
}

# Executa a função principal
main