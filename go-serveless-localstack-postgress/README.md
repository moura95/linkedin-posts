# Go Serverless com PostgreSQL e LocalStack

Este projeto demonstra como desenvolver uma aplicação serverless Go utilizando PostgreSQL como banco de dados e LocalStack para emular os serviços AWS localmente e terraform para subir infra.

## Visão Geral

Este projeto implementa uma API de tickets seguindo princípios de:
- Desenvolvimento Serverless (AWS Lambda e API Gateway)
- Testes de integração eficientes (com contêineres Docker)
- Interface de banco de dados estruturada (PostgreSQL)
- Testabilidade e facilidade de desenvolvimento local

## Estrutura do Projeto

```
go-serverless-postgres/
├── cmd/                    # Entradas da aplicação
│   └── api/                # Funções Lambda
│       └── main.go         # Ponto de entrada principal
├── internal/               # Código interno da aplicação
│   ├── handler/            # Handlers das funções Lambda
│   ├── repository/         # Acesso ao banco de dados
│   ├── model/              # Definições de tipos e modelos
│   └── db/                 # Clientes de banco de dados
├── pkg/                    # Pacotes compartilhados
├── migrations/             # Migrações de banco de dados
├── scripts/                # Scripts utilitários
├── docker-compose.yml      # Configuração de contêineres
├── template.yaml           # Template SAM
└── Makefile                # Comandos úteis
```

## Recursos Implementados

- **API de Tickets**: Criação, listagem, atualização e exclusão de tickets
- **LocalStack**: Emulação local de serviços AWS (Lambda, API Gateway)
- **Migrações**: Gerenciamento automático de schemas de banco de dados

## Pré-requisitos

- Go 1.22 ou superior
- Docker e Docker Compose
- AWS CLI

## Como Usar

### Configuração Inicial

1. Clone o repositório:
   ```
   git clone github/moura95/go-serverless-localstack-postgres.git
   cd go-serverless-localstack-postgres
   ```

### Desenvolvimento
1. Inicie o ambiente local,configure e compila usando terraform:
   ```
   make start-tf
   ```


2. Ou Inicia,Configura e Compila usando bash:
   ```
   make deploy-local
   ```


1. **Arquitetura Serverless**: Funções Lambda em vez de aplicação monolítica
2. **LocalStack**: Emulação local dos serviços AWS
3. **Implantação Simplificada**: Facilidade para testar localmente antes da implantação

## Por que esta Abordagem?

1. **Desenvolvimento Local Eficiente**: LocalStack permite testar localmente sem necessidade de conta AWS
2**Manutenção Simplificada**: Separação clara de responsabilidades (Repository, Handler, Models)
3**Baixo Custo de Operação**: Arquitetura serverless reduz custos em produção

## Referências

- [AWS Lambda com Go](https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html)
- [LocalStack Documentation](https://docs.localstack.cloud/overview/)
- [Go PostgreSQL Driver](https://github.com/lib/pq)