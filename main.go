package main

import (
	"filestorageapi/db"
	"filestorageapi/handlers"
	"filestorageapi/jwtmanager"
	"filestorageapi/middlewares"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dbURI := "mongodb://localhost:27017"
	dbManager, err := db.NewDBManager(dbURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer dbManager.Disconnect()

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenDuration := time.Hour * 24 // Change this to the duration you want
	jwtManager := jwtmanager.NewJWTManager(jwtSecret, tokenDuration)

	router := mux.NewRouter()

	// User routes
	registerHandler := handlers.NewRegisterHandler(dbManager) // Assuming this structure
	router.Handle("/register", registerHandler).Methods("POST")

	loginHandler := handlers.NewLoginHandler(dbManager, jwtManager)
	router.Handle("/login", loginHandler).Methods("POST")

	// File routes
	fileRouter := router.PathPrefix("/files").Subrouter()
	fileRouter.Use(middlewares.JWTMiddleware(jwtManager))

	uploadHandler := handlers.NewUploadFileHandler(dbManager) // Assuming this structure
	fileRouter.Handle("/upload", uploadHandler).Methods("POST")

	downloadHandler := handlers.NewDownloadFileHandler(dbManager) // Assuming this structure
	fileRouter.Handle("/download/{id}", downloadHandler).Methods("GET")

	deleteHandler := handlers.NewDeleteFileHandler(dbManager) // Assuming this structure
	fileRouter.Handle("/delete/{id}", deleteHandler).Methods("DELETE")

	err = http.ListenAndServe("localhost:8080", router)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
