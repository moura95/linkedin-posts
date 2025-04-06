# Testando Aplicações Go com TestContainers

Este repositório demonstra como usar [TestContainers](https://testcontainers.com/) para testes de integração em aplicações Go. O TestContainers permite criar e gerenciar contêineres Docker para seus testes, facilitando a implementação de testes de integração confiáveis para componentes que dependem de serviços externos, como bancos de dados.

## Visão Geral do Projeto

Este projeto de exemplo de API de Tickets demonstra:

- Uma API REST para gerenciamento de tickets e categorias
- Banco de dados PostgreSQL para persistência de dados
- Testes abrangentes com TestContainers

## Por que usar TestContainers?

- **Ambientes de teste isolados**: Cada teste executa em seu próprio contêiner limpo e isolado
- **Testes realistas**: Teste contra instâncias reais de banco de dados, não mocks
- **Resultados reproduzíveis**: Os testes são executados consistentemente em qualquer ambiente
- **Execução paralela**: Múltiplos testes podem ser executados simultaneamente sem interferência
- **Sem dependências externas**: Não é necessário configurar bancos de dados ou serviços externos para testes

## Detalhes de Implementação

### Arquivos para Implementação dos Testes

- `/internal/repository/ticket_test.go`: Testes de integração para operações do repositório de tickets
- `/internal/repository/category_test.go`: Testes de integração para operações do repositório de categorias
- `/internal/service/ticket_test.go`: Testes de camada de serviço com TestContainers
- `/internal/service/category_integration_test.go`: Testes de integração para o serviço de categorias

### Como Funciona

A parte principal da nossa implementação está na função de configuração:

```go
func setupPostgresContainer() (func(), error) {
    ctx := context.Background()

    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15.3-alpine"),
        postgres.WithInitScripts(filepath.Join("../..", "db/migrations", "000001_init_schema.up.sql")),
        postgres.WithInitScripts(filepath.Join("../..", "db/migrations", "000002_seed-categories.up.sql")),
        postgres.WithInitScripts(filepath.Join("../..", "db/migrations", "000004_add_column_user.up.sql")),
        postgres.WithDatabase("ticket-test-db"),
        postgres.WithUsername("postgres"),
        postgres.WithPassword("postgres"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).WithStartupTimeout(5*time.Second)),
    )
    if err != nil {
        return nil, err
    }

    mappedPort, err := container.MappedPort(ctx, "5432")
    if err != nil {
        return nil, err
    }

    connStr = "postgres://postgres:postgres@localhost:" + mappedPort.Port() + "/ticket-test-db?sslmode=disable"

    return func() {
        if err := container.Terminate(ctx); err != nil {
            fmt.Printf("failed to terminate pgContainer: %s", err)
        }
    }, nil
}
```

Esta função:
1. Inicia um contêiner PostgreSQL
2. Aplica migrações de banco de dados para configurar o schema
3. Carrega dados iniciais de teste
4. Retorna uma string de conexão e uma função de limpeza

Na configuração principal do teste:

```go
func TestMain(m *testing.M) {
    cleanup, err := setupPostgresContainer()
    if err != nil {
        panic(fmt.Sprintf("Failed to set up PostgreSQL container: %s", err))
    }
    defer cleanup()

    m.Run()
}
```

### Exemplo de Teste

Aqui está um exemplo de teste de integração do repositório:

```go
func TestTicketRepository_Create_SeverityHigh(t *testing.T) {
    ctx := context.Background()
    conn, err := db.ConnectPostgres(connStr)
    if err != nil {
        t.Fatalf("Failed to connect to the database: %v", err)
    }
    store := New(conn.DB())

    arg := CreateTicketParams{
        Title:       "Login Fails with Correct Credentials",
        Description: "Users can't log in despite correct credentials",
        SeverityID:  1,
        CategoryID:  3,
    }

    ticket, err := store.CreateTicket(ctx, arg)
    assert.NoError(t, err)

    assert.Equal(t, "Login Fails with Correct Credentials", ticket.Title)
    assert.Equal(t, "Users can't log in despite correct credentials", ticket.Description)
    assert.Equal(t, "OPEN", ticket.Status)
    assert.Equal(t, int32(3), ticket.CategoryID)
}
```

## Benefícios Desta Abordagem

1. **Sem necessidade de mocks**: Os testes são executados em um banco de dados PostgreSQL real
2. **Configuração consistente**: O banco de dados está sempre em um estado limpo e conhecido para cada teste
3. **Velocidade**: Os contêineres são leves e iniciam rapidamente
4. **Isolamento**: Cada conjunto de testes obtém sua própria instância de banco de dados
5. **Compatível com CI/CD**: Funciona em qualquer ambiente com suporte a Docker

## Como Começar

### Pré-requisitos

- Go 1.22+
- Docker

### Executando os Testes

```bash
go test ./... -v
```

## Estrutura do Projeto

```
api_with_test_container/
├── cmd/
│   └── main.go              # Ponto de entrada da aplicação
├── config/
│   └── config.go            # Manipulação de configuração
├── db/
│   ├── db.go                # Conexão com o banco de dados
│   ├── migrations/          # Migrações SQL
│   └── queries/             # Queries do Repository SQL
├── docs/
│   └── ...                  # Documentação da API
├── internal/
│   ├── api/                 # Camada de API
│   ├── middleware/          # Middleware HTTP
│   ├── repository/          # Camada de acesso a dados com testes
│   │   ├── category_test.go # Testes de integração com TestContainers
│   │   └── ticket_test.go   # Testes de integração com TestContainers
│   └── service/             # Lógica de negócios
│       ├── category_integration_test.go
│       └── ticket_test.go
├── pkg/
│   └── ...                  # Pacotes compartilhados
├── .env
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── Makefile
```

## Dependências Principais do TestContainers

```go
"github.com/testcontainers/testcontainers-go"
"github.com/testcontainers/testcontainers-go/modules/postgres"
"github.com/testcontainers/testcontainers-go/wait"
```
