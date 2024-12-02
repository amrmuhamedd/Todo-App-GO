package routes

import (
	"todo-api/internal/handlers"
	"todo-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	auth := r.Group("/api/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		todos := api.Group("/todos")
		{
			todos.POST("", handlers.CreateTodo)
			todos.GET("", handlers.GetTodos)
			todos.GET("/:id", handlers.GetTodo)
			todos.PUT("/:id", handlers.UpdateTodo)
			todos.DELETE("/:id", handlers.DeleteTodo)
		}
	}
}
