package main

import (
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/bergzeit/shortli/repository"
)

type PageData struct {
	OriginalURL string `json:"originalUrl"`
	ShortURL    string `json:"shortUrl"`
}

// urlShortHanlder checks if a short link is already existing in the database
// and send this original link back to the client.
func (app *application) urlShortHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Error: only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)

	shortUrlInput := r.FormValue("shortUrl")

	if len(shortUrlInput) != 0 {
		returnedOriginalURL, err := app.repo.FindOriginalLink(shortUrlInput)
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
func (app *application) createHandler(w http.ResponseWriter, r *http.Request) {

	// MethodOptions must be valid (Status: OK).
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
		return
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

	var shortUrl string
	for i := 0; i < 3; i++ {
		shortUrl, err = generateShortKey()
		if err != nil {
			http.Error(w, "Error: failed to generate a short link", http.StatusInternalServerError)
			return
		}

		err = app.repo.InsertData(pageData.OriginalURL, shortUrl)
		if err == nil {
			break
		}

		if !repository.IsDuplicateKey(err) {
			http.Error(w, "Error: database error", http.StatusInternalServerError)
			return
		}
	}

	if err != nil {
		http.Error(w, "Error: failed to create unique short URL", http.StatusConflict)
		return
	}

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

// forwardingToOriginHandler handles incoming GET requests with a shortlink path,
// looks up the corresponding original URL from the repository, and redirects the client to it.
func (app *application) forwardingToOriginHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Error: only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	pathInput := trimPath(r)

	returnedOriginalURL, err := app.repo.FindOriginalLink(pathInput)
	if err != nil {
		http.Error(w, "Error: can't find original link", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, returnedOriginalURL, http.StatusSeeOther)
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

func trimPath(r *http.Request) string {
	pathInput := strings.TrimPrefix(r.URL.Path, "/")
	return pathInput
}
