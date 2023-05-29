package internal

import (
	"database/sql"
	"log"
)

func createTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS system_resources (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		cpu_percent FLOAT,
		mem_percent FLOAT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

// Write a function to check the table exists or not
func tableExists(db *sql.DB) bool {
	query := `SELECT name FROM sqlite_master WHERE type='table' AND name='system_resources';`
	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to check if table exists: %v", err)
	}
	defer rows.Close()

	return rows.Next()
}

func conn() *sql.DB {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open the database: %v", err)
	}
	return db
}
