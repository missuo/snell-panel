/*
 * @Author: Vincent Yang
 * @Date: 2024-09-06 15:10:16
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2024-09-06 15:21:18
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
	"time"

	"golang.org/x/exp/rand"
)

type GeoIP struct {
	Organization    string  `json:"organization"`
	Longitude       float64 `json:"longitude"`
	Timezone        string  `json:"timezone"`
	ISP             string  `json:"isp"`
	Offset          int     `json:"offset"`
	ASN             int     `json:"asn"`
	ASNOrganization string  `json:"asn_organization"`
	Country         string  `json:"country"`
	IP              string  `json:"ip"`
	Latitude        float64 `json:"latitude"`
	ContinentCode   string  `json:"continent_code"`
	CountryCode     string  `json:"country_code"`
}

func getIPInfo(ip string) (GeoIP, error) {
	url := fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)
	resp, err := http.Get(url)
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
