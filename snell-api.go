/*
 * @Author: Vincent Yang
 * @Date: 2024-09-06 14:36:44
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-09-06 16:22:27
 * @FilePath: /snell-panel/snell-api.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */
package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db    *sql.DB
	token string
)

func init() {
	flag.StringVar(&token, "token", "", "API access token")
	flag.Parse()

	if token == "" {
		token = os.Getenv("API_TOKEN")
	}

	if token == "" {
		log.Fatal("API token must be provided via command line argument or environment variable")
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		providedToken := c.Query("token")
		if providedToken != token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	r := gin.Default()

	r.StaticFile("/", "./static/index.html")
	r.POST("/entry", insertEntry)
	r.GET("/entries", authMiddleware(), queryAllEntries)
	r.DELETE("/entry/:ip", deleteEntryByIP)
	r.GET("/subscribe", authMiddleware(), getSubscription)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "Path not found",
		})
	})
	r.Run(":59999")
}

func createTable() {
	statement, err := db.Prepare(`
		 CREATE TABLE IF NOT EXISTS entries (
			 id INTEGER PRIMARY KEY AUTOINCREMENT,
			 ip TEXT UNIQUE,
			 port INTEGER,
			 psk TEXT,
			 country_code TEXT,
			 isp TEXT,
			 asn TEXT,
			 node_id TEXT UNIQUE
		 )
	 `)
	if err != nil {
		log.Fatal(err)
	}
	statement.Exec()
}
