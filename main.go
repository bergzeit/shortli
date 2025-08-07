package main

import (
	"log"
	"net/http"

	"github.com/bergzeit/shortli/repository"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type application struct {
	repo *repository.Repository
}

func main() {
	db, err := repository.New("shortli.db")
	if err != nil { // Initialize (once) shortli.db database.
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close() // close DB at the end.

	app := application{
		repo: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.forwardingToOriginHandler)
	mux.HandleFunc("/api/v1/", app.urlShortHandler)
	mux.HandleFunc("/api/v1/create", app.createHandler)

	corsOptions := cors.Options{
		AllowedOrigins: []string{"http://localhost*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}
	handler := cors.New(corsOptions).Handler(mux) // Cors-Middleware

	log.Print("starting server on: 8080")       // print on console.
	err = http.ListenAndServe(":8080", handler) // error handling.
	if err != nil {
		log.Fatal(err)
	}
}
