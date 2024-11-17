package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/yourusername/go-production-level/internal/models"
	"github.com/yourusername/go-production-level/internal/services"
)

// UserController handles HTTP requests for users
type UserController struct {
	userService services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// LoginRequest represents the login request body
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register registers all user routes
func (c *UserController) Register(app *fiber.App) {
	api := app.Group("/api/v1")
	
	// Public routes
	api.Post("/login", c.Login)
	api.Post("/users", c.CreateUser)

	// Protected routes
	users := api.Group("/users")
	users.Get("/", c.ListUsers)
	users.Get("/:id", c.GetUser)
	users.Put("/:id", c.UpdateUser)
	users.Delete("/:id", c.DeleteUser)
}

// Login handles user authentication
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]
func (c *UserController) Login(ctx *fiber.Ctx) error {
	var req LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	token, err := c.userService.Login(ctx.Context(), req.Email, req.Password)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid credentials",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
	})
}

// CreateUser handles user creation
// @Summary Create new user
// @Description Create a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.User true "User object"
// @Success 201 {object} map[string]string
// @Failure 400 {array} models.ValidationError
// @Router /users [post]
func (c *UserController) CreateUser(ctx *fiber.Ctx) error {
	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate user input
	if errors := user.Validate(); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	if err := c.userService.Create(ctx.Context(), &user); err != nil {
		if err == services.ErrEmailExists {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create user",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
	})
}

// GetUser handles fetching a single user
// @Summary Get user by ID
// @Description Get user details by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	user, err := c.userService.GetByID(ctx.Context(), uint(id))
	if err != nil {
		if err == services.ErrUserNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(user)
}

// UpdateUser handles user updates
// @Summary Update user
// @Description Update user details
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User object"
// @Success 200 {object} map[string]string
// @Failure 400 {array} models.ValidationError
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate user input
	if errors := user.Validate(); errors != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errors)
	}

	user.ID = uint(id)
	if err := c.userService.Update(ctx.Context(), &user); err != nil {
		if err == services.ErrUserNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "user updated successfully",
	})
}

// DeleteUser handles user deletion
// @Summary Delete user
// @Description Delete user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	if err := c.userService.Delete(ctx.Context(), uint(id)); err != nil {
		if err == services.ErrUserNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "user not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"message": "user deleted successfully",
	})
}

// ListUsers handles fetching a list of users
// @Summary List users
// @Description Get paginated list of users
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /users [get]
func (c *UserController) ListUsers(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := c.userService.List(ctx.Context(), offset, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return ctx.JSON(fiber.Map{
		"users": users,
		"page":  page,
		"limit": limit,
	})
}
