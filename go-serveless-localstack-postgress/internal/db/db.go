package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Connection interface {
	Close() error
	DB() *sql.DB
}

type PostgresConnection struct {
	db *sql.DB
}

func NewPostgresConnection() (Connection, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "tickets_db")
	sslMode := getEnv("DB_SSL_MODE", "disable")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir conexão com PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("erro ao ping PostgreSQL: %w", err)
	}

	log.Printf("Conectado ao PostgreSQL em %s:%s/%s", host, port, dbName)

	return &PostgresConnection{db: db}, nil
}

func (c *PostgresConnection) Close() error {
	return c.db.Close()
}

func (c *PostgresConnection) DB() *sql.DB {
	return c.db
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
