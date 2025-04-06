package service

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github/moura95/ticket-api/db"
	"github/moura95/ticket-api/internal/repository"
)

var (
	pgContainer testcontainers.Container
	connStr     string
)

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

func TestMain(m *testing.M) {
	cleanup, err := setupPostgresContainer()
	if err != nil {
		panic(fmt.Sprintf("Failed to set up PostgreSQL container: %s", err))
	}
	defer cleanup()

	m.Run()
}

func NewTicketServiceTest(repo repository.Querier) *TicketService {
	return &TicketService{
		repository: repo,
	}
}

type createRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	SeverityId    int32  `json:"severity_id"`
	CategoryId    int32  `json:"category_id"`
	SubCategoryId int32  `json:"subcategory_id"`
}

func TestCreateTicket(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewTicketServiceTest(store)
	req := createRequest{
		Title:       "Install Docker",
		Description: "I need Install docker for work",
		SeverityId:  1,
		CategoryId:  3,
	}

	ti, err := service.Create(ctx, req.Title, req.Description, req.SeverityId, req.CategoryId, req.SubCategoryId)
	assert.NoError(t, err)

	ticket, _ := service.GetByID(ctx, ti.ID)
	assert.Equal(t, "Install Docker", ticket.Title)
	assert.Equal(t, "I need Install docker for work", ticket.Description)
	assert.Equal(t, int32(1), ticket.SeverityID)
	assert.Equal(t, "OPEN", ticket.Status)

	assert.NoError(t, err)

}

func TestCreateTicket1(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewTicketServiceTest(store)

	req := createRequest{
		Title:       "Install Vpn",
		Description: "I need install vpn for connect database",
		SeverityId:  2,
		CategoryId:  2,
	}

	ti, err := service.Create(ctx, req.Title, req.Description, req.SeverityId, req.CategoryId, req.SubCategoryId)
	assert.NoError(t, err)

	ticket, _ := service.GetByID(ctx, ti.ID)
	assert.Equal(t, "Install Vpn", ticket.Title)
	assert.Equal(t, "I need install vpn for connect database", ticket.Description)
	assert.Equal(t, int32(2), ticket.SeverityID)
	assert.Equal(t, int32(2), ticket.CategoryID)
	assert.Equal(t, "OPEN", ticket.Status)

	assert.NoError(t, err)

}

func TestGetTicket(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewTicketServiceTest(store)

	ticket, err := service.GetByID(ctx, 2)

	assert.NoError(t, err)
	assert.Equal(t, "Install Vpn", ticket.Title)
	assert.Equal(t, "I need install vpn for connect database", ticket.Description)
	assert.Equal(t, int32(2), ticket.SeverityID)
	assert.Equal(t, int32(2), ticket.CategoryID)
	assert.Equal(t, "OPEN", ticket.Status)

}

func TestList(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewTicketServiceTest(store)

	tickets, err := service.GetAll(ctx)

	assert.NoError(t, err)
	assert.Equal(t, len(tickets), 2)
	assert.Equal(t, "Install Docker", tickets[0].Title)
	assert.Equal(t, "I need Install docker for work", tickets[0].Description)
	assert.Equal(t, int32(1), tickets[0].SeverityID)
	assert.Equal(t, "OPEN", tickets[0].Status)

	// Ticket 2

	assert.Equal(t, "Install Vpn", tickets[1].Title)
	assert.Equal(t, "I need install vpn for connect database", tickets[1].Description)
	assert.Equal(t, int32(2), tickets[1].SeverityID)
	assert.Equal(t, int32(2), tickets[1].CategoryID)
	assert.Equal(t, "OPEN", tickets[1].Status)
}

func TestDeleteTicket(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewTicketServiceTest(store)

	tickets, err := service.GetAll(ctx)
	assert.NoError(t, err)

	// delete 1 item
	err = service.Delete(ctx, 2)
	assert.NoError(t, err)

	// confirm list
	newListTickets, err := service.GetAll(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(tickets)-1, len(newListTickets))

}
