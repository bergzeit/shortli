package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
)

type Application struct {
	db *sql.DB
}

type PageData struct {
	OriginalURL string `json:"originalUrl"`
	ShortURL    string `json:"shortUrl"`
}

// urlShortHanlder checks if a short link is already existing in the database
// and send this original link back to the client.
func (app *Application) UrlShortHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Error: only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)

	shortUrlInput := r.FormValue("shortUrl")

	if len(shortUrlInput) != 0 {
		returnedOriginalURL, err := FindOriginalLink(app.db, shortUrlInput)
		if err != nil {
			http.Error(w, "Error: can't find original link", http.StatusBadRequest)
			return
		}

		data := PageData{
			OriginalURL: returnedOriginalURL,
			ShortURL:    shortUrlInput,
		}

		// Create a JSON data.
		// Marhal-method encoding input in JSON-encoded-bytes.
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Write(jsonData)
	}
}

// createHanlder receives the original link (JSON) from client (POST-METHOD)
// and create a new shortlink.
func (app *Application) CreateHandler(w http.ResponseWriter, r *http.Request) {

	// MehtodOptions must be valid (Status: OK).
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check if the request method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Error: only POST-METHOD is avaible", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error: while reading the file", http.StatusBadRequest)
	}
	defer r.Body.Close()

	var pageData PageData
	if err := json.Unmarshal(body, &pageData); err != nil {
		http.Error(w, "Error: invalid JSON format – expected valid key-value structure", http.StatusBadRequest)
		return
	}

	if pageData.OriginalURL == "" { // Error-Handling.
		http.Error(w, "Error: empty url is not valid", http.StatusBadRequest)
		return
	}

	shortUrl, _ := generateShortKey() // generate the short url
	InsertData(app.db, pageData.OriginalURL, shortUrl)

	data := PageData{
		OriginalURL: pageData.OriginalURL,
		ShortURL:    shortUrl,
	}

	// Create a JSON data.
	// Marhal-method encoding input in JSON-encoded-bytes.
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

// generateShortKey generates and returns a random key with 7 chars.
func generateShortKey() (string, error) {
	const keyLength = 7
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, keyLength)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b), nil
}
