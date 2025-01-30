package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortHandler)

	log.Print("starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

// urlShort shortens a original URL.
func urlShortHandler(w http.ResponseWriter, r *http.Request) {

	var (
		urlMap      = make(map[string]uuid.UUID) // URL -> UUID
		shortUrlMap = make(map[uuid.UUID]string) // UUID -> ShortURL
	)

	uu := uuid.New()                 // generate UUID number
	url := "https://google.com/test" // test url
	shortUrl := generateShortKey()   // a unique shortlink

	urlMap[url] = uu
	shortUrlMap[uu] = shortUrl

	mapOne := fmt.Sprintf("\n\nurlMap:", urlMap)         // only for test and visualize
	mapTwo := fmt.Sprintf("\nshortUrlMap:", shortUrlMap) // only for test and visualize
	answer := fmt.Sprintf("Original Link: %s\nShort Link: %s", url, shortUrl)
	w.Write([]byte(answer))
	w.Write([]byte(mapOne)) // only for test and visualize
	w.Write([]byte(mapTwo)) // only for test and visualize

}

// generateShortKey generates and returns a random key with 7 chars.
func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 7

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	shortUrl := "bg/" + string(shortKey)
	return string(shortUrl)
}
