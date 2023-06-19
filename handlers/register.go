package handlers

import (
	"encoding/json"
	"filestorageapi/db"
	"filestorageapi/models"
	"net/http"
)

type RegisterHandler struct {
	DBManager *db.DBManager
}

func NewRegisterHandler(dbManager *db.DBManager) *RegisterHandler {
	return &RegisterHandler{
		DBManager: dbManager,
	}
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Missing username or password", http.StatusBadRequest)
		return
	}

	existingUser, _ := h.DBManager.GetUserByUsername(user.Username)
	if existingUser != nil {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	if err := h.DBManager.CreateUser(&user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
