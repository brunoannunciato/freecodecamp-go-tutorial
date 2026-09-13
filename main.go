package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json: "id"`
	Completed bool   `json: "completed"`
	Body      string `json: "body"`
}

func main() {
	fmt.Println("Running app")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Body field is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(200).JSON(*todo)
	})

	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if id == fmt.Sprint(todo.ID) {
				todos[i].Completed = true

				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Todo with ID %v not founded", id)})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)

				return c.Status(200).JSON(fiber.Map{"status": "todo deleted successfully!"})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": fmt.Sprintf("Todo with id %v not exists", id)})
	})

	log.Fatal(app.Listen(":" + PORT))
}
