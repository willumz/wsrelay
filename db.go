package main

import (
    "database/sql"
    "fmt"
    _ "modernc.org/sqlite"
    "os"
    "path/filepath"
)

var db *sql.DB

// InitDB initializes the database connection to the SQLite database
// located at "data/data.db". It creates the data directory if it doesn't exist.
func InitDB() (*sql.DB, error) {
    // Ensure the data directory exists
    dataDir := "data"
    if _, err := os.Stat(dataDir); os.IsNotExist(err) {
        if err := os.MkdirAll(dataDir, 0755); err != nil {
            return nil, fmt.Errorf("failed to create data directory: %w", err)
        }
    }

    // Connect to the SQLite database
    dbPath := filepath.Join(dataDir, "data.db")
    database, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    // Test the connection
    if err := database.Ping(); err != nil {
        database.Close()
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    database.Exec(`
		create table if not exists groups (
		    group_id integer primary key,
			name text,
			created_time integer,
			modified_time integer
		)
	`)
    database.Exec(`
			create table if not exists auth_keys (
			    key_id integer primary key,
				key text,
				description text,
				created_time integer,
				group_id integer,
				foreign key (group_id) references groups(group_id)
			)
		`)

    db = database
    fmt.Printf("Connected to SQLite database at %s\n", dbPath)
    return database, nil
}

func GetDB() *sql.DB {
    return db
}

func CloseDB() error {
    if db != nil {
        return db.Close()
    }
    return nil
}

func GetGroupByKey(key string) string {
    row := db.QueryRow(`select group_id from auth_keys where key = ?`, key)
    var groupId string
    err := row.Scan(&groupId)
    if err != nil {
        return ""
    }
    return groupId
}
