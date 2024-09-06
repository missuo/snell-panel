/*
 * @Author: Vincent Yang
 * @Date: 2024-09-06 14:36:44
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-09-06 15:09:56
 * @FilePath: /snell-panel/snell-api.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */
package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

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
	r.GET("/subscribe", getSubscription)

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
