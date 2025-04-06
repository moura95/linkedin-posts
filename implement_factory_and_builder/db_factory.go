package main

import (
	"fmt"
)

type DBType string

const (
	PostgresDB DBType = "postgres"
	MongoDB    DBType = "mongodb"
)

type DatabaseFactory interface {
	GetBuilder(dbType DBType) (DatabaseBuilder, error)
}

type DefaultDatabaseFactory struct {
}

func NewDefaultDatabaseFactory() DatabaseFactory {
	return &DefaultDatabaseFactory{}
}

func (f *DefaultDatabaseFactory) GetBuilder(dbType DBType) (DatabaseBuilder, error) {
	switch dbType {
	case PostgresDB:
		return NewPostgresBuilder(), nil
	case MongoDB:
		return NewMongoBuilder(), nil
	default:
		return nil, fmt.Errorf("%w: %s", UnsupportedDatabaseErr, dbType)
	}
}
