package main

import (
	"embed"
	"log"
	Application "markdog/app"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed frontend/dist/*
var frontendFiles embed.FS



func main() {
    if err := Application.DBInitialize(); err != nil {
        log.Fatalf("Failed to initialize the database: %v", err)
    }
    defer Application.DBClose()

    app := fiber.New(fiber.Config{
        AppName: "MarkDog - Markdown Editor",
    })

    app.Server().ReadTimeout= 60 * time.Second 
    app.Server().WriteTimeout= 60 * time.Second 

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())


    // Setup routes
    Application.SetupRoutes(app)
    // Setup static routes
    Application.SetupStaticRoute(app, frontendFiles)
    
    // Start server
    log.Fatal(app.Listen(":3050"))
}