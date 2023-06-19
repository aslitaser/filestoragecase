package handlers

import (
	"encoding/json"
	"filestorageapi/db"
	"filestorageapi/models"
	"io"
	"net/http"
	"os"
	"strconv"
)

type UploadFileHandler struct {
	DBManager *db.DBManager
}

func NewUploadFileHandler(dbManager *db.DBManager) *UploadFileHandler {
	return &UploadFileHandler{
		DBManager: dbManager,
	}
}

func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileSize, err := strconv.ParseInt(r.FormValue("filesize"), 10, 64)
	if err != nil || fileSize <= 0 {
		http.Error(w, "Invalid filesize", http.StatusBadRequest)
		return
	}

	filePath := "/file/path" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}

	fileMetadata := &models.FileMetadata{
		Filename: header.Filename,
		Filesize: fileSize,
	}

	if err := h.DBManager.SaveFileMetadata(fileMetadata); err != nil {
		http.Error(w, "Error saving file metadata", http.StatusInternalServerError)
		return
	}

	// Set response header content type to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Create a JSON response with the file metadata
	jsonResponse, err := json.Marshal(map[string]interface{}{
		"message":      "File uploaded successfully",
		"fileMetadata": fileMetadata,
	})

	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	// Write the JSON response to the response writer
	w.Write(jsonResponse)
}
