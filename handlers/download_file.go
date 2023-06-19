package handlers

import (
	"filestorageapi/db"
	"io"
	"net/http"
	"os"
	"strconv"
)

type DownloadFileHandler struct {
	DBManager *db.DBManager
}

func NewDownloadFileHandler(dbManager *db.DBManager) *DownloadFileHandler {
	return &DownloadFileHandler{
		DBManager: dbManager,
	}
}

func (h *DownloadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileMetadata.Filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileMetadata.Filesize, 10))
	io.Copy(w, file)
}
