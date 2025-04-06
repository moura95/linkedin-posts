package repository

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github/moura95/ticket-api/db"
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

func TestTicketRepository_Create_SeverityMedium(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := CreateTicketParams{
		Title:       "Submit Button Not Working on Contact Page",
		Description: "Submit button is unresponsive on contact form",
		SeverityID:  2,
		CategoryID:  2,
	}

	ticket, err := store.CreateTicket(ctx, arg)
	assert.NoError(t, err)

	assert.Equal(t, "Submit Button Not Working on Contact Page", ticket.Title)
	assert.Equal(t, "Submit button is unresponsive on contact form", ticket.Description)
	assert.Equal(t, "OPEN", ticket.Status)
	assert.Equal(t, int32(2), ticket.CategoryID)

}

func TestTicketRepository_Create_Failed(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := CreateTicketParams{
		Title:       "",
		Description: "Submit button is unresponsive on contact form",
		SeverityID:  999,
	}

	_, err = store.CreateTicket(ctx, arg)

	assert.Equal(t, err.Error(), "pq: insert or update on table \"tickets\" violates foreign key constraint \"tickets_severity_id_fkey\"")

}

func TestTicketRepository_GetById(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())
	id := 1
	ticket, err := store.GetTicketById(ctx, int32(id))
	assert.NoError(t, err)

	assert.Equal(t, "Login Fails with Correct Credentials", ticket.Title)
	assert.Equal(t, "Users can't log in despite correct credentials", ticket.Description)
	assert.Equal(t, "OPEN", ticket.Status)
	assert.Equal(t, int32(1), ticket.SeverityID)

}

func TestTicketRepository_GetAll(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())
	tickets, err := store.ListTickets(ctx)
	assert.NoError(t, err)

	// ticket 1
	assert.Equal(t, "Login Fails with Correct Credentials", tickets[0].Title)
	assert.Equal(t, "Users can't log in despite correct credentials", tickets[0].Description)
	assert.Equal(t, "OPEN", tickets[0].Status)
	assert.Equal(t, int32(1), tickets[0].SeverityID)
	assert.Equal(t, int32(3), tickets[0].CategoryID)

	// ticket 2
	assert.Equal(t, "Submit Button Not Working on Contact Page", tickets[1].Title)
	assert.Equal(t, "Submit button is unresponsive on contact form", tickets[1].Description)
	assert.Equal(t, "OPEN", tickets[1].Status)
	assert.Equal(t, int32(2), tickets[1].SeverityID)
	assert.Equal(t, int32(2), tickets[1].CategoryID)

}

func TestTicketRepository_UpdateTitleTicket(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := UpdateTicketParams{
		ID: 1,
		Title: sql.NullString{
			String: "Title updated",
			Valid:  true,
		},
	}

	err = store.UpdateTicket(ctx, arg)
	assert.NoError(t, err)

	// Get Ticket for confirm update

	ticket, err := store.GetTicketById(ctx, arg.ID)
	assert.Equal(t, "Title updated", ticket.Title)

}

func TestTicketRepository_UpdateTitleStatusDone(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := UpdateTicketParams{
		ID: 1,
		Status: sql.NullString{
			String: "DONE",
			Valid:  true,
		},
	}

	err = store.UpdateTicket(ctx, arg)
	assert.NoError(t, err)

	// Get Ticket for confirm update

	ticket, err := store.GetTicketById(ctx, arg.ID)
	assert.Equal(t, "DONE", ticket.Status)
}

func TestTicketRepository_UpdateAllFieldsTicket(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := UpdateTicketParams{
		ID: 2,
		Title: sql.NullString{
			String: "New Title",
			Valid:  true,
		},
		Status: sql.NullString{
			String: "CLOSED",
			Valid:  true,
		},
		Description: sql.NullString{
			String: "New Desc",
			Valid:  true,
		},
	}

	err = store.UpdateTicket(ctx, arg)
	assert.NoError(t, err)

	// Get Ticket for confirm update

	ticket, err := store.GetTicketById(ctx, arg.ID)
	assert.Equal(t, "New Title", ticket.Title)
	assert.Equal(t, "New Desc", ticket.Description)
	assert.Equal(t, "CLOSED", ticket.Status)
}

func TestTicketRepository_Delete(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())
	id := 1
	err = store.DeleteTicket(ctx, int32(id))
	assert.NoError(t, err)

}
