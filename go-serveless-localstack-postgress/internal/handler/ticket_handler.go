package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github/moura95/go-serverless-localstack-postgres/internal/model"
	"github/moura95/go-serverless-localstack-postgres/internal/repository"
)

type TicketHandler struct {
	repo repository.TicketRepository
}

func NewTicketHandler(db *sql.DB) *TicketHandler {
	repo := repository.NewPostgresTicketRepository(db)
	return &TicketHandler{
		repo: repo,
	}
}

func (h *TicketHandler) HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Request: %s %s", request.HTTPMethod, request.Path)

	switch request.HTTPMethod {
	case "GET":
		if id, ok := request.PathParameters["id"]; ok && id != "" {
			return h.getTicket(ctx, id)
		}
		return h.getAllTickets(ctx)
	case "POST":
		return h.createTicket(ctx, request)
	case "PUT", "PATCH":
		if id, ok := request.PathParameters["id"]; ok && id != "" {
			return h.updateTicket(ctx, id, request)
		}
	case "DELETE":
		if id, ok := request.PathParameters["id"]; ok && id != "" {
			return h.deleteTicket(ctx, id)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Body:       `{"error": "Method not allowed"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (h *TicketHandler) getAllTickets(ctx context.Context) (events.APIGatewayProxyResponse, error) {
	tickets, err := h.repo.GetAll(ctx)
	if err != nil {
		log.Printf("Error get all tickets: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Error get tickets: %s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseJSON, err := json.Marshal(tickets)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Error serializing tickets"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseJSON),
	}, nil
}

func (h *TicketHandler) getTicket(ctx context.Context, idStr string) (events.APIGatewayProxyResponse, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid ID "}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	ticket, err := h.repo.GetByID(ctx, int32(id))
	if err != nil {
		log.Printf("Error to get ticket %d: %v", id, err)
		if err.Error() == "sql: no rows in result set" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error": "Ticket not found: %d"}`, id),
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Error to get ticket: %s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseJSON, err := json.Marshal(ticket)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Error serializing ticket"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseJSON),
	}, nil
}

func (h *TicketHandler) createTicket(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var req model.CreateTicketRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid JSON "}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	if req.Title == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Title is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	ticket, err := h.repo.Create(ctx, req)
	if err != nil {
		log.Printf("Error to create ticket: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Error to create ticket: %s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseJSON, err := json.Marshal(ticket)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Error serializing ticket"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseJSON),
	}, nil
}

func (h *TicketHandler) updateTicket(ctx context.Context, idStr string, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid ID"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var req model.UpdateTicketRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid JSON"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	ticket, err := h.repo.Update(ctx, int32(id), req)
	if err != nil {
		log.Printf("Error to update ticket %d: %v", id, err)
		if err.Error() == "Ticket not found" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error": "Ticket not found: %d"}`, id),
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Error to update ticket: %s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseJSON, err := json.Marshal(ticket)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Error serializing ticket"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseJSON),
	}, nil
}

func (h *TicketHandler) deleteTicket(ctx context.Context, idStr string) (events.APIGatewayProxyResponse, error) {
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid ID "}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	if err := h.repo.Delete(ctx, int32(id)); err != nil {
		log.Printf("Error to delete ticket %d: %v", id, err)
		if err.Error() == "ticket not found" {
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusNotFound,
				Body:       fmt.Sprintf(`{"error": "Error to delete ticket: %d"}`, id),
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf(`{"error": "Error to delete ticket: %s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       `{"message": "Ticket deleted"}`,
	}, nil
}
