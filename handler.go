/*
 * @Author: Vincent Yang
 * @Date: 2024-09-06 15:09:31
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-09-06 15:20:34
 * @FilePath: /snell-panel/handler.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */

package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID   int    `json:"id"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
	PSK  string `json:"psk"`
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

func getSubscription(c *gin.Context) {
	rows, err := db.Query("SELECT ip, port, psk FROM entries")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var subscriptionLines []string
	for rows.Next() {
		var ip string
		var port int
		var psk string
		if err := rows.Scan(&ip, &port, &psk); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ipInfo, err := getIPInfo(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		nodeName := ipInfo.CountryCode + " AS" + strconv.Itoa(ipInfo.ASN) + " " + ipInfo.ISP + " " + generateRandomString()
		line := fmt.Sprintf("%s = snell, %s, %d, psk=%s, version=4", nodeName, ip, port, psk)
		subscriptionLines = append(subscriptionLines, line)
	}

	c.String(http.StatusOK, strings.Join(subscriptionLines, "\n"))
}
