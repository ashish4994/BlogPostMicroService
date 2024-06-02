package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() {
	var err error

	// Read database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("SSL_MODE")
	sslRootCert := os.Getenv("SSL_ROOT_CERT") // Path to the root certificate

	// Construct the connection string using the environment variables
	serviceURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	conn, err := url.Parse(serviceURI)
	if err != nil {
		log.Fatal(err)
	}

	// If SSL is enabled, add the sslrootcert parameter
	if sslMode == "verify-ca" && sslRootCert != "" {
		conn.RawQuery = fmt.Sprintf("sslmode=%s&sslrootcert=%s", sslMode, sslRootCert)
	}
	db, err = sql.Open("postgres", conn.String())
	if err != nil {
		log.Fatal("Indisde err1")
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Indisde err2")
		log.Fatal(err)
	}

	fmt.Println("Connected to the database!")

	rows, err := db.Query("SELECT version()")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var result string
		err = rows.Scan(&result)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Version: %s\n", result)
	}

}
