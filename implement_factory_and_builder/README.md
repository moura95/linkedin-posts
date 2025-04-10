# 🚀 Construindo Aplicações Multi-Database com Factory e Builder em Go

Este repositório demonstra como implementar os padrões Factory e Builder em Go para gerenciar múltiplos bancos de dados na mesma aplicação.

## Estrutura do Projeto

```
implement_factory_and_builder/
├── db_interface.go           # Interface comum para todos os bancos de dados
├── db_builder.go             # Implementação do padrão Builder para Postgres e MongoDB
├── db_factory.go             # Implementação do padrão Factory para os Builders
├── mongo_db.go               # Implementação do MongoDB
├── postgres_db.go            # Implementação do PostgreSQL
├── docker-compose.yaml       # Configuração Docker para testes
└── main.go                   # Exemplo de uso
```

## Padrões Implementados

### Factory Pattern

O Factory Pattern é utilizado para encapsular a lógica de criação de objetos. Neste projeto, usamos uma `DatabaseFactory` que cria diferentes tipos de `DatabaseBuilder` dependendo do tipo de banco de dados solicitado.

**Vantagens:**
- Encapsula a lógica de criação de objetos complexos
- Permite adicionar novos tipos facilmente
- Melhora a manutenção do código

### Builder Pattern

O Builder Pattern é utilizado para construir objetos complexos passo a passo. Neste projeto, implementamos um `DatabaseBuilder` que configura diferentes aspectos de uma conexão de banco de dados (host, porta, credenciais, etc.) antes de finalmente construir o objeto.

**Vantagens:**
- Permite a construção passo a passo de objetos complexos
- Usa a mesma interface de construção para diferentes representações
- Separa a construção de um objeto complexo da sua representação

## Implementação

### Como Funcionam Juntos

1. A `DatabaseFactory` determina qual tipo específico de `DatabaseBuilder` criar
2. O método de construção fluente do `DatabaseBuilder` permite configurar a conexão
3. Finalmente, o método `Build()` cria a instância específica do banco de dados

### Exemplo de Uso

```go
// Instance factory
factory := &DefaultDatabaseFactory{}

// Set Builder
pgBuilder, err := factory.GetBuilder(PostgresDB)
if err != nil {
    log.Fatalf("Error setting builder: %v", err)
}

// Configure to builder
pgBuilder.SetHost("localhost").
    SetPort(5432).
    SetCredentials("postgres", "postgres123").
    SetDatabase("sales_db")

// Build instance
pgDB, err := pgBuilder.Build()
if err != nil {
    log.Fatalf("Error building instance: %v", err)
}

// Connect to instance
err = pgDB.Connect()
if err != nil {
    log.Printf("Error connecting: %v", err)
} else {
    fmt.Printf("Connected: %v\n", pgDB.IsConnected())
    defer pgDB.Disconnect()
}
```

## Como executar

Para executar o exemplo, certifique-se de ter o Docker e Go instalados e execute:

```bash
# Inicie os contêineres Docker para MongoDB e PostgreSQL
docker-compose up -d

# Execute o exemplo
go run .
```

## Quando Usar

- **Factory Pattern**: Quando a criação de um objeto envolve lógica que não é apenas uma simples instanciação
- **Builder Pattern**: Quando um objeto tem muitos parâmetros de construção, alguns opcionais ou quando diferentes representações são necessárias

## Alternativas e Combinações

Este exemplo combina Factory e Builder, mas eles também podem ser usados separadamente:

- Factory sem Builder: Quando a configuração do objeto é simples
- Builder sem Factory: Quando há apenas um tipo de objeto a ser construído

## Extensões Possíveis

- Adicionar validação nos métodos Builder
- Criar um Director para coordenar a construção de objetos com configurações predefinidas
- Adicionar suporte para outros bancos de dados (MySQL, SQLite, etc.)
- Implementar lógica de failover entre bancos de dados
- Criar uma camada ORM agnóstica que funcione com qualquer banco implementado

## Por que usar múltiplos bancos de dados?

Aplicações modernas frequentemente necessitam de diferentes tipos de bancos de dados para diferentes casos de uso:

- **PostgreSQL** para dados relacionais e transações complexas
- **MongoDB** para armazenamento de documentos e consultas flexíveis
- **Redis** para cache e filas de mensagens
- **Elasticsearch** para busca de texto completo

Com a abordagem de Factory e Builder, você pode facilmente integrar múltiplos bancos de dados mantendo seu código limpo e gerenciável!

> 💡 **Dica Pro**: Esta mesma abordagem pode ser estendida para outros tipos de serviços além de bancos de dados, como sistemas de mensageria, serviços de armazenamento em nuvem, ou qualquer dependência externa que sua aplicação precise gerenciar.