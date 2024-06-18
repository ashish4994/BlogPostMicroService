package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	connectDB()
	defer db.Close()

	r := mux.NewRouter()

	// Set up CORS options
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		Debug:            true, // Set to false in production
	})

	//r.Use(c.Handler)

	r.HandleFunc("/api/blogposts", saveBlogPost).Methods("POST")
	r.HandleFunc("/api/blogposts", getAllPosts).Methods("GET")
	r.HandleFunc("/api/blogposts/{id}", getBlogPostByID).Methods("GET")
	r.HandleFunc("/", serviceStatusHandler)

	handler := c.Handler(r)

	// Get the port number from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))

}
