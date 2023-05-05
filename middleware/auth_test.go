package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Setup
	r := gin.New()
	r.Use(AuthMiddleware())

	// Test cases
	testCases := []struct {
		name          string
		token         string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "Valid token",
			token:         "valid_token",
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		{
			name:          "Missing token",
			token:         "",
			expectedCode:  http.StatusUnauthorized,
			expectedError: "unauthorized",
		},
		{
			name:          "Invalid token",
			token:         "invalid_token",
			expectedCode:  http.StatusUnauthorized,
			expectedError: "Unauthorized",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new request with Authorization header
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Authorization", tc.token)

			// Perform the request
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Check response code and error message
			assert.Equal(t, tc.expectedCode, w.Code)
			if tc.expectedError != "" {
				assert.Contains(t, w.Body.String(), tc.expectedError)
			}
		})
	}
}
