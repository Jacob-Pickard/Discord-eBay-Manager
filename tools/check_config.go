package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Simple config checker to verify .env setup without exposing secrets
func main() {
	fmt.Println("🔍 eBay Manager Bot - Configuration Checker")
	fmt.Println("==========================================")

	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("❌ Error: .env file not found!")
		fmt.Println("   Create one by copying .env.example to .env")
		os.Exit(1)
	}

	fmt.Println("✅ .env file found")
	fmt.Println()

	// Check required variables
	checks := map[string]string{
		"DISCORD_BOT_TOKEN":       os.Getenv("DISCORD_BOT_TOKEN"),
		"EBAY_APP_ID":             os.Getenv("EBAY_APP_ID"),
		"EBAY_CERT_ID":            os.Getenv("EBAY_CERT_ID"),
		"EBAY_DEV_ID":             os.Getenv("EBAY_DEV_ID"),
		"EBAY_REDIRECT_URI":       os.Getenv("EBAY_REDIRECT_URI"),
		"EBAY_ENVIRONMENT":        os.Getenv("EBAY_ENVIRONMENT"),
		"WEBHOOK_PORT":            os.Getenv("WEBHOOK_PORT"),
		"WEBHOOK_VERIFY_TOKEN":    os.Getenv("WEBHOOK_VERIFY_TOKEN"),
		"NOTIFICATION_CHANNEL_ID": os.Getenv("NOTIFICATION_CHANNEL_ID"),
	}

	optionalChecks := map[string]string{
		"EBAY_ACCESS_TOKEN":  os.Getenv("EBAY_ACCESS_TOKEN"),
		"EBAY_REFRESH_TOKEN": os.Getenv("EBAY_REFRESH_TOKEN"),
	}

	allGood := true

	// Check required configs
	fmt.Println("📋 Required Configuration:")
	for key, value := range checks {
		if value == "" || strings.Contains(value, "your_") || strings.Contains(value, "change_this") {
			fmt.Printf("   ❌ %s - NOT SET or uses placeholder\n", key)
			allGood = false
		} else {
			// Show first few characters only
			preview := value
			if len(value) > 10 {
				preview = value[:10] + "..."
			}
			fmt.Printf("   ✅ %s - %s\n", key, preview)
		}
	}

	// Check optional configs (OAuth tokens)
	fmt.Println("\n🔐 OAuth Tokens (Generated via /ebay-authorize):")
	for key, value := range optionalChecks {
		if value == "" {
			fmt.Printf("   ⚪ %s - Not generated yet (run /ebay-authorize)\n", key)
		} else {
			preview := value
			if len(value) > 10 {
				preview = value[:10] + "..."
			}
			fmt.Printf("   ✅ %s - %s\n", key, preview)
		}
	}

	// Environment check
	fmt.Println("\n🌍 Environment:")
	env := os.Getenv("EBAY_ENVIRONMENT")
	if env == "SANDBOX" {
		fmt.Println("   ✅ SANDBOX mode (safe for testing)")
	} else if env == "PRODUCTION" {
		fmt.Println("   ⚠️  PRODUCTION mode (real transactions!)")
	} else {
		fmt.Println("   ❌ Invalid environment! Must be SANDBOX or PRODUCTION")
		allGood = false
	}

	// Final status
	fmt.Println("\n" + strings.Repeat("=", 42))
	if allGood {
		fmt.Println("✅ Configuration looks good!")
		fmt.Println("\nNext steps:")
		fmt.Println("1. Run: .\\ebaymanager-bot.exe")
		fmt.Println("2. Go to Discord and run: /ebay-authorize")
		fmt.Println("3. Follow the OAuth flow to get tokens")
		fmt.Println("4. Start testing with SANDBOX_TEST_CHECKLIST.md")
	} else {
		fmt.Println("❌ Configuration incomplete!")
		fmt.Println("\nPlease update your .env file with actual values.")
		fmt.Println("See .env.example for reference.")
		os.Exit(1)
	}
}
