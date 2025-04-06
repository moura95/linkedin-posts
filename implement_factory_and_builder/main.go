package main

import (
	"fmt"
	"log"
)

func main() {
	// Instance factory
	factory := &DefaultDatabaseFactory{}

	// Connect to PostgreSQL
	fmt.Println("=== Connecting to PostgreSQL ===")

	// Set Builder
	pgBuilder, err := factory.GetBuilder(PostgresDB)
	if err != nil {
		log.Fatalf("Error setting builder: %v", err)
	}

	// Configure to builder
	pgBuilder.SetHost("localhost").
		SetPort(5432).
		SetCredentials("postgres", "postgres123").
		SetDatabase("sales_db")

	// Build instance
	pgDB, err := pgBuilder.Build()
	if err != nil {
		log.Fatalf("Error building instance: %v", err)
	}

	// Connect to instance
	err = pgDB.Connect()
	if err != nil {
		log.Printf("Error connecting: %v", err)
	} else {
		fmt.Printf("Connected: %v\n", pgDB.IsConnected())
		defer pgDB.Disconnect()
	}

	// Connect to MongoDB
	fmt.Println("\n=== Connect to MongoDB ===")

	// Set Builder
	mongoBuilder, err := factory.GetBuilder(MongoDB)
	if err != nil {
		log.Fatalf("Error setting builder: %v", err)
	}

	// Configure to builder
	mongoBuilder.SetHost("localhost").
		SetPort(27017).
		SetCredentials("mongodb", "mongodb123").
		SetDatabase("inventory_db")

	// Build instance
	mongoDB, err := mongoBuilder.Build()
	if err != nil {
		log.Fatalf("Error building instance: %v", err)
	}

	// Connect to instance
	err = mongoDB.Connect()
	if err != nil {
		log.Printf("Error connecting: %v", err)
	} else {
		fmt.Printf("Connected: %v\n", mongoDB.IsConnected())
		defer mongoDB.Disconnect()
	}

	fmt.Println("\n=== Finished ===")
}
