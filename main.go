/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:25:17
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:25:20
 * @FilePath: /snell-panel/main.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package main

import (
	"log"

	"snell-panel/config"
	"snell-panel/database"
	"snell-panel/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := database.InitDB(cfg.DatabaseURL)
	defer database.CloseDB(db)

	// Initialize router
	router := service.Router(cfg)

	// Start server
	log.Printf("Server starting on port %d...", cfg.Port)
	if err := router.Run(cfg.GetPortString()); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
