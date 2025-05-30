/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:24:49
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:25:05
 * @FilePath: /snell-panel/handlers/handlers.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"snell-panel/models"
	"snell-panel/utils"
)

// Handlers contains the HTTP request handlers
type Handlers struct {
	DB    *sql.DB
	Token string
}

// NewHandlers creates a new Handlers instance
func NewHandlers(db *sql.DB, token string) *Handlers {
	return &Handlers{
		DB:    db,
		Token: token,
	}
}

// CorsMiddleware returns a CORS middleware configured for the API
func CorsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	return cors.New(config)
}

// AuthMiddleware returns a middleware that checks for API token
func (h *Handlers) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		providedToken := c.Query("token")
		if providedToken != h.Token {
			c.JSON(http.StatusUnauthorized, models.ApiResponse{
				Status:  "error",
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// Welcome handles the root route
func (h *Handlers) Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Welcome to Snell Panel. Please use the API to manage the entries.\n https://github.com/missuo/snell-panel",
	})
}

// NotFound handles not found routes
func (h *Handlers) NotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, models.ApiResponse{
		Status:  "error",
		Message: "Path not found",
	})
}

// InsertEntry handles creating a new entry
func (h *Handlers) InsertEntry(c *gin.Context) {
	var entry models.Entry
	if err := c.BindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Get IP information
	ipInfo, err := utils.GetIPInfo(entry.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: "Failed to get IP info",
		})
		return
	}

	// Update entry with IP information
	entry.CountryCode = ipInfo.CountryCode
	entry.ISP = ipInfo.ISP
	entry.ASN = ipInfo.ASN
	entry.NodeID = utils.GenerateUUID()

	// Insert entry into database
	var id int
	err = h.DB.QueryRow(`
		 INSERT INTO entries (ip, port, psk, country_code, isp, asn, node_id, node_name) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		 RETURNING id`,
		entry.IP, entry.Port, entry.PSK, entry.CountryCode, entry.ISP, entry.ASN, entry.NodeID, entry.NodeName).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Set ID in entry
	entry.ID = id

	c.JSON(http.StatusCreated, models.ApiResponse{
		Status:  "success",
		Message: "Entry created successfully",
		Data:    entry,
	})
}

// DeleteEntryByIP handles deleting an entry by IP
func (h *Handlers) DeleteEntryByIP(c *gin.Context) {
	ip := c.Param("ip")

	result, err := h.DB.Exec("DELETE FROM entries WHERE ip = $1", ip)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "Entry not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Entry deleted successfully",
	})
}

// QueryAllEntries handles retrieving all entries
func (h *Handlers) QueryAllEntries(c *gin.Context) {
	rows, err := h.DB.Query(`
		 SELECT id, ip, port, psk, country_code, isp, asn, node_id, node_name 
		 FROM entries
	 `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var entries []models.Entry
	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(
			&entry.ID, &entry.IP, &entry.Port, &entry.PSK,
			&entry.CountryCode, &entry.ISP, &entry.ASN,
			&entry.NodeID, &entry.NodeName,
		); err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}
		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "warning",
			Message: "No entries found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Entries retrieved successfully",
		Data:    entries,
	})
}

// GetSubscription handles generating a subscription string
func (h *Handlers) GetSubscription(c *gin.Context) {
	rows, err := h.DB.Query(`
		 SELECT ip, port, psk, country_code, isp, asn, node_id, node_name 
		 FROM entries
	 `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}
	defer rows.Close()

	var subscriptionLines []string
	for rows.Next() {
		var entry models.Entry
		if err := rows.Scan(
			&entry.IP, &entry.Port, &entry.PSK,
			&entry.CountryCode, &entry.ISP, &entry.ASN,
			&entry.NodeID, &entry.NodeName,
		); err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		emojiFlag := utils.CountryCodeToFlagEmoji(entry.CountryCode)
		nodeName := entry.NodeName
		if nodeName == "" {
			nodeName = fmt.Sprintf("%s %s AS%d %s %s",
				emojiFlag, entry.CountryCode, entry.ASN, entry.ISP, entry.NodeID)
		} else {
			nodeName = fmt.Sprintf("%s %s", emojiFlag, entry.NodeName)
		}

		line := fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = 4",
			nodeName, entry.IP, entry.Port, entry.PSK)
		subscriptionLines = append(subscriptionLines, line)
	}

	if len(subscriptionLines) == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "No entries found for subscription",
		})
		return
	}

	c.String(http.StatusOK, strings.Join(subscriptionLines, "\n"))
}

// ModifyNodeByNodeID handles modifying a node by its NodeID
func (h *Handlers) ModifyNodeByNodeID(c *gin.Context) {
	nodeID := c.Param("id")

	var modifyReq models.ModifyRequest
	if err := c.BindJSON(&modifyReq); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// Build query based on which fields are provided
	query := "UPDATE entries SET"
	var args []interface{}
	var setStatements []string
	paramIndex := 1

	if modifyReq.NodeName != "" {
		setStatements = append(setStatements, fmt.Sprintf(" node_name = $%d", paramIndex))
		args = append(args, modifyReq.NodeName)
		paramIndex++
	}

	if modifyReq.IP != "" {
		setStatements = append(setStatements, fmt.Sprintf(" ip = $%d", paramIndex))
		args = append(args, modifyReq.IP)
		paramIndex++

		// If IP is updated, we should also update geolocation info
		ipInfo, err := utils.GetIPInfo(modifyReq.IP)
		if err == nil {
			setStatements = append(setStatements, fmt.Sprintf(" country_code = $%d", paramIndex))
			args = append(args, ipInfo.CountryCode)
			paramIndex++

			setStatements = append(setStatements, fmt.Sprintf(" isp = $%d", paramIndex))
			args = append(args, ipInfo.ISP)
			paramIndex++

			setStatements = append(setStatements, fmt.Sprintf(" asn = $%d", paramIndex))
			args = append(args, ipInfo.ASN)
			paramIndex++
		}
	}

	// If no fields to update, return error
	if len(setStatements) == 0 {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status:  "error",
			Message: "No fields to update",
		})
		return
	}

	// Combine all set statements
	query += strings.Join(setStatements, ",")
	query += fmt.Sprintf(" WHERE node_id = $%d", paramIndex)
	args = append(args, nodeID)

	// Execute update
	result, err := h.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status:  "error",
			Message: "Node ID not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status:  "success",
		Message: "Node updated successfully",
	})
}
