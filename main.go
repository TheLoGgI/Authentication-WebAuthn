package main

import (
	"log"

	"github.com/TheLoGgI/database"
	"github.com/TheLoGgI/models"
	"github.com/TheLoGgI/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

const Port string = "3000"

func createServer() models.Server {
	app := fiber.New()
	database := database.GetMongoDatabase()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3001",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	server := models.Server{
		Database: database,
		App:      app,
	}

	return server
}

func main() {
	// Create server
	server := createServer()

	// Routes
	server.App.Static("/", "./static")
	routes.Providers(server)
	routes.Users(server)

	// Listen for port
	log.Printf("Starting server at port " + Port + "\n")
	log.Fatal(server.App.Listen("127.0.0.1:" + Port))
}
