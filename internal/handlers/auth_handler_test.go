package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-api/internal/test"

	"github.com/stretchr/testify/assert"
)

func TestSignup(t *testing.T) {
	router := setupTestRouter()
	router.POST("/auth/signup", Signup)

	tests := []struct {
		name       string
		reqBody    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid signup",
			reqBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "password123",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Missing required fields",
			reqBody: map[string]interface{}{
				"email": "test@example.com",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear test data before each test
			db, _ := test.SetupTestDB()
			test.ClearTestData(db)

			jsonBody, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestLogin(t *testing.T) {
	router := setupTestRouter()
	router.POST("/auth/login", Login)

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
			name: "Valid credentials",
			reqBody: map[string]interface{}{
				"email":    testUser.Email,
				"password": "testpassword",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Missing credentials",
			reqBody: map[string]interface{}{
				"email": testUser.Email,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid credentials",
			reqBody: map[string]interface{}{
				"email":    testUser.Email,
				"password": "wrongpassword",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.reqBody)
			req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			if w.Code == http.StatusOK {
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "token")
			}
		})
	}
}
