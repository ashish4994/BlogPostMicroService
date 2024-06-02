package main

type BlogPost struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	ImageURL string   `json:"image_url"`
	Content  string   `json:"content"`
	PostedBy string   `json:"posted_by"`
	Tags     []string `json:"tags"`
}

// Define other models and database interaction functions here
