package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type Todo struct {
	_ID       int    `json: "id" bson: "_id"`
	Body      string `json: "body"`
	Completed bool   `json: "completed"`
}

func main() {
	fmt.Println("App running")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}
}
