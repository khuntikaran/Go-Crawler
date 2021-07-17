package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Database {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	MongoDb := os.Getenv("MONGOURL")
	// Set client options
	clientOptions := options.Client().ApplyURI(MongoDb)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		fmt.Println("error aavo")
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	db := client.Database("movies")

	return db
}
