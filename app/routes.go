package app

import "github.com/gofiber/fiber/v2"

func SetupRoutes(app *fiber.App) {
    controller := LoadControllers()

    // API routes
    api := app.Group("/api")
    api.Post("/convert", controller.ConvertMarkdown)
    api.Post("/upload", controller.HandleFileUpload)
}