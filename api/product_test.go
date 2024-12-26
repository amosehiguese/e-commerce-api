package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/amosehiguese/ecommerce-api/pkg/config"
	"github.com/amosehiguese/ecommerce-api/routes"
	"github.com/amosehiguese/ecommerce-api/server"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductEndpoint(t *testing.T) {
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

	product := map[string]interface{}{
		"name":        "Test Product",
		"description": "This is a test product",
		"price":       100.50,
	}

	productJSON, _ := json.Marshal(product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Product created successfully", response["msg"])
}
func TestGetProductEndpoint(t *testing.T) {
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

	// Create a sample product first
	product := map[string]interface{}{
		"name":        "Test Product",
		"description": "This is a test product",
		"price":       100.50,
	}

	productJSON, _ := json.Marshal(product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Check for response body status
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 or 200, got %v. Response: %s", w.Code, w.Body.String())
	}

	// Extract product ID from response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response: %v, Response Body: %s", err, w.Body.String())
	}

	productID, ok := response["product"].(string)
	if !ok || productID == "" {
		t.Fatalf("expected product ID in response, but got: %v", response["product"])
	}

	// Now test getting the product by ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/"+productID, nil)
	router.ServeHTTP(w, req)

	// Check for unexpected response codes
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %v. Response: %s", w.Code, w.Body.String())
	}

	// Verify the response contains the product data
	var productResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &productResponse)
	if err != nil {
		t.Fatalf("error unmarshaling product response: %v, Response Body: %s", err, w.Body.String())
	}

	// Check if the product name is correct
	if productResponse["name"] != "Test Product" {
		t.Fatalf("expected product name 'Test Product', but got: %v", productResponse["name"])
	}
}
func TestUpdateProductEndpoint(t *testing.T) {
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

	// Create a sample product first
	product := map[string]interface{}{
		"name":        "Test Product",
		"description": "This is a test product",
		"price":       100.50,
	}

	productJSON, _ := json.Marshal(product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Check for response status
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 or 200, got %v. Response: %s", w.Code, w.Body.String())
	}

	// Extract product ID from response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response: %v, Response Body: %s", err, w.Body.String())
	}

	productID, ok := response["product"].(string)
	if !ok || productID == "" {
		t.Fatalf("expected product ID in response, but got: %v", response["product"])
	}

	// Update the product
	updatedProduct := map[string]interface{}{
		"name":        "Updated Test Product",
		"description": "This is an updated test product",
		"price":       150.75,
	}

	updatedProductJSON, _ := json.Marshal(updatedProduct)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/products/"+productID, bytes.NewBuffer(updatedProductJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Check for response status
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %v. Response: %s", w.Code, w.Body.String())
	}

	// Verify the response contains the updated product data
	var updatedProductResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &updatedProductResponse)
	if err != nil {
		t.Fatalf("error unmarshaling updated product response: %v, Response Body: %s", err, w.Body.String())
	}

	// Check if the updated product name is correct
	if updatedProductResponse["name"] != "Updated Test Product" {
		t.Fatalf("expected product name 'Updated Test Product', but got: %v", updatedProductResponse["name"])
	}

	// Check if the updated product description is correct
	if updatedProductResponse["description"] != "This is an updated test product" {
		t.Fatalf("expected product description 'This is an updated test product', but got: %v", updatedProductResponse["description"])
	}

	// Check if the updated product price is correct
	if updatedProductResponse["price"] != 150.75 {
		t.Fatalf("expected product price 150.75, but got: %v", updatedProductResponse["price"])
	}
}
func TestDeleteProductEndpoint(t *testing.T) {
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

	// Create a sample product first
	product := map[string]interface{}{
		"name":        "Test Product",
		"description": "This is a test product",
		"price":       100.50,
	}

	productJSON, _ := json.Marshal(product)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/products", bytes.NewBuffer(productJSON))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	// Check for response status
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 or 200, got %v. Response: %s", w.Code, w.Body.String())
	}

	// Extract product ID from response
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("error unmarshaling response: %v, Response Body: %s", err, w.Body.String())
	}

	productID, ok := response["product"].(string)
	if !ok || productID == "" {
		t.Fatalf("expected product ID in response, but got: %v", response["product"])
	}

	// Delete the product
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/products/"+productID, nil)

	router.ServeHTTP(w, req)

	// Check for response status
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %v. Response: %s", w.Code, w.Body.String())
	}

	// Verify that the product is deleted by checking the response body
	var deleteResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &deleteResponse)
	if err != nil {
		t.Fatalf("error unmarshaling delete product response: %v, Response Body: %s", err, w.Body.String())
	}

	// Check if the response indicates successful deletion
	if deleteResponse["message"] != "Product deleted successfully" {
		t.Fatalf("expected 'Product deleted successfully' message, but got: %v", deleteResponse["message"])
	}

	// Verify that the product no longer exists by trying to fetch the deleted product
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/products/"+productID, nil)

	router.ServeHTTP(w, req)

	// Check for response status - should be 404 if the product doesn't exist
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for deleted product, got %v. Response: %s", w.Code, w.Body.String())
	}
}
