package main

import (
	"testing"
)

func TestAPIUrl(t *testing.T) {
	if gApiUrl != "https://api.qrclip.io" {
		t.Errorf("API URL BADLY CONFIGURED")
	}

	if gSpaUrl != "https://app.qrclip.io" {
		t.Errorf("SPA URL BADLY CONFIGURED")
	}
}
