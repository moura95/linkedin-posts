package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github/moura95/go-serverless-localstack-postgres/internal/model"
)

type TicketRepository interface {
	GetAll(ctx context.Context) ([]model.Ticket, error)
	GetByID(ctx context.Context, id int32) (*model.TicketDetailResponse, error)
	Create(ctx context.Context, ticket model.CreateTicketRequest) (*model.Ticket, error)
	Update(ctx context.Context, id int32, ticket model.UpdateTicketRequest) (*model.Ticket, error)
	Delete(ctx context.Context, id int32) error
}

type PostgresTicketRepository struct {
	db *sql.DB
}

func NewPostgresTicketRepository(db *sql.DB) TicketRepository {
	return &PostgresTicketRepository{db: db}
}

func (r *PostgresTicketRepository) GetAll(ctx context.Context) ([]model.Ticket, error) {
	query := `
		SELECT id, title, description, status, severity_id, category_id, subcategory_id, 
		       created_at, updated_at, completed_at
		FROM tickets
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching tickets: %w", err)
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var t model.Ticket
		if err := rows.Scan(
			&t.ID, &t.Title, &t.Description, &t.Status, &t.SeverityID,
			&t.CategoryID, &t.SubcategoryID, &t.CreatedAt, &t.UpdatedAt, &t.CompletedAt,
		); err != nil {
			return nil, fmt.Errorf("error reading ticket: %w", err)
		}
		tickets = append(tickets, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tickets: %w", err)
	}

	return tickets, nil
}

func (r *PostgresTicketRepository) GetByID(ctx context.Context, id int32) (*model.TicketDetailResponse, error) {
	query := `
		SELECT t.id, t.title, t.description, t.status, 
		       t.severity_id, s.name AS severity_name,
		       t.category_id, c.name AS category_name,
		       t.subcategory_id, sc.name AS subcategory_name,
		       t.created_at, t.updated_at, t.completed_at
		FROM tickets t
		JOIN severities s ON t.severity_id = s.id
		JOIN categories c ON t.category_id = c.id
		LEFT JOIN categories sc ON t.subcategory_id = sc.id
		WHERE t.id = $1
	`

	var ticket model.TicketDetailResponse
	var subcategoryID sql.NullInt32
	var subcategoryName sql.NullString
	var completedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
		&ticket.SeverityID, &ticket.SeverityName,
		&ticket.CategoryID, &ticket.CategoryName,
		&subcategoryID, &subcategoryName,
		&ticket.CreatedAt, &ticket.UpdatedAt, &completedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("ticket not found: %w", err)
		}
		return nil, fmt.Errorf("error fetching ticket: %w", err)
	}

	// Convert nullables to pointers
	if subcategoryID.Valid {
		ticket.SubcategoryID = &subcategoryID.Int32
	}

	if subcategoryName.Valid {
		ticket.SubcategoryName = &subcategoryName.String
	}

	if completedAt.Valid {
		ticket.CompletedAt = &completedAt.Time
	}

	return &ticket, nil
}

func (r *PostgresTicketRepository) Create(ctx context.Context, req model.CreateTicketRequest) (*model.Ticket, error) {
	query := `
		INSERT INTO tickets (title, description, status, severity_id, category_id, subcategory_id, created_at, updated_at)
		VALUES ($1, $2, 'OPEN', $3, $4, $5, NOW(), NOW())
		RETURNING id, title, description, status, severity_id, category_id, subcategory_id, created_at, updated_at, completed_at
	`

	var subcategoryID sql.NullInt32
	if req.SubcategoryID != nil {
		subcategoryID = sql.NullInt32{Int32: *req.SubcategoryID, Valid: true}
	}

	var ticket model.Ticket
	err := r.db.QueryRowContext(ctx, query,
		req.Title, req.Description, req.SeverityID, req.CategoryID, subcategoryID,
	).Scan(
		&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
		&ticket.SeverityID, &ticket.CategoryID, &ticket.SubcategoryID,
		&ticket.CreatedAt, &ticket.UpdatedAt, &ticket.CompletedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating ticket: %w", err)
	}

	return &ticket, nil
}

func (r *PostgresTicketRepository) Update(ctx context.Context, id int32, req model.UpdateTicketRequest) (*model.Ticket, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM tickets WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking ticket: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("ticket not found")
	}

	query := "UPDATE tickets SET updated_at = NOW()"
	args := []interface{}{}
	paramCount := 1

	if req.Title != nil {
		query += fmt.Sprintf(", title = $%d", paramCount)
		args = append(args, *req.Title)
		paramCount++
	}

	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", paramCount)
		args = append(args, *req.Description)
		paramCount++
	}

	if req.Status != nil {
		query += fmt.Sprintf(", status = $%d", paramCount)
		args = append(args, *req.Status)
		paramCount++

		// If status changed to DONE or CLOSED, set completed_at
		if *req.Status == "DONE" || *req.Status == "CLOSED" {
			query += ", completed_at = NOW()"
		}
	}

	if req.SeverityID != nil {
		query += fmt.Sprintf(", severity_id = $%d", paramCount)
		args = append(args, *req.SeverityID)
		paramCount++
	}

	if req.CategoryID != nil {
		query += fmt.Sprintf(", category_id = $%d", paramCount)
		args = append(args, *req.CategoryID)
		paramCount++
	}

	if req.SubcategoryID != nil {
		query += fmt.Sprintf(", subcategory_id = $%d", paramCount)
		args = append(args, *req.SubcategoryID)
		paramCount++
	}

	// Adding the WHERE condition and RETURNING
	query += fmt.Sprintf(" WHERE id = $%d ", paramCount)
	args = append(args, id)

	query += `
		RETURNING id, title, description, status, severity_id, category_id, 
		subcategory_id, created_at, updated_at, completed_at
	`

	var ticket model.Ticket
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&ticket.ID, &ticket.Title, &ticket.Description, &ticket.Status,
		&ticket.SeverityID, &ticket.CategoryID, &ticket.SubcategoryID,
		&ticket.CreatedAt, &ticket.UpdatedAt, &ticket.CompletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating ticket: %w", err)
	}

	return &ticket, nil
}

func (r *PostgresTicketRepository) Delete(ctx context.Context, id int32) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tickets WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found")
	}

	return nil
}
