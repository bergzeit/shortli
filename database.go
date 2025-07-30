package main

import "database/sql"

// Struct for the connection pool.
type Database struct {
	DB *sql.DB
}

var dbInstance *Database

func DbConnection(dbFile string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(1) // to prevent ("database is locked").
	db.SetMaxIdleConns(1) // to prevent ("database is locked").

	dbInstance = &Database{DB: db}
	return nil
}

// GetDB returns the global database instance.
// To Access the database connection pool everywhere.
func GetDB() *Database {
	return dbInstance
}

func (d *Database) CloseDB() error {
	return d.DB.Close()
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
