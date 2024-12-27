package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amosehiguese/ecommerce-api/api/payload"
	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/routes"
	"github.com/amosehiguese/ecommerce-api/server"
	"github.com/stretchr/testify/assert"
)

func TestAuthEndpoints(t *testing.T) {
	// Set up the application
	ta, err := server.SpawnApp()
	if err != nil {
		t.Fatalf("failed to spawn app: %v", err)
	}
	router := routes.SetUp(ta.DB, config.Get())

	// Ensure cleanup after the test
	defer func() {
		err := server.DropTestDatabase(ta.DB, ta.DB_Name)
		if err != nil {
			t.Errorf("failed to drop test database: %v", err)
		}
	}()

	// Test Register endpoint
	t.Run("Test Register", func(t *testing.T) {
		registerPayload := payload.RegisterPayload{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "securepassword123",
		}

		req, err := http.NewRequest("POST", "/api/auth/register", createJSONRequestBody(registerPayload))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Log response body for debugging
		if w.Code != http.StatusOK {
			t.Log("Response Body:", w.Body.String())
		}

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "User created successfully")
	})

	// Test Login endpoint
	t.Run("Test Login", func(t *testing.T) {
		loginPayload := payload.LoginPayload{
			Email:    "john.doe@example.com",
			Password: "securepassword123",
		}

		req, err := http.NewRequest("POST", "/api/auth/login", createJSONRequestBody(loginPayload))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "access")
	})

}

func createJSONRequestBody(v interface{}) *bytes.Buffer {
	body, _ := json.Marshal(v)
	return bytes.NewBuffer(body)
}
