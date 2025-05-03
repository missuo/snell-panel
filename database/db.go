/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:23:16
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:23:27
 * @FilePath: /snell-panel/database/db.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB(dbURL string) *sql.DB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create tables if they don't exist
	createTables(db)

	return db
}

// createTables creates the necessary tables in the database
func createTables(db *sql.DB) {
	// Create entries table
	_, err := db.Exec(`
		 CREATE TABLE IF NOT EXISTS entries (
			 id SERIAL PRIMARY KEY,
			 ip TEXT UNIQUE,
			 port INTEGER,
			 psk TEXT,
			 country_code TEXT,
			 isp TEXT,
			 asn INTEGER,
			 node_id TEXT UNIQUE,
			 node_name TEXT
		 )
	 `)
	if err != nil {
		log.Fatalf("Failed to create entries table: %v", err)
	}
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
