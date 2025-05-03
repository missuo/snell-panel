/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:22:16
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:23:34
 * @FilePath: /snell-panel/models/models.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package models

// Entry represents a snell proxy entry
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

// ModifyRequest represents a request to modify an entry
type ModifyRequest struct {
	NodeName string `json:"node_name,omitempty"`
	IP       string `json:"ip,omitempty"`
}

// GeoIP represents IP geolocation information
type GeoIP struct {
	Organization    string `json:"organization"`
	ISP             string `json:"isp"`
	ASN             int    `json:"asn"`
	ASNOrganization string `json:"asn_organization"`
	Country         string `json:"country"`
	IP              string `json:"ip"`
	ContinentCode   string `json:"continent_code"`
	CountryCode     string `json:"country_code"`
}

// ApiResponse represents a standardized API response
type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
