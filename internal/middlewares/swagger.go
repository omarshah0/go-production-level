package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SwaggerMiddleware serves Swagger UI
func SwaggerMiddleware() fiber.Handler {
	return swagger.New(swagger.Config{
		Title: "Go Production Level API",
		URL:   "/api/v1/docs/swagger.yaml",
	})
}
