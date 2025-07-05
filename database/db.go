/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:23:16
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-07-05 20:10:02
 * @FilePath: /snell-panel/database/db.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package database

import (
	"database/sql"
	"fmt"
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
	// Check if the table already exists
	var tableExists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_name = 'entries'
		)
	`).Scan(&tableExists)

	if err != nil {
		log.Fatalf("Failed to check if entries table exists: %v", err)
	}

	if !tableExists {
		// Create entries table without IP UNIQUE constraint (for new installations)
		_, err = db.Exec(`
			CREATE TABLE entries (
				id SERIAL PRIMARY KEY,
				ip TEXT,
				port INTEGER,
				psk TEXT,
				country_code TEXT,
				isp TEXT,
				asn INTEGER,
				node_id TEXT UNIQUE,
				node_name TEXT,
				version TEXT DEFAULT '4'
			)
		`)
		if err != nil {
			log.Fatalf("Failed to create entries table: %v", err)
		}
		log.Printf("Created entries table without IP unique constraint")
	} else {
		// For existing installations, add version column if it doesn't exist
		_, err = db.Exec(`
			ALTER TABLE entries ADD COLUMN IF NOT EXISTS version TEXT DEFAULT '4'
		`)
		if err != nil {
			log.Fatalf("Failed to add version column: %v", err)
		}

		// Remove UNIQUE constraint from ip column for existing installations
		removeIPUniqueConstraint(db)
	}
}

// removeIPUniqueConstraint removes the UNIQUE constraint from the ip column
func removeIPUniqueConstraint(db *sql.DB) {
	// First, check if there are any duplicate IPs that would prevent constraint removal
	var duplicateCount int
	err := db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT ip FROM entries 
			GROUP BY ip 
			HAVING COUNT(*) > 1
		) as duplicates
	`).Scan(&duplicateCount)

	if err != nil {
		log.Printf("Warning: Failed to check for duplicate IPs: %v", err)
		return
	}

	if duplicateCount > 0 {
		log.Printf("Warning: Found %d duplicate IP addresses. The UNIQUE constraint cannot be removed safely.", duplicateCount)
		log.Printf("Please review and resolve duplicate IPs before proceeding.")
		return
	}

	// Check if the constraint exists
	var constraintExists bool
	err = db.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints 
			WHERE table_name = 'entries' 
			AND constraint_type = 'UNIQUE' 
			AND constraint_name LIKE '%ip%'
		)
	`).Scan(&constraintExists)

	if err != nil {
		log.Printf("Warning: Failed to check IP constraint existence: %v", err)
		return
	}

	if !constraintExists {
		log.Printf("IP unique constraint does not exist, skipping removal")
		return
	}

	// Get the constraint name
	var constraintName string
	err = db.QueryRow(`
		SELECT constraint_name 
		FROM information_schema.table_constraints 
		WHERE table_name = 'entries' 
		AND constraint_type = 'UNIQUE' 
		AND constraint_name LIKE '%ip%'
		LIMIT 1
	`).Scan(&constraintName)

	if err != nil {
		log.Printf("Warning: Failed to get IP constraint name: %v", err)
		return
	}

	// Log the action for audit purposes
	log.Printf("Attempting to remove IP unique constraint: %s", constraintName)

	// Drop the constraint
	_, err = db.Exec(fmt.Sprintf("ALTER TABLE entries DROP CONSTRAINT IF EXISTS %s", constraintName))
	if err != nil {
		log.Printf("Warning: Failed to drop IP unique constraint: %v", err)
		log.Printf("This may indicate duplicate IP addresses in the database.")
		log.Printf("Please backup your data and manually resolve any conflicts.")
	} else {
		log.Printf("Successfully removed IP unique constraint: %s", constraintName)
		log.Printf("IP addresses can now be duplicated across entries.")
	}
}

// CloseDB closes the database connection
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
