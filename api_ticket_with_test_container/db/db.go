package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection interface {
	Close() error
	DB() *sqlx.DB
}

type conn struct {
	db *sqlx.DB
}

func ConnectPostgres(connStr string) (Connection, error) {

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &conn{db: db}, nil
}

func (c *conn) Close() error {
	return c.db.Close()
}

func (c *conn) DB() *sqlx.DB {
	return c.db
}
