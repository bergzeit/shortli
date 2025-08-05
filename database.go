package main

import (
	"database/sql"
)

func DbConnection(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // to prevent ("database is locked").
	db.SetMaxIdleConns(1) // to prevent ("database is locked").

	dbScript(db)
	return db, nil
}

// dbScript executes a SQL script to create the 'links' table if it doesn't already exist.
// It prepares and executes the statement using the provided database connection.
func dbScript(db *sql.DB) error {
	sqlScript := `CREATE TABLE IF NOT EXISTS links (id INTEGER PRIMARY KEY, longlink TEXT NOT NULL, shortlink TEXT NOT NULL UNIQUE);`

	statement, err := db.Prepare(sqlScript) // Prepare SQL statement
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	return nil
}

// FindOriginalLink looks for the original Link in the .db database.
func FindOriginalLink(db *sql.DB, shortlink string) (string, error) {
	result := ""
	sql := `SELECT longlink FROM links WHERE shortlink = ?` // retrieve shortlink by longlink
	row := db.QueryRow(sql, shortlink)                      // execute the SELECT statement

	err := row.Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// InstertData saves the longlink and shortlink in sql database.
func InsertData(db *sql.DB, longlink string, shortlink string) error {
	insertLinks := `INSERT OR IGNORE INTO links(longlink, shortlink) VALUES (?, ?);` // ignore duplicates

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
