package model

import (
	"database/sql"
	"time"
)

type Ticket struct {
	ID            int32         `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	Status        string        `json:"status"`
	SeverityID    int32         `json:"severity_id"`
	CategoryID    int32         `json:"category_id"`
	SubcategoryID sql.NullInt32 `json:"subcategory_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	CompletedAt   sql.NullTime  `json:"completed_at,omitempty"`
}

type CreateTicketRequest struct {
	Title         string `json:"title" validate:"required"`
	Description   string `json:"description"`
	SeverityID    int32  `json:"severity_id" validate:"required"`
	CategoryID    int32  `json:"category_id" validate:"required"`
	SubcategoryID *int32 `json:"subcategory_id,omitempty"`
}

type UpdateTicketRequest struct {
	Title         *string `json:"title,omitempty"`
	Description   *string `json:"description,omitempty"`
	Status        *string `json:"status,omitempty"`
	SeverityID    *int32  `json:"severity_id,omitempty"`
	CategoryID    *int32  `json:"category_id,omitempty"`
	SubcategoryID *int32  `json:"subcategory_id,omitempty"`
}

type TicketDetailResponse struct {
	ID              int32      `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Status          string     `json:"status"`
	SeverityID      int32      `json:"severity_id"`
	SeverityName    string     `json:"severity_name"`
	CategoryID      int32      `json:"category_id"`
	CategoryName    string     `json:"category_name"`
	SubcategoryID   *int32     `json:"subcategory_id,omitempty"`
	SubcategoryName *string    `json:"subcategory_name,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type Category struct {
	ID       int32         `json:"id"`
	Name     string        `json:"name"`
	ParentID sql.NullInt32 `json:"parent_id,omitempty"`
}

type Severity struct {
	ID          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
