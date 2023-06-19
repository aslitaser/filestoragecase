package handlers_test

import (
	"bytes"
	"filestorageapi/db"
	"filestorageapi/handlers"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestUploadFile(t *testing.T) {
	// Set up the necessary dependencies
	dbManager, err := db.NewDBManager("mongodb://localhost:27017")
	if err != nil {
		t.Fatalf("Failed to connect to DB: %v", err)
	}

	uploadFileHandler := handlers.NewUploadFileHandler(dbManager)

	// Upload a file to the server and store metadata in MongoDB.
	t.Run("Upload a file to the server and store metadata in MongoDB", func(t *testing.T) {
		filePath := "testfile.txt"
		fileSize, _ := createTestFile(filePath)

		req, err := createUploadRequest(filePath, fileSize)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		uploadFileHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		os.Remove(filePath)
	})

	// Test for missing filename or filesize query parameters.
	t.Run("Test for missing filename or filesize query parameters", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/upload", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		uploadFileHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test for invalid filesize.
	t.Run("Test for invalid filesize", func(t *testing.T) {
		filePath := "testfile.txt"
		createTestFile(filePath)
		req, err := createUploadRequest(filePath, -1)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		uploadFileHandler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}

		os.Remove(filePath)
	})
}

func createTestFile(filePath string) (int64, error) {
	content := []byte("This is a test file")
	err := ioutil.WriteFile(filePath, content, 0644)
	if err != nil {
		return 0, err
	}

	fileInfo, _ := os.Stat(filePath)
	return fileInfo.Size(), nil
}

func createUploadRequest(filePath string, fileSize int64) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	err = writer.WriteField("filesize", strconv.FormatInt(fileSize, 10))
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}
