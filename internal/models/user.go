package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// User represents the user model in the database
// @Description User account information
type User struct {
	ID        uint           `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email" validate:"required,email" example:"user@example.com"`
	Password  string         `json:"password,omitempty" validate:"required,min=6" example:"password123"`
	Name      string         `json:"name" validate:"required" example:"John Doe"`
	Role      string         `json:"role" validate:"required,oneof=admin user" example:"user"`
}

// UserResponse represents the user response without sensitive information
// @Description User information for API responses
type UserResponse struct {
	ID        uint      `json:"id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	Email     string    `json:"email" example:"user@example.com"`
	Name      string    `json:"name" example:"John Doe"`
	Role      string    `json:"role" example:"user"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Validate validates the user model and returns an array of validation errors
func (u *User) Validate() []ValidationError {
	validate := validator.New()
	err := validate.Struct(u)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		var element ValidationError
		element.Field = err.Field()
		element.Error = getErrorMsg(err)
		errors = append(errors, element)
	}
	return errors
}

// getErrorMsg returns a human-readable error message for validation errors
func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Should be at least " + err.Param() + " characters long"
	case "oneof":
		return "Should be one of: " + err.Param()
	}
	return "Unknown validation error"
}
