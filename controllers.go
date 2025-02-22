package main

import (
	"log"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"
)

type Controller struct{}

func LoadControllers() *Controller {
	return &Controller{}
}

func (c *Controller) ConvertMarkdown(ctx *fiber.Ctx) error {
	var input struct {
		Markdown string `json:"markdown"`
	}

	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	html := markdown.ToHTML([]byte(input.Markdown), nil, nil)
	return ctx.JSON(fiber.Map{
		"html": string(html),
	})
}

func (c *Controller) HandleFileUpload(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	content, err := readFile(file)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	return ctx.JSON(fiber.Map{
		"content": content,
	})
}

func (c *Controller) ServeStaticFiles(ctx *fiber.Ctx) error {
	path := ctx.Path()
	if path == "/" {
		path = "frontend/dist/index.html"
	} else {
		path = "frontend/dist" + path
	}

	content, err := frontendFiles.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file %s: %v", path, err)
		return ctx.Status(fiber.StatusNotFound).SendString("Not found")
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		ctx.Type("html")
	case ".js":
		ctx.Type("javascript")
	case ".css":
		ctx.Type("css")
	case ".png", ".jpg", ".jpeg", ".gif":
		ctx.Type("image/" + ext[1:])
	}

	return ctx.Send(content)
}
