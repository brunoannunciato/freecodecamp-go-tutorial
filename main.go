package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Todo struct {
	_ID       int    `json: "id" bson: "_id"`
	Body      string `json: "body"`
	Completed bool   `json: "completed"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("App running")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected with MongoDB Atlas")
}
