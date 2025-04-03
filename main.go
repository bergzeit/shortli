package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	OriginalURL string
	ShortURL    string
	Error       error
}

var urls = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	mux.HandleFunc("/", urlShortHandler)

	log.Print("starting server on :8080")    // print on console.
	err := http.ListenAndServe(":8080", mux) // error handling.
	if err != nil {
		log.Fatal(err)
	}
}

// urlShortHanlder shortens a original URL.
func urlShortHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Check if the request method is POST.
	// If true, extract the "url" value from the form data.
	if r.Method == http.MethodPost {
		originalURL := r.FormValue("url")

		if originalURL == "" { // Error-Handling.
			http.Error(w, "Error (empty url is not valid)", http.StatusBadRequest)
			return
		}

		shortUrl := generateShortKey() // generate the short url
		urls[shortUrl] = originalURL   // fill up the map

		db, err := sql.Open("sqlite3", "shortli.db")
		if err != nil {
			tmpl.Execute(w, PageData{Error: err})
			return
		}
		// instert the both urls in the database (shortli.db)
		insertData(db, originalURL, shortUrl) // error template

		// Data will show on the interface.
		data := PageData{
			OriginalURL: originalURL,
			ShortURL:    shortUrl,
		}

		// Writing the output to the HTTP response.
		tmpl.Execute(w, data)
		return
	}

	// Check if the request method is GET.
	if r.Method == http.MethodGet {
		ShortUrlInput := r.FormValue("shortUrl")

		db, err := sql.Open("sqlite3", "shortli.db")
		if err != nil {
			tmpl.Execute(w, PageData{Error: err})
			return
		}
		defer db.Close()

		if len(ShortUrlInput) != 0 {
			answer, err := findOriginlaLink(db, ShortUrlInput)
			if err != nil {
				tmpl.Execute(w, PageData{Error: err})
				return
			}

			data := PageData{
				OriginalURL: answer,
				ShortURL:    ShortUrlInput,
			}
			tmpl.Execute(w, data)
			return
		}
	}

	// Execute the template with no data and write the output to the HTTP response.
	tmpl.Execute(w, nil)
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

// findOriginlaLink looks for the original Link in the .db database.
func findOriginlaLink(db *sql.DB, shortlink string) (string, error) {
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
	insertLinks := `INSERT INTO links(longlink, shortlink) VALUES (?,?)`
	statement, err := db.Prepare(insertLinks) // Prepare SQL statement

	if err != nil {
		return err
	}
	_, err = statement.Exec(longlink, shortlink)
	if err != nil {
		return err
	}
	return nil
}
