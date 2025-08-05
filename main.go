package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func main() {
	db, err := DbConnection("shortli.db")
	if err != nil { // Initialize (once) shortli.db database.
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close() // close DB at the end.

	app := Application{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.UrlShortHandler)
	mux.HandleFunc("/create", app.CreateHandler)

	corsOptions := cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}
	handler := cors.New(corsOptions).Handler(mux) // Cors-Middleware

	log.Print("starting server on :8080")       // print on console.
	err = http.ListenAndServe(":8080", handler) // error handling.
	if err != nil {
		log.Fatal(err)
	}
}
