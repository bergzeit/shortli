package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	OriginalURL string `json:"originalUrl"`
	ShortURL    string `json:"shortUrl"`
	Error       string `json:"error"`
}

var urls = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortHandler)
	mux.HandleFunc("/create", createHandler)

	log.Print("starting server on :8080")    // print on console.
	err := http.ListenAndServe(":8080", mux) // error handling.
	if err != nil {
		log.Fatal(err)
	}
}

// CORS-Header to allow Cross-Origin-Requests with localhost:5173
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// urlShortHanlder checks if a short link is already existing in the database
// and send this original link back to the client.
func urlShortHandler(w http.ResponseWriter, r *http.Request) {

	enableCORS(w)

	// Check if the request method is GET.
	if r.Method == http.MethodGet {
		ShortUrlInput := r.FormValue("shortUrl")

		db, err := sql.Open("sqlite3", "shortli.db")
		if err != nil {
			http.Error(w, "Error (database)", http.StatusBadRequest)
			return
		}
		defer db.Close()

		if len(ShortUrlInput) != 0 {
			answer, err := findOriginalLink(db, ShortUrlInput)
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

			return
		}
	}
}

// createHanlder receives the original link (JSON) from client (POST-METHOD)
// and create a new shortlink.
func createHandler(w http.ResponseWriter, r *http.Request) {

	enableCORS(w)

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

	var dataVue PageData
	err := json.NewDecoder(r.Body).Decode(&dataVue)
	if err != nil {
		http.Error(w, "something went wrong with JSON", http.StatusBadRequest)
		return
	}

	if dataVue.OriginalURL == "" { // Error-Handling.
		http.Error(w, "Error (empty url is not valid)", http.StatusBadRequest)
		return
	}

	shortUrl := generateShortKey()       // generate the short url
	urls[shortUrl] = dataVue.OriginalURL // fill up the map

	db, err := sql.Open("sqlite3", "shortli.db")
	if err != nil {
		http.Error(w, "Error (database)", http.StatusBadRequest)
		return
	}

	// instert the both urls in the database (shortli.db)
	insertData(db, dataVue.OriginalURL, shortUrl)

	data := PageData{
		OriginalURL: dataVue.OriginalURL,
		ShortURL:    shortUrl,
	}

	// Create a JSON data.
	// Marhal-method encoding input in a JSON-encoded-string.
	jsonStr, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// findOriginalLink looks for the original Link in the .db database.
func findOriginalLink(db *sql.DB, shortlink string) (string, error) {
	sql := `SELECT longlink FROM links WHERE shortlink = ?` // retrieve shortlink by longlink
	row := db.QueryRow(sql, shortlink)                      // execute the SELECT statement
	result := ""

	err := row.Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// instertData saves the longlink and shortlink in sql database.
func insertData(db *sql.DB, longlink string, shortlink string) error {
	insertLinks := `INSERT OR IGNORE INTO links(longlink, shortlink) VALUES (?, ?);` // ignore duplicates
	statement, err := db.Prepare(insertLinks)                                        // Prepare SQL statement

	if err != nil {
		return err
	}
	_, err = statement.Exec(longlink, shortlink)
	if err != nil {
		return err
	}
	return nil
}
