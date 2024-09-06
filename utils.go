/*
 * @Author: Vincent Yang
 * @Date: 2024-09-06 15:10:16
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-09-06 15:29:28
 * @FilePath: /snell-panel/utils.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2024 by Vincent, All Rights Reserved.
 */
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

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

func getIPInfo(ip string) (GeoIP, error) {
	url := fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GeoIP{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return GeoIP{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GeoIP{}, err
	}

	var geoIP GeoIP
	err = json.Unmarshal(body, &geoIP)
	if err != nil {
		return GeoIP{}, err
	}

	return geoIP, nil
}

func generateRandomString() string {
	rand.Seed(uint64(time.Now().UnixNano())) // Convert int64 to uint64
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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
