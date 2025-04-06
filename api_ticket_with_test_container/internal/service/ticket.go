package service

import (
	"context"
	"database/sql"
	"fmt"

	"github/moura95/ticket-api/config"
	"github/moura95/ticket-api/internal/repository"
	"go.uber.org/zap"
)

type TicketService struct {
	repository repository.Querier
	config     config.Config
	logger     *zap.SugaredLogger
}

func NewTicketService(repo repository.Querier, cfg config.Config, log *zap.SugaredLogger) *TicketService {
	return &TicketService{
		repository: repo,
		config:     cfg,
		logger:     log,
	}
}

func (s *TicketService) GetAll(ctx context.Context) ([]repository.Ticket, error) {
	tickets, err := s.repository.ListTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket %s", err.Error())
	}

	return tickets, nil
}

func (s *TicketService) GetByID(ctx context.Context, id int32) (*repository.GetTicketByIdRow, error) {
	ticket, err := s.repository.GetTicketById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get ticket %s", err.Error())
	}

	return &ticket, nil
}

func (s *TicketService) Create(ctx context.Context, title, desc string, severityId, categoryId, subCategoryId int32) (*repository.Ticket, error) {
	arg := repository.CreateTicketParams{
		Title:       title,
		Description: desc,
		SeverityID:  severityId,
		CategoryID:  categoryId,
		SubcategoryID: sql.NullInt32{
			Int32: subCategoryId,
			Valid: subCategoryId > 0,
		},
	}
	ticket, err := s.repository.CreateTicket(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s", err.Error())
	}
	return &ticket, nil
}

func (s *TicketService) Delete(ctx context.Context, id int32) error {
	err := s.repository.DeleteTicket(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get ticket %s", err.Error())
	}

	return nil
}

func (s *TicketService) Update(ctx context.Context, id, userId, severityId, categoryId, subCategoryId int32, title, desc, status string) error {
	arg := repository.UpdateTicketParams{
		ID: id,
		Title: sql.NullString{
			String: title,
			Valid:  title != "",
		},
		Description: sql.NullString{
			String: desc,
			Valid:  desc != "",
		},
		UserID: sql.NullInt32{
			Int32: userId,
			Valid: userId > 0,
		},
		SeverityID: sql.NullInt32{
			Int32: severityId,
			Valid: severityId > 0 && severityId <= 4,
		},
		Status: sql.NullString{
			String: status,
			Valid:  status != "",
		},
		CategoryID: sql.NullInt32{
			Int32: categoryId,
			Valid: categoryId > 0,
		},
		SubcategoryID: sql.NullInt32{
			Int32: subCategoryId,
			Valid: subCategoryId > 0,
		},
	}
	err := s.repository.UpdateTicket(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update %s", err.Error())
	}
	return nil
}
