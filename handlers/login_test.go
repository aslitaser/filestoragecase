package handlers

import (
	"filestorageapi/db"
	"filestorageapi/jwtmanager"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	jwtSecret := os.Getenv("JWT_SECRET")
	dbManager, err := db.NewDBManager("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	tokenDuration := time.Hour
	jwtManager := jwtmanager.NewJWTManager(jwtSecret, tokenDuration)

	loginHandler := NewLoginHandler(dbManager, jwtManager)

	// Authenticate a user and return a JWT.
	t.Run("Authenticate a user and return a JWT", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"testuser", "password":"testpassword"}`))
		resp := httptest.NewRecorder()
		loginHandler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	// Test for invalid username or password.
	t.Run("Test for invalid username or password", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"wronguser", "password":"wrongpassword"}`))
		resp := httptest.NewRecorder()
		loginHandler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusUnauthorized {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusUnauthorized)
		}
	})

	// Test for successful login and JWT generation.
	t.Run("Test for successful login and JWT generation", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username":"testuser", "password":"testpassword"}`))
		resp := httptest.NewRecorder()
		loginHandler.ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// Check that the Authorization header was set correctly.
		authHeader := resp.Header().Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			t.Errorf("handler returned wrong Authorization header: got %v want prefix %v",
				authHeader, "Bearer ")
		}
	})
}
