/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:25:37
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:25:40
 * @FilePath: /snell-panel/api/entrypoint.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package api

import (
	"net/http"

	"snell-panel/config"
	"snell-panel/service"

	"github.com/gin-gonic/gin"
)

var (
	cfg *config.Config
	app *gin.Engine
)

func init() {
	// Initialize the configuration
	cfg = config.LoadConfig()

	// Initialize the router
	app = service.Router(cfg)
}

// Entrypoint is the serverless function handler for Vercel
func Entrypoint(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

// Type assertion to ensure our function matches http.HandlerFunc
var _ http.HandlerFunc = Entrypoint
