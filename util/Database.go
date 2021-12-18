package util

import (
	"database/sql"
	"log"
)

func DatabaseConnect() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/dev.db")
	Check(err)
	return db
}

func InsertUser(db *sql.DB, username string, hash string) error {
	log.Println("Inserting user into database...")
	insertUserSQL := `
		INSERT INTO users(username, hash) VALUES (?, ?)
	`
	stmt, err := db.Prepare(insertUserSQL) // Prepared statements to prevent sql injection
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, hash)
	if err != nil {
		return err
	}

	return nil
}

func GetHashByUsername(db *sql.DB, username string) (string, error) {
	selectUserSQL := `
		SELECT hash FROM users WHERE username = ?
	`
	stmt, err := db.Prepare(selectUserSQL)
	if err != nil {
		return "", err
	}

	row, err := stmt.Query(username)
	if err != nil {
		return "", err
	}

	cols, err := row.Columns()
	if err != nil {
		return "", err
	}

	if len(cols) < 1 {
		return "", nil
	}

	row.Next()

	var hash string
	row.Scan(&hash)

	return hash, nil
}
