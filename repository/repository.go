package repository

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	db *sql.DB
}

func New(dbFile string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // to prevent ("database is locked").
	db.SetMaxIdleConns(1) // to prevent ("database is locked").

	r := Repository{
		db: db,
	}

	err = r.createTable()
	if err != nil {
		return &Repository{}, fmt.Errorf("can't create table: %w", err)
	}
	return &r, nil
}

// createTable executes a SQL script to create the 'links' table if it doesn't already exist.
// It prepares and executes the statement using the provided database connection.
func (r *Repository) createTable() error {
	sqlScript := `CREATE TABLE IF NOT EXISTS links (longlink TEXT NOT NULL, shortlink TEXT NOT NULL UNIQUE PRIMARY KEY);`

	statement, err := r.db.Prepare(sqlScript) // Prepare SQL statement
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
func (r *Repository) FindOriginalLink(shortlink string) (string, error) {
	result := ""
	sql := `SELECT longlink FROM links WHERE shortlink = ?` // retrieve shortlink by longlink
	row := r.db.QueryRow(sql, shortlink)                    // execute the SELECT statement

	err := row.Scan(&result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// InstertData saves the longlink and shortlink in sql database.
func (r *Repository) InsertData(longlink string, shortlink string) error {
	insertLinks := `INSERT INTO links(longlink, shortlink) VALUES (?, ?);` // ignore duplicates

	statement, err := r.db.Prepare(insertLinks) // Prepare SQL statement
	if err != nil {
		return err
	}
	_, err = statement.Exec(longlink, shortlink)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) Close() {
	r.db.Close()
}
