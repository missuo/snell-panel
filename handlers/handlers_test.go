package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

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

func TestFormatMihomoSubscriptionLine(t *testing.T) {
	entry := models.Entry{
		IP:      "example.com",
		Port:    443,
		PSK:     "your_psk_here",
		Version: "5",
		TFO:     true,
	}

	got := formatMihomoSubscriptionLine(entry, "Custom Node Name", "Parent Proxy")
	want := `  - {name: "Custom Node Name", server: "example.com", port: 443, type: snell, psk: "your_psk_here", version: 5, tfo: true, dialer-proxy: "Parent Proxy"}`

	if got != want {
		t.Fatalf("formatMihomoSubscriptionLine() = %q, want %q", got, want)
	}
}

func TestFormatMihomoSubscriptionContent(t *testing.T) {
	lines := []string{`  - {name: "Custom Node Name", server: "example.com", port: 443, type: snell, psk: "your_psk_here", version: 5}`}

	got := formatSubscriptionContent(lines, subscriptionFormatMihomo)
	want := "proxies:\n" + lines[0]

	if got != want {
		t.Fatalf("formatSubscriptionContent() = %q, want %q", got, want)
	}
}

func TestMihomoSnellVersionIgnoresInvalidValues(t *testing.T) {
	if got := mihomoSnellVersion("5, injected: true"); got != "" {
		t.Fatalf("mihomoSnellVersion() = %q, want empty string", got)
	}
}

func TestParseSubscriptionFormatMihomo(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []string{
		"/subscribe?format=mihomo",
		"/subscribe?mihomo=true",
	}

	for _, target := range tests {
		t.Run(target, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Request = httptest.NewRequest("GET", target, nil)

			if got := parseSubscriptionFormat(c); got != subscriptionFormatMihomo {
				t.Fatalf("parseSubscriptionFormat() = %q, want %q", got, subscriptionFormatMihomo)
			}
		})
	}
}
