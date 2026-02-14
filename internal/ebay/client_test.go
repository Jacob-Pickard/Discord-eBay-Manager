package ebay

import (
	"strings"
	"testing"

	"ebaymanager-bot/internal/config"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		config     config.EbayConfig
		expectURL  string
		expectAuth string
	}{
		{
			name: "Sandbox Environment",
			config: config.EbayConfig{
				Environment: "SANDBOX",
				AppID:       "test-app-id",
				CertID:      "test-cert-id",
			},
			expectURL:  sandboxAPIURL,
			expectAuth: sandboxAuthURL,
		},
		{
			name: "Production Environment",
			config: config.EbayConfig{
				Environment: "PRODUCTION",
				AppID:       "test-app-id",
				CertID:      "test-cert-id",
			},
			expectURL:  productionAPIURL,
			expectAuth: productionAuthURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)

			if client.baseURL != tt.expectURL {
				t.Errorf("Expected baseURL %s, got %s", tt.expectURL, client.baseURL)
			}

			if client.authURL != tt.expectAuth {
				t.Errorf("Expected authURL %s, got %s", tt.expectAuth, client.authURL)
			}
		})
	}
}

func TestCheckConnection(t *testing.T) {
	tests := []struct {
		name        string
		accessToken string
		expectError bool
	}{
		{
			name:        "No Access Token",
			accessToken: "",
			expectError: true,
		},
		{
			name:        "Valid Access Token",
			accessToken: "v^1.1#i^1#test-token",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.EbayConfig{
				Environment: "SANDBOX",
				AccessToken: tt.accessToken,
			}
			client := NewClient(cfg)

			result := client.CheckConnection()

			if tt.expectError && !strings.Contains(result, "No access token") {
				t.Errorf("Expected error message about no access token, got: %s", result)
			}

			if !tt.expectError && !strings.Contains(result, "Connected to eBay API") {
				t.Errorf("Expected success message, got: %s", result)
			}
		})
	}
}

func TestGetOrdersNoToken(t *testing.T) {
	cfg := config.EbayConfig{
		Environment: "SANDBOX",
		AccessToken: "",
	}
	client := NewClient(cfg)

	_, err := client.GetOrders(10)
	if err == nil {
		t.Error("Expected error when no access token provided")
	}
}

func TestGetOffersNoToken(t *testing.T) {
	cfg := config.EbayConfig{
		Environment: "SANDBOX",
		AccessToken: "",
	}
	client := NewClient(cfg)

	_, err := client.GetOffers()
	if err == nil {
		t.Error("Expected error when no access token provided")
	}
}
