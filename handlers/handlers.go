/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:24:49
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-07-05 20:40:15
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

	// Resolve domain to IP if needed and get IP information
	// Keep original domain/IP in entry.IP, only use resolved IP for getting geo info
	originalIP := entry.IP
	_, ipInfo, err := utils.GetIPInfoFromDomainOrIP(entry.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to resolve domain/IP or get IP info: %v", err),
		})
		return
	}

	// Update entry with IP information (but keep original domain/IP)
	entry.IP = originalIP // Keep the original domain/IP address
	entry.CountryCode = ipInfo.CountryCode
	entry.ISP = ipInfo.ISP
	entry.ASN = ipInfo.ASN
	entry.NodeID = utils.GenerateUUID()

	// Set default version if not provided
	if entry.Version == "" {
		entry.Version = "4"
	}

	// Insert entry into database
	var id int
	err = h.DB.QueryRow(`
		 INSERT INTO entries (ip, port, psk, country_code, isp, asn, node_id, node_name, version) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		 RETURNING id`,
		entry.IP, entry.Port, entry.PSK, entry.CountryCode, entry.ISP, entry.ASN, entry.NodeID, entry.NodeName, entry.Version).Scan(&id)
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

// DeleteEntryByNodeID handles deleting an entry by node ID
func (h *Handlers) DeleteEntryByNodeID(c *gin.Context) {
	nodeID := c.Param("node_id")

	result, err := h.DB.Exec("DELETE FROM entries WHERE node_id = $1", nodeID)
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
		 SELECT id, ip, port, psk, country_code, isp, asn, node_id, node_name, version 
		 FROM entries
		 ORDER BY id
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
			&entry.NodeID, &entry.NodeName, &entry.Version,
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
	// Get the via and filter parameters from query string
	via := c.Query("via")
	filter := c.Query("filter")
	
	var query string
	var args []interface{}
	
	if filter != "" {
		// Filter nodes by node name containing the keyword
		query = `
			SELECT ip, port, psk, country_code, isp, asn, node_id, node_name, version 
			FROM entries
			WHERE node_name LIKE ?
			ORDER BY id
		`
		args = []interface{}{"%" + filter + "%"}
	} else {
		query = `
			SELECT ip, port, psk, country_code, isp, asn, node_id, node_name, version 
			FROM entries
			ORDER BY id
		`
	}
	
	rows, err := h.DB.Query(query, args...)
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
			&entry.NodeID, &entry.NodeName, &entry.Version,
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
		
		// Add [Via xxx] suffix to node name when via parameter is provided
		if via != "" {
			nodeName = fmt.Sprintf("%s [Via %s]", nodeName, via)
		}

		var line string
		if via != "" {
			// Include underlying-proxy parameter when via is specified
			line = fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = %s, underlying-proxy = %s",
				nodeName, entry.IP, entry.Port, entry.PSK, entry.Version, via)
		} else {
			line = fmt.Sprintf("%s = snell, %s, %d, psk = %s, version = %s",
				nodeName, entry.IP, entry.Port, entry.PSK, entry.Version)
		}
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
		// Resolve domain to IP if needed and get IP information
		// Keep original domain/IP in database, only use resolved IP for getting geo info
		_, ipInfo, err := utils.GetIPInfoFromDomainOrIP(modifyReq.IP)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ApiResponse{
				Status:  "error",
				Message: fmt.Sprintf("Failed to resolve domain/IP or get IP info: %v", err),
			})
			return
		}

		// Use original domain/IP address (not resolved IP)
		setStatements = append(setStatements, fmt.Sprintf(" ip = $%d", paramIndex))
		args = append(args, modifyReq.IP)
		paramIndex++

		// Update geolocation info
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
