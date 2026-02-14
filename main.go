package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ebaymanager-bot/internal/bot"
	"ebaymanager-bot/internal/config"
	"ebaymanager-bot/internal/ebay"
	"ebaymanager-bot/internal/webhook"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize eBay client
	ebayClient := ebay.NewClient(cfg.EbayConfig)

	// Create Discord session
	discord, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	log.Println("Connecting to Discord...")
	// Open Discord connection first
	if err := discord.Open(); err != nil {
		log.Fatalf("Failed to open Discord connection: %v", err)
	}
	defer discord.Close()

	log.Printf("Connected as: %s#%s (ID: %s)", discord.State.User.Username, discord.State.User.Discriminator, discord.State.User.ID)

	// Start webhook server in background first
	webhookServer := webhook.NewServer(discord, cfg.NotificationChannelID, cfg.WebhookVerifyToken, cfg.WebhookPort)
	webhook.SetEbayClient(ebayClient) // Set eBay client for OAuth (package-level)
	go func() {
		if err := webhookServer.Start(); err != nil {
			log.Printf("⚠️ Webhook server error: %v", err)
		}
	}()

	// Initialize bot and register commands after connection is open
	botHandler := bot.NewHandler(discord, ebayClient)
	botHandler.SetWebhookServer(webhookServer) // Pass webhook server for OAuth
	botHandler.RegisterCommands()

	fmt.Println("eBay Manager Bot is now running. Press CTRL+C to exit.")

	// Wait for interrupt signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	fmt.Println("\nShutting down gracefully...")
}
