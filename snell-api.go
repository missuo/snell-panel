package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Entry struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
	PSK  string `json:"psk"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	r := gin.Default()

	r.POST("/entry", insertEntry)
	r.GET("/entries", queryAllEntries)
	r.DELETE("/entry/:ip", deleteEntryByIP)

	r.Run(":59999")
}

func createTable() {
	statement, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ip TEXT UNIQUE,
			port INTEGER,
			psk TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}

func insertEntry(c *gin.Context) {
	var entry Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO entries (ip, port, psk) VALUES (?, ?, ?)",
		entry.IP, entry.Port, entry.PSK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	entry.ID = int(id)

	c.JSON(http.StatusCreated, entry)
}

func queryAllEntries(c *gin.Context) {
	rows, err := db.Query("SELECT id, ip, port, psk FROM entries")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.IP, &entry.Port, &entry.PSK); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		entries = append(entries, entry)
	}

	c.JSON(http.StatusOK, entries)
}

func deleteEntryByIP(c *gin.Context) {
	ip := c.Param("ip")

	result, err := db.Exec("DELETE FROM entries WHERE ip = ?", ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}
