package main

import (
	"bytes"
	"embed"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"path/filepath"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gomarkdown/markdown"
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
	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the content
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: "MarkDog - Markdown Editor",
		
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// API routes
	api := app.Group("/api")

	// Convert markdown to PDF
	api.Post("/convert", func(c *fiber.Ctx) error {
		var input struct {
			Markdown string `json:"markdown"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		html := markdown.ToHTML([]byte(input.Markdown), nil, nil)
		return c.JSON(fiber.Map{
			"html": string(html),
			})
	})

	// Convert markdown to PDF
	api.Post("/export-pdf", func(c *fiber.Ctx) error {
		var input struct {
			Markdown string `json:"markdown"`
		}

		if err := c.BodyParser(&input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid input",
			})
		}

		// Convert markdown to HTML
		html := markdown.ToHTML([]byte(input.Markdown), nil, nil)

		// Create PDF generator
		pdfg, err := wkhtmltopdf.NewPDFGenerator()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create PDF generator",
			})
		}

		// Create a new template and parse the HTML
		tmpl, err := template.New("pdf").Parse(htmlTemplate)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse template",
			})
		}

		// Execute the template with the HTML content
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, struct{ Content template.HTML }{Content: template.HTML(html)})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to execute template",
			})
		}

		// Add HTML page to PDF generator
		page := wkhtmltopdf.NewPageReader(bytes.NewReader(buf.Bytes()))
		page.EnableLocalFileAccess.Set(true)
		pdfg.AddPage(page)

		// Generate PDF
		err = pdfg.Create()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate PDF",
			})
		}

		// Set response headers
		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", "attachment; filename=document.pdf")

		return c.Send(pdfg.Bytes())
	})

	// Handle file upload
	api.Post("/upload", func(c *fiber.Ctx) error {
		// Get uploaded file
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No file uploaded",
			})
		}

		// Read file content
		content, err := readFile(file)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to read file",
			})
		}

		return c.JSON(fiber.Map{
			"content": content,
		})
	})

	// Serve embedded static files
	app.Use("/", func(c *fiber.Ctx) error {
		path := c.Path()
		if path == "/" {
			path = "frontend/dist/index.html"
		} else {
			// Remove leading slash and add frontend/dist prefix
			path = "frontend/dist" + path
		}

		content, err := frontendFiles.ReadFile(path)
		if err != nil {
			log.Printf("Failed to read file %s: %v", path, err)
			return c.Status(fiber.StatusNotFound).SendString("Not found")
		}

		// Set content type based on file extension
		ext := filepath.Ext(path)
		switch ext {
		case ".html":
			c.Type("html")
		case ".js":
			c.Type("javascript")
		case ".css":
			c.Type("css")
		case ".png", ".jpg", ".jpeg", ".gif":
			c.Type("image/" + ext[1:])
		}

		return c.Send(content)
	})

	// Start server
	log.Fatal(app.Listen(":3050"))
}
