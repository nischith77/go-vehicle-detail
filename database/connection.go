package db

import (
	"database/sql"
	"fmt"
	"go-vehicle-detail/models"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func ConnectDB(connStr string) *sql.DB {
	// Example connection string:
	// "host=localhost port=5432 user=postgres password=secret dbname=mydb sslmode=disable"

	if connStr == "" {
		log.Fatal("Connection string is empty  here")
	}

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	// Test connection
	fmt.Println("Testing database connection...")
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")
	return database
}

func InsertToDB(conn *sql.DB, result models.Response) {
	for _, vehicle := range result.Data {
		_, err := conn.Exec(
			`INSERT INTO models (id, make_id, make, name) VALUES ($1, $2, $3, $4)
		     ON CONFLICT (id) DO NOTHING`,
			vehicle.Id, vehicle.MakeId, vehicle.Make, vehicle.Name,
		)
		if err != nil {
			log.Printf("Failed to insert model %d: %v", vehicle.Id, err)
		}
	}
}
