package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

func init() {
	err := ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}
func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /api/healthchecker", HealthCheckerHandler)
	router.HandleFunc("PATCH /api/notes/{noteId}", UpdateNote)
	router.HandleFunc("GET /api/notes/{noteId}", FindNoteById)
	router.HandleFunc("DELETE /api/notes/{noteId}", DeleteNote)
	router.HandleFunc("POST /api/notes", CreateNoteHandler)
	router.HandleFunc("GET /api/notes", FindNotes)

	//Custom CORS configuration
	corsConfig := cors.New(cors.Options{
		AllowedHeaders:   []string{"Origin", "Authorization", "Accept", "Content-Type"},
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	})
	//wrap the router with the logRequest middleware
	loggedRouter := logReuests(router)

	//Create a new CORS handler
	corsHandler := corsConfig.Handler(loggedRouter)

	server := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler,
	}
	log.Println("Starting server on port 8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server Error: %v", err)
	}

}

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
func logReuests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{w, http.StatusOK}
		next.ServeHTTP(wrapped, r)

		elapsed := time.Since(start)
		log.Printf("Received request: %d %s %s", wrapped.statusCode, r.Method, r.URL.Path, elapsed)
	})
}
func HealthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "success",
		"message": "Welcome to the health checker",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
