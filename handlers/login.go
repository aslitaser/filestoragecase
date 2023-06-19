package handlers

import (
	"encoding/json"
	"filestorageapi/db"
	"filestorageapi/jwtmanager"
	"filestorageapi/models"
	"net/http"
)

type LoginHandler struct {
	DBManager  *db.DBManager
	JWTManager *jwtmanager.JWTManager
}

func NewLoginHandler(dbManager *db.DBManager, jwtManager *jwtmanager.JWTManager) *LoginHandler {
	return &LoginHandler{
		DBManager:  dbManager,
		JWTManager: jwtManager,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.DBManager.GetUserByUsername(credentials.Username)
	if err != nil || user.Password != credentials.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.JWTManager.GenerateToken(user.ID.Hex())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set response header content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Create a JSON response with the token
	jsonResponse, err := json.Marshal(map[string]string{
		"token": token,
	})

	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the response writer
	w.Write(jsonResponse)
}
