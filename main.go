package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	connectDB()
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/blogposts", saveBlogPost).Methods("POST")
	r.HandleFunc("/api/blogposts", getAllPosts).Methods("GET")
	r.HandleFunc("/api/blogposts/{id}", getBlogPostByID).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
