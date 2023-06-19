package handlers

import (
	"filestorageapi/db"
	"net/http"
	"os"
)

type DeleteFileHandler struct {
	DBManager *db.DBManager
}

func NewDeleteFileHandler(dbManager *db.DBManager) *DeleteFileHandler {
	return &DeleteFileHandler{
		DBManager: dbManager,
	}
}

func (h *DeleteFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	fileMetadata, err := h.DBManager.GetFileMetadata(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	filePath := "/path/to/store/files/" + fileMetadata.Filename
	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Error deleting the file", http.StatusInternalServerError)
		return
	}

	if err := h.DBManager.DeleteFileMetadata(fileID); err != nil {
		http.Error(w, "Error deleting file metadata", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
