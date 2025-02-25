package app

import (
	"log"
	"path/filepath"
	"embed"
	"github.com/gofiber/fiber/v2"
	"strings"

)

func SetupStaticRoute(app *fiber.App, frontendFiles embed.FS) {
	
	app.Use(
		func(ctx *fiber.Ctx) error {

			if  strings.HasPrefix(ctx.Path(),"/api"){
				return ctx.Next()
			}
			return ServeStaticFiles(ctx, frontendFiles)
		},
	)
	

}

func  ServeStaticFiles(ctx *fiber.Ctx ,
	frontendFiles embed.FS,
	) error {
	requestPath := ctx.Path()
	

	if requestPath == "/" {
		requestPath = "frontend/dist/index.html"
	} else {
		requestPath = "frontend/dist" + requestPath
	}
    

	content, err := frontendFiles.ReadFile(requestPath)

	if err != nil {
		log.Printf("Failed to read file %s: %v", requestPath, err)
		return ctx.Status(fiber.StatusNotFound).SendString("Not found")
	}

	ext := filepath.Ext(requestPath)
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
