/*
 * @Author: Vincent Yang
 * @Date: 2025-05-03 04:23:55
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-07-05 20:43:56
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
	"net"
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

// IsValidIP checks if a string is a valid IP address
func IsValidIP(addr string) bool {
	return net.ParseIP(addr) != nil
}

// ResolveDomainToIP resolves a domain name to an IP address
func ResolveDomainToIP(domain string) (string, error) {
	// Check if it's already an IP address
	if IsValidIP(domain) {
		return domain, nil
	}

	// Resolve domain to IP addresses
	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", fmt.Errorf("failed to resolve domain %s: %w", domain, err)
	}

	// Find the first IPv4 address
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	// If no IPv4 found, try IPv6
	for _, ip := range ips {
		if ipv6 := ip.To16(); ipv6 != nil {
			return ipv6.String(), nil
		}
	}

	return "", fmt.Errorf("no valid IP address found for domain %s", domain)
}

// GetIPInfoFromDomainOrIP resolves domain to IP if needed and gets geolocation info
func GetIPInfoFromDomainOrIP(domainOrIP string) (string, models.GeoIP, error) {
	// Resolve domain to IP if necessary
	resolvedIP, err := ResolveDomainToIP(domainOrIP)
	if err != nil {
		return "", models.GeoIP{}, err
	}

	// Get IP geolocation information
	geoIP, err := GetIPInfo(resolvedIP)
	if err != nil {
		return resolvedIP, models.GeoIP{}, err
	}

	return resolvedIP, geoIP, nil
}
