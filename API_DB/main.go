package main

import (
	"context"
	"log"
	"os"
)


func main() {
	ctx := context.Background()

	db, err := NewDB(ctx)
	if err != nil {
		log.Fatal("Error conect DB:", err)
	}
	log.Println("Conect to CockroachDB")

	baseURL := os.Getenv("API_BASE_URL")
	err = fetchAndStoreData(db, baseURL, "")
	if err != nil {
		log.Fatal("Error Fetch:", err)
	}
	log.Println("Finish Fetch")
}
