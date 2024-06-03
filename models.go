package main

import "time"

type BlogPost struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	ImageURL   string    `json:"image_url"`
	Content    string    `json:"content"`
	PostedBy   string    `json:"posted_by"`
	Tags       []string  `json:"tags"`
	PostedDate time.Time `json:"posted_date"`
}

// Define other models and database interaction functions here
