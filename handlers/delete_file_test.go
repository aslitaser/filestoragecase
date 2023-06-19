package handlers

import (
	"filestorageapi/db"
	"filestorageapi/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDeleteFileHandler(t *testing.T) {
	// Setup
	dbURI := "mongodb://localhost:27017"
	dbManager, err := db.NewDBManager(dbURI)
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	fileMetadata := &models.FileMetadata{
		Filename: "testfile.txt",
		Filesize: 100,
	}
	dbManager.SaveFileMetadata(fileMetadata)

	filePath := "/Users/asliftw/GolandProjects/filestorageapi" + fileMetadata.Filename
	err = ioutil.WriteFile(filePath, []byte("This is a test file"), 0644)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}

	deleteFileHandler := NewDeleteFileHandler(dbManager)

	// Test
	req, err := http.NewRequest("DELETE", "/delete_file?id="+fileMetadata.ID, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	recorder := httptest.NewRecorder()
	deleteFileHandler.ServeHTTP(recorder, req)

	// Assert
	if status := recorder.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	_, err = os.Stat(filePath)
	if !os.IsNotExist(err) {
		t.Errorf("Test file not deleted: %s", filePath)
	}

	_, err = dbManager.GetFileMetadata(fileMetadata.ID)
	if err == nil {
		t.Error("File metadata not deleted")
	}
}
