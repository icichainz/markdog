package app

import (	
	"github.com/gofiber/fiber/v2"
	"github.com/gomarkdown/markdown"

)





type Controller struct{
}

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

