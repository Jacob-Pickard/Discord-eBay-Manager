package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Save original env vars
	originalVars := map[string]string{
		"DISCORD_BOT_TOKEN":       os.Getenv("DISCORD_BOT_TOKEN"),
		"EBAY_APP_ID":             os.Getenv("EBAY_APP_ID"),
		"EBAY_CERT_ID":            os.Getenv("EBAY_CERT_ID"),
		"EBAY_DEV_ID":             os.Getenv("EBAY_DEV_ID"),
		"EBAY_REDIRECT_URI":       os.Getenv("EBAY_REDIRECT_URI"),
		"EBAY_ENVIRONMENT":        os.Getenv("EBAY_ENVIRONMENT"),
		"NOTIFICATION_CHANNEL_ID": os.Getenv("NOTIFICATION_CHANNEL_ID"),
	}

	// Restore env vars after test
	defer func() {
		for key, value := range originalVars {
			if value != "" {
				os.Setenv(key, value)
			} else {
				os.Unsetenv(key)
			}
		}
	}()

	tests := []struct {
		name      string
		setupEnv  func()
		wantError bool
	}{
		{
			name: "Valid Configuration",
			setupEnv: func() {
				os.Setenv("DISCORD_BOT_TOKEN", "test_token")
				os.Setenv("EBAY_APP_ID", "test_app_id")
				os.Setenv("EBAY_CERT_ID", "test_cert_id")
				os.Setenv("EBAY_DEV_ID", "test_dev_id")
				os.Setenv("EBAY_REDIRECT_URI", "https://test.com")
				os.Setenv("EBAY_ENVIRONMENT", "SANDBOX")
				os.Setenv("NOTIFICATION_CHANNEL_ID", "123456789")
			},
			wantError: false,
		},
		{
			name: "Missing Discord Token",
			setupEnv: func() {
				os.Unsetenv("DISCORD_BOT_TOKEN")
				os.Setenv("EBAY_APP_ID", "test_app_id")
				os.Setenv("EBAY_CERT_ID", "test_cert_id")
				os.Setenv("EBAY_DEV_ID", "test_dev_id")
				os.Setenv("EBAY_REDIRECT_URI", "https://test.com")
				os.Setenv("EBAY_ENVIRONMENT", "SANDBOX")
			},
			wantError: true,
		},
		{
			name: "Missing eBay App ID",
			setupEnv: func() {
				os.Setenv("DISCORD_BOT_TOKEN", "test_token")
				os.Unsetenv("EBAY_APP_ID")
				os.Setenv("EBAY_CERT_ID", "test_cert_id")
				os.Setenv("EBAY_DEV_ID", "test_dev_id")
				os.Setenv("EBAY_REDIRECT_URI", "https://test.com")
				os.Setenv("EBAY_ENVIRONMENT", "SANDBOX")
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()

			cfg, err := Load()

			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.wantError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.wantError && cfg == nil {
				t.Error("Expected config but got nil")
			}
		})
	}
}

func TestDefaultValues(t *testing.T) {
	// Setup minimal env
	os.Setenv("DISCORD_BOT_TOKEN", "test_token")
	os.Setenv("EBAY_APP_ID", "test_app_id")
	os.Setenv("EBAY_CERT_ID", "test_cert_id")
	os.Setenv("EBAY_DEV_ID", "test_dev_id")
	os.Setenv("EBAY_REDIRECT_URI", "https://test.com")
	os.Setenv("EBAY_ENVIRONMENT", "SANDBOX")
	os.Unsetenv("WEBHOOK_PORT")
	os.Unsetenv("NOTIFICATION_CHANNEL_ID")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.WebhookPort != "8081" {
		t.Errorf("Expected default webhook port 8081, got %s", cfg.WebhookPort)
	}
}
