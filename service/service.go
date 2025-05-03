/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:24:26
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:24:33
 * @FilePath: /snell-panel/service/service.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package service

import (
	"database/sql"
	"fmt"
	"strings"

	"snell-panel/config"
	"snell-panel/database"
	"snell-panel/handlers"
	"snell-panel/models"
	"snell-panel/utils"

	"github.com/gin-gonic/gin"
)

// Service represents the main service layer
type Service struct {
	DB     *sql.DB
	Config *config.Config
}

// InitConfig initializes and returns the application configuration
func InitConfig() *config.Config {
	return config.LoadConfig()
}

// NewService creates a new service instance
func NewService(cfg *config.Config) *Service {
	db := database.InitDB(cfg.DatabaseURL)
	return &Service{
		DB:     db,
		Config: cfg,
	}
}

// Router initializes and returns the gin router
func Router(cfg *config.Config) *gin.Engine {
	// Set gin mode based on environment
	if cfg.IsDevelopment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create service
	svc := NewService(cfg)

	// Create handlers with the service
	h := handlers.NewHandlers(svc.DB, cfg.ApiToken)

	// Initialize router
	r := gin.Default()

	// Add CORS middleware
	r.Use(handlers.CorsMiddleware())

	// Routes
	r.GET("/", h.Welcome)
	r.POST("/entry", h.AuthMiddleware(), h.InsertEntry)
	r.GET("/entries", h.AuthMiddleware(), h.QueryAllEntries)
	r.DELETE("/entry/:ip", h.AuthMiddleware(), h.DeleteEntryByIP)
	r.GET("/subscribe", h.AuthMiddleware(), h.GetSubscription)
	r.PUT("/modify/:id", h.AuthMiddleware(), h.ModifyNodeByNodeID)
	r.NoRoute(h.NotFound)

	return r
}

// InsertEntry inserts a new entry into the database
func (s *Service) InsertEntry(entry *models.Entry) (*models.Entry, error) {
	// Get IP information
	ipInfo, err := utils.GetIPInfo(entry.IP)
	if err != nil {
		return nil, fmt.Errorf("failed to get IP info: %w", err)
	}

	// Update entry with IP information
	entry.CountryCode = ipInfo.CountryCode
	entry.ISP = ipInfo.ISP
	entry.ASN = ipInfo.ASN
	entry.NodeID = utils.GenerateUUID()

	// Insert entry into database
	var id int
	err = s.DB.QueryRow(`
		 INSERT INTO entries (ip, port, psk, country_code, isp, asn, node_id, node_name) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		 RETURNING id`,
		entry.IP, entry.Port, entry.PSK, entry.CountryCode, entry.ISP, entry.ASN, entry.NodeID, entry.NodeName).Scan(&id)
	if err != nil {
		return nil, err
	}

	// Set ID in entry
	entry.ID = id

	return entry, nil
}

// DeleteEntryByIP deletes an entry by IP address
func (s *Service) DeleteEntryByIP(ip string) error {
	result, err := s.DB.Exec("DELETE FROM entries WHERE ip = $1", ip)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("entry not found")
	}

	return nil
}

// QueryAllEntries retrieves all entries from the database
func (s *Service) QueryAllEntries() ([]models.Entry, error) {
	rows, err := s.DB.Query(`
		 SELECT id, ip, port, psk, country_code, isp, asn, node_id, node_name 
		 FROM entries
	 `)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

// GetSubscription generates a subscription string for all entries
func (s *Service) GetSubscription() (string, error) {
	rows, err := s.DB.Query(`
		 SELECT ip, port, psk, country_code, isp, asn, node_id, node_name 
		 FROM entries
	 `)
	if err != nil {
		return "", err
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
			return "", err
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
		return "", fmt.Errorf("no entries found for subscription")
	}

	return strings.Join(subscriptionLines, "\n"), nil
}

// ModifyNodeByNodeID modifies node name and/or IP by node ID
func (s *Service) ModifyNodeByNodeID(nodeID string, modifyReq *models.ModifyRequest) error {
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
		return fmt.Errorf("no fields to update")
	}

	// Combine all set statements
	query += strings.Join(setStatements, ",")
	query += fmt.Sprintf(" WHERE node_id = $%d", paramIndex)
	args = append(args, nodeID)

	// Execute update
	result, err := s.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("node ID not found")
	}

	return nil
}
