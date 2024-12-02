package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"todo-api/internal/database"
	"todo-api/internal/models"
)

type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
}

// @Summary Create a new todo
// @Description Create a new todo for the authenticated user
// @Tags todos
// @Accept json
// @Produce json
// @Security Bearer
// @Param todo body models.Todo true "Todo object"
// @Success 201 {object} models.Todo
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/todos [post]
func CreateTodo(c *gin.Context) {
	userID, _ := c.Get("userID")

	var todo models.Todo
	if err := c.BindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	todo.UserID = userID.(uint)

	result := database.GetDB().Create(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: result.Error.Error()})
		return
	}

	database.GetDB().Preload("User").First(&todo, todo.ID)

	c.JSON(http.StatusCreated, todo)
}

// @Summary Get all todos
// @Description Get all todos for the authenticated user
// @Tags todos
// @Produce json
// @Security Bearer
// @Success 200 {array} models.Todo
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/todos [get]
func GetTodos(c *gin.Context) {
	userID, _ := c.Get("userID")

	var todos []models.Todo
	result := database.GetDB().Preload("User").Where("user_id = ?", userID).Find(&todos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

// @Summary Get a todo
// @Description Get a specific todo by ID
// @Tags todos
// @Produce json
// @Security Bearer
// @Param id path int true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/todos/{id} [get]
func GetTodo(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	var todo models.Todo

	result := database.GetDB().Preload("User").Where("id = ? AND user_id = ?", id, userID).First(&todo)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// @Summary Update a todo
// @Description Update a specific todo by ID
// @Tags todos
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Todo ID"
// @Param todo body models.Todo true "Todo object"
// @Success 200 {object} models.Todo
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/todos/{id} [put]
func UpdateTodo(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	var todo models.Todo

	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo not found"})
		return
	}

	var updateData models.Todo
	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	result := database.GetDB().Model(&todo).Updates(map[string]interface{}{
		"title":       updateData.Title,
		"description": updateData.Description,
		"completed":   updateData.Completed,
	})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: result.Error.Error()})
		return
	}

	database.GetDB().Preload("User").First(&todo, id)

	c.JSON(http.StatusOK, todo)
}

// @Summary Delete a todo
// @Description Delete a specific todo by ID
// @Tags todos
// @Security Bearer
// @Param id path int true "Todo ID"
// @Success 204 "No Content"
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/todos/{id} [delete]
func DeleteTodo(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	var todo models.Todo

	if err := database.GetDB().Where("id = ? AND user_id = ?", id, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Todo not found"})
		return
	}

	result := database.GetDB().Delete(&todo)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
