package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("Running app")

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "ok"})
	})

	log.Fatal(app.Listen(":8080"))
}
