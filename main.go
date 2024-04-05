package main

import (
	"fmt"
	mg "github.com/ooyeku/mongoutil/mgutil"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		fmt.Println("MONGO_URI environment variable is not set")
		return
	}
	client, err := mg.NewMongoClient(mongoURI)
	if err != nil {
		// Handle error
	}
	defer client.Close()

	// Read
	var results []bson.M
	err = client.FindAll("Fitness", "exercises", &results)
	if err != nil {
		fmt.Println("Failed to read documents:", err)
	}
	for _, result := range results {
		log.Println(result)
	}
}
