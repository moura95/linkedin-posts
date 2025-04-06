package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github/moura95/go-serverless-localstack-postgres/internal/db"
	"github/moura95/go-serverless-localstack-postgres/internal/handler"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	conn, err := db.NewPostgresConnection()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error": "Internal Error"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}
	defer conn.Close()

	if err := db.RunMigrationsWithRetry(conn.DB(), 3); err != nil {
		log.Printf("Error to execute migrations: %v", err)
	}

	ticketHandler := handler.NewTicketHandler(conn.DB())

	return ticketHandler.HandleRequest(ctx, request)
}

func main() {
	lambda.Start(handleRequest)
}
