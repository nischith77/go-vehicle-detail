package main

import (
	"encoding/json"
	"fmt"
	db "go-vehicle-detail/database"
	"go-vehicle-detail/models"
	"log"
	"net/http"
	"os"
)

func main() {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		log.Fatal("DB_CONN environment variable is required")
	}

	fmt.Printf("Attempting to connect with: %s\n", connStr)
	conn := db.ConnectDB(connStr)
	defer conn.Close()
	resp, err := http.Get("https://carapi.app/api/models/v2")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var result models.Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	db.InsertToDB(conn, result)
	if err != nil {
		panic(err)
	}
	next := result.Collection.Next
	for {
		url := "https://carapi.app" + next
		resp, err := http.Get(url)

		if err != nil {
			panic(err)
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close() // Close immediately after reading
		next = result.Collection.Next
		db.InsertToDB(conn, result)
		if err != nil {
			panic(err)
		}
		if next == "" {
			break
		}

	}

}
