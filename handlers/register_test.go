package handlers_test

import (
	"bytes"
	"encoding/json"
	"filestorageapi/db"
	"filestorageapi/handlers"
	"filestorageapi/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegister(t *testing.T) {
	// Initialize the real DBManager
	dbManager, err := db.NewDBManager("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	registerHandler := handlers.NewRegisterHandler(dbManager)

	t.Run("Create a new user account", func(t *testing.T) {
		newUser := models.User{Username: "testuser", Password: "testpassword"}
		userJSON, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		registerHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})

	t.Run("Test for missing username or password query parameters", func(t *testing.T) {
		newUser := models.User{Username: "", Password: "testpassword"}
		userJSON, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		registerHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("Test for existing user", func(t *testing.T) {
		existingUser := models.User{Username: "existinguser", Password: "existingpassword"}
		_ = dbManager.CreateUser(&existingUser)

		userJSON, _ := json.Marshal(existingUser)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		registerHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}
	})

	t.Run("Test for successful registration", func(t *testing.T) {
		newUser := models.User{Username: "newuser", Password: "password123"}
		userJSON, _ := json.Marshal(newUser)

		req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(userJSON))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		registerHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}
	})
}
