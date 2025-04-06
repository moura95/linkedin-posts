package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	AuthDB   string
}

type MongoDatabase struct {
	BaseDatabase
	Config MongoConfig
	Client *mongo.Client
	DB     *mongo.Database
	ctx    context.Context
}

func (db *MongoDatabase) Connect() error {
	authDB := db.Config.AuthDB
	if authDB == "" {
		authDB = "admin"
	}

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		db.Config.User, db.Config.Password, db.Config.Host, db.Config.Port, db.Config.Database)

	uri = fmt.Sprintf("%s?authSource=%s", uri, authDB)

	clientOptions := options.Client().
		ApplyURI(uri).
		SetRetryWrites(true).
		SetRetryReads(true)

	db.ctx = context.Background()
	ctx, cancel := context.WithTimeout(db.ctx, 10*time.Second)
	defer cancel()

	var err error
	db.Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = db.Client.Ping(ctx, readpref.Primary())
	if err != nil {
		db.Client.Disconnect(ctx)
		return fmt.Errorf("error pinging MongoDB: %w", err)
	}

	db.DB = db.Client.Database(db.Config.Database)

	db.Connected = true
	fmt.Printf("Connected to MongoDB: %s@%s:%d/%s (auth: %s)\n",
		db.Config.User, db.Config.Host, db.Config.Port, db.Config.Database, authDB)

	return nil
}

func (db *MongoDatabase) Disconnect() error {
	if db.Connected && db.Client != nil {
		ctx, cancel := context.WithTimeout(db.ctx, 5*time.Second)
		defer cancel()

		err := db.Client.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("error disconnecting from MongoDB: %w", err)
		}
		db.Connected = false
		fmt.Println("Disconnected from MongoDB")
	}
	return nil
}
