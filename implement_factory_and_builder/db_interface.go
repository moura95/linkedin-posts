package main

import (
	"errors"
)

var (
	UnsupportedDatabaseErr = errors.New("type unsupported database")
	ConfigurationErr       = errors.New("invalid configuration")
)

type Database interface {
	Connect() error
	Disconnect() error
	IsConnected() bool
	GetName() string
}

type BaseDatabase struct {
	Connected bool
	Name      string
}

func (db *BaseDatabase) IsConnected() bool {
	return db.Connected
}

func (db *BaseDatabase) GetName() string {
	return db.Name
}
