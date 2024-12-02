package auth

import (
	"os"
	"testing"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// Set a test JWT secret
	os.Setenv("JWT_SECRET", "test-secret")
	defer os.Unsetenv("JWT_SECRET")

	tests := []struct {
		name    string
		userID  uint
		wantErr bool
	}{
		{
			name:    "Valid user ID",
			userID:  1,
			wantErr: false,
		},
		{
			name:    "Zero user ID",
			userID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Only validate token if we expect no error
				userID, err := ValidateToken(token)
				if err != nil {
					t.Errorf("ValidateToken() error = %v", err)
					return
				}
				if userID != tt.userID {
					t.Errorf("ValidateToken() userID = %v, want %v", userID, tt.userID)
				}
			}
		})
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Invalid format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
