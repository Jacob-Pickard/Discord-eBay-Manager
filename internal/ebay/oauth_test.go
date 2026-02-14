package ebay

import (
	"encoding/json"
	"os"
	"testing"

	"ebaymanager-bot/internal/config"
)

func TestExchangeCodeForToken(t *testing.T) {
	// Skip if no credentials (CI/CD environment)
	if os.Getenv("EBAY_APP_ID") == "" {
		t.Skip("Skipping OAuth test: no credentials configured")
	}

	cfg := config.EbayConfig{
		AppID:       os.Getenv("EBAY_APP_ID"),
		CertID:      os.Getenv("EBAY_CERT_ID"),
		Environment: "SANDBOX",
		RedirectURI: "https://localhost:3000/callback",
	}

	client := NewClient(cfg)

	// Test with invalid code (should fail gracefully)
	_, err := client.ExchangeCodeForToken("invalid_code_12345")
	if err == nil {
		t.Error("Expected error with invalid authorization code")
	}

	// Error should be descriptive
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}

func TestRefreshAccessToken(t *testing.T) {
	cfg := config.EbayConfig{
		AppID:       "test_app_id",
		CertID:      "test_cert_id",
		Environment: "SANDBOX",
	}

	client := NewClient(cfg)

	// Test with empty refresh token
	_, err := client.RefreshAccessToken("")
	if err == nil {
		t.Error("Expected error when refresh token is empty")
	}
}

func TestGetUserAuthorizationURL(t *testing.T) {
	cfg := config.EbayConfig{
		AppID:       "test-app-id",
		Environment: "SANDBOX",
		RedirectURI: "https://localhost:3000/callback",
	}

	client := NewClient(cfg)
	authURL := client.GetUserAuthorizationURL("test-state")

	// Verify URL contains required components
	if authURL == "" {
		t.Error("Authorization URL should not be empty")
	}

	// Should contain client_id
	if !contains(authURL, "client_id=test-app-id") {
		t.Error("Authorization URL should contain app ID")
	}

	// Should contain redirect_uri
	if !contains(authURL, "redirect_uri=") {
		t.Error("Authorization URL should contain redirect URI")
	}

	// Should contain response_type=code
	if !contains(authURL, "response_type=code") {
		t.Error("Authorization URL should contain response_type=code")
	}

	// Should use sandbox URL
	if !contains(authURL, "auth.sandbox.ebay.com") {
		t.Error("Should use sandbox auth URL")
	}
}

func TestGetUserAuthorizationURLProduction(t *testing.T) {
	cfg := config.EbayConfig{
		AppID:       "test-app-id",
		Environment: "PRODUCTION",
		RedirectURI: "https://localhost:3000/callback",
	}

	client := NewClient(cfg)
	authURL := client.GetUserAuthorizationURL("prod-state")

	// Should use production URL
	if !contains(authURL, "auth.ebay.com") {
		t.Error("Should use production auth URL")
	}

	// Should NOT contain sandbox
	if contains(authURL, "sandbox") {
		t.Error("Production URL should not contain 'sandbox'")
	}
}

func TestTokenResponseParsing(t *testing.T) {
	tokenJSON := `{
		"access_token": "v^1.1#i^1#test_token",
		"expires_in": 7200,
		"refresh_token": "v^1.1#i^1#refresh_token",
		"refresh_token_expires_in": 47304000,
		"token_type": "User Access Token"
	}`

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
	}

	err := json.Unmarshal([]byte(tokenJSON), &tokenResp)
	if err != nil {
		t.Fatalf("Failed to parse token response: %v", err)
	}

	if tokenResp.AccessToken != "v^1.1#i^1#test_token" {
		t.Error("Access token not parsed correctly")
	}

	if tokenResp.ExpiresIn != 7200 {
		t.Error("ExpiresIn not parsed correctly")
	}

	if tokenResp.RefreshToken != "v^1.1#i^1#refresh_token" {
		t.Error("Refresh token not parsed correctly")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
