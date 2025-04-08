# Implementação dos Padrões de Design Factory e Strategy em Go

Este repositório demonstra uma implementação prática dos padrões de design Factory e Strategy em Go, utilizando como exemplo um sistema de notificações multicanal.

## Visão Geral

O sistema permite enviar notificações através de diferentes canais (Email, SMS, Push, Slack) usando uma interface unificada. Os padrões implementados permitem:

- Alternar facilmente entre diferentes canais de notificação
- Adicionar novos canais sem modificar código existente
- Manter uma estrutura de código limpa e manutenível

## Padrões de Design Utilizados

### Padrão Strategy

O padrão Strategy permite definir uma família de estrategias, encapsular cada uma deles em classes separadas e torná-los intercambiáveis. No nosso sistema:

- `NotificationStrategy` é a interface que define o comportamento comum
- `EmailStrategy`, `SMSStrategy`, `PushStrategy` e `SlackStrategy` são as implementações específicas
- Cada estratégia encapsula a lógica de envio para um canal específico

### Padrão Factory

O padrão Factory fornece uma interface para criar objetos sem especificar suas classes concretas. No nosso sistema:

- `NotificationFactory` é a interface que define o método de criação
- `notificationFactory` é a implementação concreta que cria as estratégias apropriadas
- O cliente não precisa conhecer os detalhes de criação das estratégias

## Estrutura do Projeto

```
implement_factory_and_strategy/
├── main.go                      # Ponto de entrada da aplicação
├── notification_service.go      # Serviço que utiliza as estratégias
├── notification_factory.go      # Factory para criar estratégias
└── strategies/
    ├── notification_interface.go  # Interface e base das estratégias
    ├── email_strategy.go          # Implementação para Email
    ├── sms_strategy.go            # Implementação para SMS
    ├── push_strategy.go           # Implementação para Push
    └── slack_strategy.go          # Implementação para Slack
```

## Código Chave Explicado

### Interface Strategy

```go
// strategies/notification_interface.go
type NotificationStrategy interface {
    Send(message string, recipient string) error
    GetChannelName() string
}
```

Esta interface define o contrato para todas as estratégias de notificação. Qualquer novo canal deve implementar estes métodos.

### Factory

```go
// notification_factory.go
type NotificationFactory interface {
    CreateStrategy(channelType string) (strategy.NotificationStrategy, error)
}

func (f *notificationFactory) CreateStrategy(channelType string) (strategy.NotificationStrategy, error) {
    switch channelType {
    case "email":
        return &strategy.EmailStrategy{}, nil
    case "sms":
        return &strategy.SMSStrategy{}, nil
    // ... outros casos
    }
}
```

A factory decide qual implementação concreta criar com base no tipo solicitado, abstraindo essa lógica do cliente.

### Serviço de Notificação

```go
// notification_service.go
type NotificationService interface {
    SetStrategy(channelType string) error
    Notify(message string, recipient string) error
    GetCurrentChannel() string
}
```

O serviço utiliza a estratégia atual para enviar notificações e pode alternar entre diferentes estratégias em tempo de execução.

## Como Executar

```bash
git clone https://github.com/moura95/linkedin-posts/tree/main/implement_factory_and_strategy
cd implement_factory_and_strategy
go run  main.go notification_factory.go notification_service.go
```

## Exemplo de Uso

```go
// Criar factory e serviço com estratégia padrão (email)
factory := &notificationFactory{}
service, _ := NewNotificationService(factory, "email")

// Enviar notificação por email
service.Notify("Olá!", "usuario@exemplo.com")

// Mudar para SMS e enviar novamente
service.SetStrategy("sms")
service.Notify("Olá!", "+5511987654321")
```

## Benefícios da Implementação

1. **Desacoplamento**: O serviço de notificação não conhece os detalhes das implementações específicas
2. **Extensibilidade**: Novos canais podem ser adicionados sem modificar o código existente
3. **Flexibilidade**: Estratégias podem ser trocadas dinamicamente em tempo de execução
4. **Testabilidade**: Cada componente pode ser testado de forma isolada

## Adicionando um Novo Canal

Para adicionar um novo canal (ex: WhatsApp), basta:

1. Criar uma nova estratégia que implemente a interface `NotificationStrategy`
2. Adicionar um novo caso na factory para criar a nova estratégia

Todo o restante do código continuará funcionando sem alterações!