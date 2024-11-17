package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// HealthController handles health check endpoints
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Register registers health check routes
func (c *HealthController) Register(app *fiber.App) {
	app.Get("/health", c.HealthCheck)
}

// HealthCheck handles the health check endpoint
// @Summary Health check endpoint
// @Description Check if the service is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (c *HealthController) HealthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status": "healthy",
		"message": "Service is up and running",
	})
}
