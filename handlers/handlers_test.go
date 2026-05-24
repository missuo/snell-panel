package handlers

import (
	"testing"

	"snell-panel/models"
)

func TestFormatShadowrocketSubscriptionLine(t *testing.T) {
	entry := models.Entry{
		IP:      "example.com",
		Port:    443,
		PSK:     "your_psk_here",
		Version: "5",
		TFO:     true,
	}

	got := formatShadowrocketSubscriptionLine(entry, "")
	want := "snell://Y2hhY2hhMjAtaWV0Zi1wb2x5MTMwNTp5b3VyX3Bza19oZXJlQGV4YW1wbGUuY29tOjQ0Mw?tfo=1&version=5"

	if got != want {
		t.Fatalf("formatShadowrocketSubscriptionLine() = %q, want %q", got, want)
	}
}

func TestFormatSurgeSubscriptionLine(t *testing.T) {
	entry := models.Entry{
		IP:      "example.com",
		Port:    443,
		PSK:     "your_psk_here",
		Version: "5",
		TFO:     true,
	}

	got := formatSurgeSubscriptionLine(entry, "Custom Node Name", "")
	want := "Custom Node Name = snell, example.com, 443, psk = your_psk_here, version = 5, tfo = true"

	if got != want {
		t.Fatalf("formatSurgeSubscriptionLine() = %q, want %q", got, want)
	}
}
