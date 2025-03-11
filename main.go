package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type PageData struct {
	OriginalURL string
	ShortURL    string
}

var urls = make(map[string]string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", urlShortHandler)

	log.Print("starting server on :8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

// urlShortHanlder shortens a original URL.
func urlShortHandler(w http.ResponseWriter, r *http.Request) {

	// Check if the request method is POST.
	// If true, extract the "url" value from the form data.
	if r.Method == http.MethodPost {
		originalURL := r.FormValue("url")

		if originalURL == "" { // Error-Handling.
			http.Error(w, "Error (empty url is not valid)", http.StatusBadRequest) //ändern
			return
		}

		shortUrl := generateShortKey() // generate the short url
		urls[shortUrl] = originalURL   // fill up the map

		db, err := sql.Open("sqlite3", "shortli.db")
		if err != nil { //ändern
			fmt.Println(err) // return (Datenbank nicht verfügbar)
		}
		// instert the both urls in the database (shortli.db)
		insertData(db, originalURL, shortUrl) // error template

		// Data will show on the interface.
		data := PageData{
			OriginalURL: originalURL,
			ShortURL:    shortUrl,
		}

		// Writing the output to the HTTP response.
		tmpl := template.Must(template.ParseFiles("templates/index.html")) //..
		tmpl.Execute(w, data)
		return
	}

	if r.Method == http.MethodGet {
		ShortUrlInput := r.FormValue("shortUrl")

		db, err := sql.Open("sqlite3", "shortli.db")
		if err != nil {
			fmt.Println(err)
		}

		allData, _ := findAll(db)
		for _, v := range allData {
			if v.ShortURL == ShortUrlInput {

				data := PageData{
					OriginalURL: v.OriginalURL,
					ShortURL:    v.ShortURL,
				}
				// Writing the output to the HTTP response.
				tmpl := template.Must(template.ParseFiles("templates/index.html"))
				tmpl.Execute(w, data)
				return
			}
		}
	}
	// Execute the template with no data and write the output to the HTTP response.
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
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

// FindAll looks for all data in the shortli.db database and returns it in a slice.
func findAll(db *sql.DB) ([]PageData, error) { // Where (SQL)
	sql := `SELECT * FROM links`

	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []PageData
	for rows.Next() {
		c := &PageData{}
		err := rows.Scan(&c.OriginalURL, &c.ShortURL)
		if err != nil {
			return nil, err
		}
		links = append(links, *c)
	}
	return links, nil
}

// instertData saves the longlink and shortlink in sql database.
func insertData(db *sql.DB, longlink string, shortlink string) {
	insertLinks := `INSERT INTO links(longlink, shortlink) VALUES (?,?)`
	statement, err := db.Prepare(insertLinks)

	if err != nil { // ändern error return
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(longlink, shortlink)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
