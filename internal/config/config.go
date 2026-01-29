package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	DiscordToken      string
	EbayConfig        EbayConfig
	WebhookPort       string
	WebhookVerifyToken string
	NotificationChannelID string
}

// EbayConfig holds eBay API configuration
type EbayConfig struct {
	AppID       string
	CertID      string
	DevID       string
	RedirectURI string
	AccessToken string
	RefreshToken string
	Environment string // PRODUCTION or SANDBOX
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	if discordToken == "" {
		return nil, fmt.Errorf("DISCORD_BOT_TOKEN not set")
	}

	ebayAppID := os.Getenv("EBAY_APP_ID")
	if ebayAppID == "" {
		return nil, fmt.Errorf("EBAY_APP_ID not set")
	}

	ebayEnvironment := os.Getenv("EBAY_ENVIRONMENT")
	if ebayEnvironment == "" {
		ebayEnvironment = "SANDBOX"
	}

	webhookPort := os.Getenv("WEBHOOK_PORT")
	if webhookPort == "" {
		webhookPort = "8080"
	}

	webhookVerifyToken := os.Getenv("WEBHOOK_VERIFY_TOKEN")
	if webhookVerifyToken == "" {
		webhookVerifyToken = "default_verify_token_change_me"
	}

	return &Config{
		DiscordToken: discordToken,
		EbayConfig: EbayConfig{
			AppID:        ebayAppID,
			CertID:       os.Getenv("EBAY_CERT_ID"),
			DevID:        os.Getenv("EBAY_DEV_ID"),
			RedirectURI:  os.Getenv("EBAY_REDIRECT_URI"),
			AccessToken:  os.Getenv("EBAY_ACCESS_TOKEN"),
			RefreshToken: os.Getenv("EBAY_REFRESH_TOKEN"),
			Environment:  ebayEnvironment,
		},
		WebhookPort:           webhookPort,
		WebhookVerifyToken:    webhookVerifyToken,
		NotificationChannelID: os.Getenv("NOTIFICATION_CHANNEL_ID"),
	}, nil
}
