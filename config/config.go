/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:22:58
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:23:38
 * @FilePath: /snell-panel/config/config.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents application configuration
type Config struct {
	ApiToken      string
	DatabaseURL   string
	Port          int
	IsDevelopment bool
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Load API token
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		log.Fatal("API_TOKEN is required")
	}

	// Load database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Load port with default value of 8080
	portStr := os.Getenv("PORT")
	port := 8080
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Printf("Invalid PORT value: %s, using default: 8080", portStr)
			port = 8080
		}
	}

	// Check if we're in development mode
	isDev := os.Getenv("ENV") == "development"

	return &Config{
		ApiToken:      apiToken,
		DatabaseURL:   dbURL,
		Port:          port,
		IsDevelopment: isDev,
	}
}

// GetPortString returns the port string for HTTP serving
func (c *Config) GetPortString() string {
	return fmt.Sprintf(":%d", c.Port)
}
