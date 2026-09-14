package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Todo struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Body      string        `json:"body"`
	Completed bool          `json:"completed"`
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

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected with MongoDB Atlas")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	app.Listen(":" + PORT)
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo

		if err := cursor.Decode(&todo); err != nil {
			return err
		}

		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		log.Fatal(err)
	}

	if todo.Body == "" {
		c.Status(400).JSON(fiber.Map{"error": "body is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		log.Fatal(err)
	}

	todo.ID = insertResult.InsertedID.(bson.ObjectID)

	return c.Status(200).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": "ID not founded"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	return c.Status(200).JSON(fiber.Map{"status": "todo updated"})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectId, err := bson.ObjectIDFromHex(id)

	if err != nil {
		c.Status(400).JSON(fiber.Map{"error": "ID not founded"})
	}

	filter := bson.M{"_id": objectId}

	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID not founded"})
	}

	return c.Status(200).JSON(fiber.Map{"status": "todo deleted"})
}
