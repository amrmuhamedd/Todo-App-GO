package test

import (
	"todo-api/internal/models"
	"todo-api/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schemas
	err = db.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		return nil, err
	}

	// Set the global DB instance
	database.DB = db

	return db, nil
}

// ClearTestData cleans up test data from the database
func ClearTestData(db *gorm.DB) error {
	err := db.Exec("DELETE FROM todos").Error
	if err != nil {
		return err
	}

	err = db.Exec("DELETE FROM users").Error
	if err != nil {
		return err
	}

	return nil
}

// CreateTestUser creates a test user in the database
func CreateTestUser(db *gorm.DB) (*models.User, error) {
	user := &models.User{
		Email:    "test@example.com",
		Password: "testpassword",
	}

	// Hash the password before saving
	if err := user.HashPassword(); err != nil {
		return nil, err
	}

	result := db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// CreateTestTodo creates a test todo item in the database
func CreateTestTodo(db *gorm.DB, userID uint) (*models.Todo, error) {
	todo := &models.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		UserID:      userID,
		Completed:   false,
	}

	result := db.Create(todo)
	if result.Error != nil {
		return nil, result.Error
	}

	return todo, nil
}

// SetupTestRouter creates a new Gin router in test mode
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// MockAuthMiddleware creates a mock authentication middleware for testing
func MockAuthMiddleware(userID uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	}
}
