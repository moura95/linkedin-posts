package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

type PostgresDatabase struct {
	BaseDatabase
	Config PostgresConfig
	DB     *sql.DB
}

func (db *PostgresDatabase) Connect() error {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Config.Host, db.Config.Port, db.Config.User,
		db.Config.Password, db.Config.Database, db.Config.SSLMode,
	)

	var err error
	db.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return fmt.Errorf("Error connecting to PostgreSQL: %w", err)
	}

	err = db.DB.Ping()
	if err != nil {
		db.DB.Close()
		return fmt.Errorf("Error pinging PostgreSQL: %w", err)
	}

	db.DB.SetMaxOpenConns(25)
	db.DB.SetMaxIdleConns(5)

	db.Connected = true
	fmt.Printf("Connected PostgreSQL: %s@%s:%d/%s\n",
		db.Config.User, db.Config.Host, db.Config.Port, db.Config.Database)

	return nil
}

func (db *PostgresDatabase) Disconnect() error {
	if db.Connected && db.DB != nil {
		err := db.DB.Close()
		if err != nil {
			return fmt.Errorf("Error disconnecting from PostgreSQL: %w", err)
		}
		db.Connected = false
		fmt.Println("Disconnected from PostgreSQL")
	}
	return nil
}
