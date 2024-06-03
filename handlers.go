package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func serviceStatusHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now()
	statusMessage := fmt.Sprintf("BlogPostMicroService is up and running at %s", currentTime.Format(time.RFC1123))
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, statusMessage)
}

func saveBlogPost(w http.ResponseWriter, r *http.Request) {
	var post BlogPost
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert the blog post into the blog_posts table
	var blogPostID int
	err = tx.QueryRow("INSERT INTO blog_posts (name, image_url, content, posted_by, posted_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		post.Name, post.ImageURL, post.Content, post.PostedBy, time.Now().UTC()).Scan(&blogPostID)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert tags into the tags table and associate them with the blog post
	for _, tagName := range post.Tags {
		var tagID int
		err = tx.QueryRow("INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id", tagName).Scan(&tagID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("INSERT INTO blog_post_tags (blog_post_id, tag_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", blogPostID, tagID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT bp.id, bp.name, bp.image_url, bp.content, bp.posted_by, array_agg(t.name) as tags
		FROM blog_posts bp
		LEFT JOIN blog_post_tags bpt ON bp.id = bpt.blog_post_id
		LEFT JOIN tags t ON bpt.tag_id = t.id
		GROUP BY bp.id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		var tags []sql.NullString
		if err := rows.Scan(&post.ID, &post.Name, &post.ImageURL, &post.Content, &post.PostedBy, pq.Array(&tags)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for _, tag := range tags {
			if tag.Valid {
				post.Tags = append(post.Tags, tag.String)
			}
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func getBlogPostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid blog post ID", http.StatusBadRequest)
		return
	}

	var post BlogPost
	var tags []sql.NullString
	err = db.QueryRow(`
		SELECT bp.id, bp.name, bp.image_url, bp.content, bp.posted_by, array_agg(t.name) as tags
		FROM blog_posts bp
		LEFT JOIN blog_post_tags bpt ON bp.id = bpt.blog_post_id
		LEFT JOIN tags t ON bpt.tag_id = t.id
		WHERE bp.id = $1
		GROUP BY bp.id
	`, id).Scan(&post.ID, &post.Name, &post.ImageURL, &post.Content, &post.PostedBy, pq.Array(&tags))
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
