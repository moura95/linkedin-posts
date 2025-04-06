package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github/moura95/ticket-api/config"
	"github/moura95/ticket-api/internal/repository"
)

func TestTicketService_GetAll(t *testing.T) {

	tests := []struct {
		name          string
		setupMocked   func(repo *repository.QuerierMocked)
		expected      []repository.Ticket
		expectedError error
	}{
		{
			name: "Get All Tickets",
			setupMocked: func(repo *repository.QuerierMocked) {
				tickets := []repository.Ticket{
					{
						ID:          1,
						Title:       "Login Fails with Correct Credentials",
						Status:      "OPEN",
						Description: "Users can't log in despite correct credentials",
						SeverityID:  1,
						CategoryID:  2,
					},
					{
						ID:          2,
						Title:       "Printer Not Responding",
						Status:      "IN_PROGRESS",
						Description: "Office printer is not responding to any commands",
						SeverityID:  2,
						CategoryID:  1,
					},
				}
				repo.EXPECT().ListTickets(mock.Anything).Return(tickets, nil)

			},
			expected: []repository.Ticket{
				{
					ID:          1,
					Title:       "Login Fails with Correct Credentials",
					Status:      "OPEN",
					Description: "Users can't log in despite correct credentials",
					SeverityID:  1,
					CategoryID:  2,
				},
				{
					ID:          2,
					Title:       "Printer Not Responding",
					Status:      "IN_PROGRESS",
					Description: "Office printer is not responding to any commands",
					SeverityID:  2,
					CategoryID:  1,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewTicketService(repoMocked, config.Config{}, nil)
			tickets, err := service.GetAll(context.Background())
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expected, tickets)
		})
	}
}

func TestTicketService_Create(t *testing.T) {

	tests := []struct {
		name          string
		title         string
		desc          string
		severityId    int32
		categoryId    int32
		subCategoryId int32
		setupMocked   func(repo *repository.QuerierMocked)
		expected      *repository.Ticket
		expectedError error
	}{
		{
			name: "Create Ticket Printer not responding",
			setupMocked: func(repo *repository.QuerierMocked) {
				tickt := repository.Ticket{
					ID:          1,
					Title:       "Printer Not Responding",
					Status:      "OPEN",
					Description: "Office printer is not responding to any commands",
					SeverityID:  2,
					CategoryID:  1,
				}
				repo.EXPECT().CreateTicket(mock.Anything, mock.Anything).Return(tickt, nil)

			},
			expected: &repository.Ticket{
				ID:          1,
				Title:       "Printer Not Responding",
				Status:      "OPEN",
				Description: "Office printer is not responding to any commands",
				SeverityID:  2,
				CategoryID:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewTicketService(repoMocked, config.Config{}, nil)
			ticket, err := service.Create(context.Background(), tt.title, tt.desc, tt.severityId, tt.categoryId, tt.subCategoryId)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expected, ticket)
		})
	}
}

func TestTicketService_Delete(t *testing.T) {

	tests := []struct {
		name          string
		id            int32
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError error
	}{
		{
			name: "Delete Ticket By Id 1",
			id:   1,
			setupMocked: func(repo *repository.QuerierMocked) {
				repo.EXPECT().DeleteTicket(mock.Anything, mock.Anything).Return(nil)

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewTicketService(repoMocked, config.Config{}, nil)
			err := service.Delete(context.Background(), tt.id)
			assert.Equal(t, tt.expectedError, err)

		})
	}
}
