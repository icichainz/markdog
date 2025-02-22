package main

import (
    "embed"
    "io/ioutil"
    "log"
    "mime/multipart"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Arial, sans-serif;
            line-height: 1.6;
            padding: 2em;
            max-width: 800px;
            margin: 0 auto;
        }
        pre {
            background-color: #f5f5f5;
            padding: 1em;
            border-radius: 4px;
        }
        code {
            font-family: 'Courier New', Courier, monospace;
        }
        img {
            max-width: 100%;
        }
    </style>
</head>
<body>
    {{.Content}}
</body>
</html>
`

//go:embed frontend/dist/*
var frontendFiles embed.FS

// readFile reads the content of an uploaded file
func readFile(file *multipart.FileHeader) (string, error) {
    f, err := file.Open()
    if err != nil {
        return "", err
    }
    defer f.Close()

    content, err := ioutil.ReadAll(f)
    if err != nil {
        return "", err
    }

    return string(content), nil
}

func main() {
    if err := DBInitialize(); err != nil {
        log.Fatalf("Failed to initialize the database: %v", err)
    }
    defer DBClose()

    app := fiber.New(fiber.Config{
        AppName: "MarkDog - Markdown Editor",
    })

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())

    // Setup routes
    SetupRoutes(app)

    // Start server
    log.Fatal(app.Listen(":3050"))
}