package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"todo-api/internal/middleware"
	"todo-api/internal/models"
	"todo-api/internal/test"
)

func TestCreateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test user
	db, _ := test.SetupTestDB()
	testUser, err := test.CreateTestUser(db)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		reqBody    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid todo creation",
			reqBody: map[string]interface{}{
				"title":       "Test Todo",
				"description": "Test Description",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing title",
			reqBody: map[string]interface{}{
				"description": "Test Description",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test case
			router := gin.New()

			// Clear test data before each test
			test.ClearTestData(db)

			router.POST("/todos", func(c *gin.Context) {
				c.Set("userID", testUser.ID)
				CreateTodo(c)
			})

			jsonBody, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			
			if w.Code == http.StatusCreated {
				var response models.Todo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.reqBody["title"], response.Title)
				assert.Equal(t, testUser.ID, response.UserID)
			}
		})
	}
}

func TestGetTodos(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test user and some todos
	db, _ := test.SetupTestDB()
	testUser, err := test.CreateTestUser(db)
	assert.NoError(t, err)

	// Create test todos
	todo1 := &models.Todo{
		Title:       "Test Todo 1",
		Description: "Test Description 1",
		UserID:      testUser.ID,
	}
	todo2 := &models.Todo{
		Title:       "Test Todo 2",
		Description: "Test Description 2",
		UserID:      testUser.ID,
	}
	db.Create(todo1)
	db.Create(todo2)

	tests := []struct {
		name       string
		setupAuth  bool
		wantStatus int
		wantCount  int
	}{
		{
			name:       "Get user todos",
			setupAuth:  true,
			wantStatus: http.StatusOK,
			wantCount:  2,
		},
		{
			name:       "Unauthorized request",
			setupAuth:  false,
			wantStatus: http.StatusUnauthorized,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test case
			router := gin.New()
			
			// Clear test data before each test
			test.ClearTestData(db)
			// Recreate test todos
			db.Create(todo1)
			db.Create(todo2)

			if tt.setupAuth {
				router.GET("/todos", func(c *gin.Context) {
					c.Set("userID", testUser.ID)
					GetTodos(c)
				})
			} else {
				router.GET("/todos", middleware.AuthMiddleware(), GetTodos)
			}

			req, _ := http.NewRequest("GET", "/todos", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			
			if w.Code == http.StatusOK {
				var response []models.Todo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantCount, len(response))
			}
		})
	}
}

func TestUpdateTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test user and todo
	db, _ := test.SetupTestDB()
	testUser, err := test.CreateTestUser(db)
	assert.NoError(t, err)

	todo := &models.Todo{
		Title:       "Original Todo",
		Description: "Original Description",
		UserID:      testUser.ID,
	}
	result := db.Create(todo)
	assert.NoError(t, result.Error)

	tests := []struct {
		name       string
		todoID     string
		reqBody    map[string]interface{}
		setupAuth  bool
		wantStatus int
	}{
		{
			name:   "Valid update",
			todoID: "1",
			reqBody: map[string]interface{}{
				"title":       "Updated Todo",
				"description": "Updated Description",
			},
			setupAuth:  true,
			wantStatus: http.StatusOK,
		},
		{
			name:   "Todo not found",
			todoID: "999",
			reqBody: map[string]interface{}{
				"title":       "Updated Todo",
				"description": "Updated Description",
			},
			setupAuth:  true,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test case
			router := gin.New()

			// Clear test data before each test
			test.ClearTestData(db)
			// Recreate test todo
			result := db.Create(todo)
			assert.NoError(t, result.Error)

			if tt.setupAuth {
				router.PUT("/todos/:id", func(c *gin.Context) {
					c.Set("userID", testUser.ID)
					UpdateTodo(c)
				})
			} else {
				router.PUT("/todos/:id", UpdateTodo)
			}

			jsonBody, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("PUT", "/todos/"+tt.todoID, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if w.Code == http.StatusOK {
				var response models.Todo
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.reqBody["title"], response.Title)
				assert.Equal(t, tt.reqBody["description"], response.Description)
				assert.Equal(t, testUser.ID, response.UserID)
			}
		})
	}
}

func TestDeleteTodo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a test user and todo
	db, _ := test.SetupTestDB()
	testUser, err := test.CreateTestUser(db)
	assert.NoError(t, err)

	todo := &models.Todo{
		Title:       "Test Todo",
		Description: "Test Description",
		UserID:      testUser.ID,
	}
	result := db.Create(todo)
	assert.NoError(t, result.Error)

	tests := []struct {
		name       string
		todoID     string
		setupAuth  bool
		wantStatus int
	}{
		{
			name:       "Valid deletion",
			todoID:     "1",
			setupAuth:  true,
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "Todo not found",
			todoID:     "999",
			setupAuth:  true,
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new router for each test case
			router := gin.New()

			// Clear test data before each test
			test.ClearTestData(db)
			// Recreate test todo
			result := db.Create(todo)
			assert.NoError(t, result.Error)

			if tt.setupAuth {
				router.DELETE("/todos/:id", func(c *gin.Context) {
					c.Set("userID", testUser.ID)
					DeleteTodo(c)
				})
			} else {
				router.DELETE("/todos/:id", DeleteTodo)
			}

			req, _ := http.NewRequest("DELETE", "/todos/"+tt.todoID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if w.Code == http.StatusNoContent {
				// Verify todo was actually deleted
				var count int64
				db.Model(&models.Todo{}).Where("id = ?", tt.todoID).Count(&count)
				assert.Equal(t, int64(0), count)
			}
		})
	}
}
