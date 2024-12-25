package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/routes"
	"github.com/amosehiguese/ecommerce-api/server"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T) {
	ta, err := server.SpawnApp()
	router := routes.SetUp(ta.DB, config.Get())
	if err != nil {
		t.Fatalf("failed to spawn app: %v", err)
	}

	// Ensure cleanup after the test
	defer func() {
		// Drop the database after the test completes
		err := server.DropTestDatabase(ta.DB, ta.DB_Name)
		if err != nil {
			t.Errorf("failed to drop test database: %v", err)
		}
	}()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/_healthz", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}
