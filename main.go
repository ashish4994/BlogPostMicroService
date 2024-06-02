package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connectDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/blogposts", saveBlogPost).Methods("POST")
	r.HandleFunc("/api/blogposts", getAllPosts).Methods("GET")
	r.HandleFunc("/api/blogposts/{id}", getBlogPostByID).Methods("GET")
	r.HandleFunc("/", serviceStatusHandler)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
