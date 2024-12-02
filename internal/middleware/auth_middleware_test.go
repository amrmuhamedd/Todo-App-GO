package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo-api/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Set up test environment
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		setupAuth  func() string
		wantStatus int
	}{
		{
			name: "Valid token",
			setupAuth: func() string {
				token, _ := auth.GenerateToken(1)
				return token
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Missing token",
			setupAuth: func() string {
				return ""
			},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			setupAuth: func() string {
				return "invalid-token"
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AuthMiddleware())
			
			// Add a test endpoint
			router.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			
			token := tt.setupAuth()
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}

			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
		want       string
		wantErr    bool
	}{
		{
			name:       "Valid bearer token",
			authHeader: "Bearer valid-token",
			want:       "valid-token",
			wantErr:    false,
		},
		{
			name:       "Empty header",
			authHeader: "",
			want:       "",
			wantErr:    true,
		},
		{
			name:       "Missing bearer prefix",
			authHeader: "valid-token",
			want:       "",
			wantErr:    true,
		},
		{
			name:       "Only bearer prefix",
			authHeader: "Bearer",
			want:       "",
			wantErr:    true,
		},
		{
			name:       "Wrong prefix",
			authHeader: "Basic valid-token",
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractBearerToken(tt.authHeader)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
