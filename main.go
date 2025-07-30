package main

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type PageData struct {
	OriginalURL string `json:"originalUrl"`
	ShortURL    string `json:"shortUrl"`
	Error       string `json:"error"`
}

var urls = make(map[string]string)

func main() {
	if err := DbConnection("shortli.db"); err != nil { // Initialize (once) shortli.db database.
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer GetDB().CloseDB() // close DB at the end.

	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortHandler)
	mux.HandleFunc("/create", createHandler)

	corsOptions := cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}
	handler := cors.New(corsOptions).Handler(mux) // Cors-Middleware

	log.Print("starting server on :8080")        // print on console.
	err := http.ListenAndServe(":8080", handler) // error handling.
	if err != nil {
		log.Fatal(err)
	}
}

// urlShortHanlder checks if a short link is already existing in the database
// and send this original link back to the client.
func urlShortHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET.
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)
	}

	ShortUrlInput := r.FormValue("shortUrl")
	db := GetDB()

	if len(ShortUrlInput) != 0 {
		answer, err := FindOriginalLink(db.DB, ShortUrlInput)
		if err != nil {
			http.Error(w, "Error (can't find original link)", http.StatusBadRequest)
			return
		}

		data := PageData{
			OriginalURL: answer,
			ShortURL:    ShortUrlInput,
		}

		// Create a JSON data.
		// Marhal-method encoding input in a JSON-encoded-string.
		jsonStr, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(jsonStr)
	}
}

// createHanlder receives the original link (JSON) from client (POST-METHOD)
// and create a new shortlink.
func createHandler(w http.ResponseWriter, r *http.Request) {

	// MehtodOptions must be valid (Status: OK).
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if the request method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "only POST-METHOD is avaible", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error while reading the file", http.StatusBadRequest)
	}
	defer r.Body.Close()

	var dataVue PageData
	if err := json.Unmarshal(body, &dataVue); err != nil {
		http.Error(w, "something went wrong with JSON", http.StatusBadRequest)
		return
	}

	if dataVue.OriginalURL == "" { // Error-Handling.
		http.Error(w, "Error (empty url is not valid)", http.StatusBadRequest)
		return
	}

	shortUrl := generateShortKey()       // generate the short url
	urls[shortUrl] = dataVue.OriginalURL // fill up the map

	db := GetDB()
	InsertData(db.DB, dataVue.OriginalURL, shortUrl)

	data := PageData{
		OriginalURL: dataVue.OriginalURL,
		ShortURL:    shortUrl,
	}

	// Create a JSON data.
	// Marhal-method encoding input in a JSON-encoded-string.
	jsonStr, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonStr)
}

// generateShortKey generates and returns a random key with 7 chars.
func generateShortKey() string {
	const keyLength = 7
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, keyLength)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
