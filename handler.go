package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Entry struct {
	ID          int    `json:"id"`
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	PSK         string `json:"psk"`
	CountryCode string `json:"country_code"`
	ISP         string `json:"isp"`
	ASN         int    `json:"asn"`
	NodeID      string `json:"node_id"`
	NodeName    string `json:"node_name"`
}

func insertEntry(c *gin.Context) {
	var entry Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ipInfo, err := getIPInfo(entry.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get IP info"})
		return
	}

	entry.CountryCode = ipInfo.CountryCode
	entry.ISP = ipInfo.ISP
	entry.ASN = ipInfo.ASN
	entry.NodeID = generateRandomString()

	result, err := db.Exec("INSERT INTO entries (ip, port, psk, country_code, isp, asn, node_id, node_name) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		entry.IP, entry.Port, entry.PSK, entry.CountryCode, entry.ISP, entry.ASN, entry.NodeID, entry.NodeName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := result.LastInsertId()
	entry.ID = int(id)

	c.JSON(http.StatusCreated, entry)
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

func queryAllEntries(c *gin.Context) {
	rows, err := db.Query("SELECT id, ip, port, psk, country_code, isp, asn, node_id, node_name FROM entries")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.ID, &entry.IP, &entry.Port, &entry.PSK, &entry.CountryCode, &entry.ISP, &entry.ASN, &entry.NodeID, &entry.NodeName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No entries found"})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func getSubscription(c *gin.Context) {
	rows, err := db.Query("SELECT ip, port, psk, country_code, isp, asn, node_id, node_name FROM entries")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var subscriptionLines []string
	for rows.Next() {
		var entry Entry
		if err := rows.Scan(&entry.IP, &entry.Port, &entry.PSK, &entry.CountryCode, &entry.ISP, &entry.ASN, &entry.NodeID, &entry.NodeName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		emojiFlag := CountryCodeToFlagEmoji(entry.CountryCode)
		nodeName := emojiFlag + " " + entry.NodeName
		if nodeName == "" {
			nodeName = fmt.Sprintf("%s %s AS%d %s %s", emojiFlag, entry.CountryCode, entry.ASN, entry.ISP, entry.NodeID)
		}
		line := fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = 4", nodeName, entry.IP, entry.Port, entry.PSK)
		subscriptionLines = append(subscriptionLines, line)
	}

	if len(subscriptionLines) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No entries found for subscription"})
		return
	}

	c.String(http.StatusOK, strings.Join(subscriptionLines, "\n"))
}
