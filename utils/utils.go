/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:23:55
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-05-03 04:24:00
 * @FilePath: /snell-panel/utils/utils.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"snell-panel/models"

	"github.com/google/uuid"
)

// GenerateUUID generates a random UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// GetIPInfo retrieves geolocation information for an IP address
func GetIPInfo(ip string) (models.GeoIP, error) {
	url := fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.GeoIP{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return models.GeoIP{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.GeoIP{}, err
	}

	var geoIP models.GeoIP
	err = json.Unmarshal(body, &geoIP)
	if err != nil {
		return models.GeoIP{}, err
	}

	return geoIP, nil
}

// CountryCodeToFlagEmoji converts a country code to a flag emoji
func CountryCodeToFlagEmoji(countryCode string) string {
	if len(countryCode) != 2 {
		return countryCode // Return original string if it's not a 2-letter code
	}

	// Convert ASCII to regional indicator symbols
	regionalIndicatorA := rune(0x1F1E6)
	flagEmoji := ""

	for _, char := range strings.ToUpper(countryCode) {
		if char < 'A' || char > 'Z' {
			return countryCode // Return original string if it contains non-letter characters
		}
		flagEmoji += string(regionalIndicatorA + rune(char) - 'A')
	}

	return flagEmoji
}
