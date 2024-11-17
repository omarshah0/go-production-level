package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yourusername/go-production-level/config"
	"github.com/yourusername/go-production-level/internal/controllers"
	"github.com/yourusername/go-production-level/internal/middlewares"
	"github.com/yourusername/go-production-level/internal/models"
	"github.com/yourusername/go-production-level/internal/repository"
	"github.com/yourusername/go-production-level/internal/services"
	"github.com/yourusername/go-production-level/internal/utils"
)

// @title Go Production Level API
// @version 1.0
// @description A production-ready RESTful API built with Go
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := utils.InitDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate database
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize Redis
	redis, err := utils.InitRedis(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize repositories
	gormRepo := repository.NewGormRepository(db)
	userRepo := repository.NewUserRepository(gormRepo)

	// Initialize services
	userService := services.NewUserService(userRepo, redis, cfg)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	healthController := controllers.NewHealthController()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Serve Swagger documentation
	app.Static("/api/v1/docs", "./docs")

	// Swagger documentation UI
	app.Get("/swagger/*", middlewares.SwaggerMiddleware())

	// Register routes
	api := app.Group("/api/v1")

	// Health check route (before other routes)
	healthController.Register(app)

	// Public routes
	userController.Register(app)

	// Protected routes
	protected := api.Group("/protected")
	protected.Use(middlewares.AuthMiddleware(cfg))

	// Admin routes
	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	log.Printf("Swagger documentation available at http://localhost:%s/swagger/", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
