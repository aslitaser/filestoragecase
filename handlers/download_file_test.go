package handlers_test

import (
	"filestorageapi/db"
	"filestorageapi/handlers"
	"filestorageapi/models"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	// Set up the necessary dependencies
	dbManager, err := db.NewDBManager("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	downloadFileHandler := handlers.NewDownloadFileHandler(dbManager)

	// Download a file from the server using its ID.
	t.Run("Download a file from the server using its ID", func(t *testing.T) {
		fileMetadata := &models.FileMetadata{
			Filename: "testfile.txt",
			Filesize: 18,
		}

		err := dbManager.SaveFileMetadata(fileMetadata)
		if err != nil {
			t.Fatal(err)
		}

		fileContent := []byte("This is a test file")
		err = ioutil.WriteFile("/path/to/store/files/"+fileMetadata.Filename, fileContent, 0644)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("GET", "/download?id="+fileMetadata.ID, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		downloadFileHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		os.Remove("/path/to/store/files/" + fileMetadata.Filename)
	})

	// Test for invalid file ID.
	t.Run("Test for invalid file ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/download?id=invalid_file_id", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		downloadFileHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
		}
	})
}
