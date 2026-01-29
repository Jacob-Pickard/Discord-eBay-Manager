package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Server handles incoming eBay webhook notifications
type Server struct {
	discord      *discordgo.Session
	channelID    string
	verifyToken  string
	port         string
}

// NewServer creates a new webhook server
func NewServer(discord *discordgo.Session, channelID, verifyToken, port string) *Server {
	return &Server{
		discord:     discord,
		channelID:   channelID,
		verifyToken: verifyToken,
		port:        port,
	}
}

// Start begins listening for webhook notifications
func (s *Server) Start() error {
	http.HandleFunc("/webhook/ebay/notification", s.handleNotification)
	http.HandleFunc("/webhook/ebay/challenge", s.handleChallenge)
	http.HandleFunc("/webhook/health", s.handleHealth)

	addr := ":" + s.port
	log.Printf("🎣 Webhook server starting on %s", addr)
	log.Printf("📍 Notification endpoint: http://localhost%s/webhook/ebay/notification", addr)
	log.Printf("📍 Challenge endpoint: http://localhost%s/webhook/ebay/challenge", addr)
	
	return http.ListenAndServe(addr, nil)
}

// handleChallenge responds to eBay's endpoint verification challenge
func (s *Server) handleChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the challenge code from query parameter
	challengeCode := r.URL.Query().Get("challenge_code")
	if challengeCode == "" {
		http.Error(w, "Missing challenge_code parameter", http.StatusBadRequest)
		return
	}

	log.Printf("📨 Received eBay challenge: %s", challengeCode)

	// Create challenge response with SHA-256 hash
	hash := sha256.New()
	hash.Write([]byte(challengeCode))
	hash.Write([]byte(s.verifyToken))
	hash.Write([]byte(r.URL.Path))
	
	challengeResponse := base64.StdEncoding.EncodeToString(hash.Sum(nil))

	// Return JSON response
	response := map[string]string{
		"challengeResponse": challengeResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	
	log.Printf("✅ Challenge response sent successfully")
}

// handleNotification processes incoming eBay notifications
func (s *Server) handleNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read notification body: %v", err)
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Verify signature (optional but recommended)
	signature := r.Header.Get("X-EBAY-SIGNATURE")
	if signature != "" && !s.verifySignature(body, signature) {
		log.Printf("❌ Invalid signature for notification")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse notification
	var notification EbayNotification
	if err := json.Unmarshal(body, &notification); err != nil {
		log.Printf("❌ Failed to parse notification: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("📨 Received eBay notification: %s", notification.NotificationEventType)

	// Process and send to Discord
	go s.processNotification(&notification)

	// Respond with 200 OK
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// handleHealth provides a health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}

// verifySignature verifies the eBay notification signature
func (s *Server) verifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(s.verifyToken))
	mac.Write(body)
	expectedMAC := mac.Sum(nil)
	expectedSignature := base64.StdEncoding.EncodeToString(expectedMAC)
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// processNotification handles the notification and sends to Discord
func (s *Server) processNotification(notification *EbayNotification) {
	if s.channelID == "" {
		log.Println("⚠️ No Discord channel configured for notifications")
		return
	}

	embed := s.buildDiscordEmbed(notification)
	
	_, err := s.discord.ChannelMessageSendEmbed(s.channelID, embed)
	if err != nil {
		log.Printf("❌ Failed to send Discord notification: %v", err)
		return
	}
	
	log.Printf("✅ Notification sent to Discord channel %s", s.channelID)
}

// buildDiscordEmbed creates a rich embed for the Discord notification
func (s *Server) buildDiscordEmbed(notification *EbayNotification) *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Title:     "🔔 eBay Notification",
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     0x0099ff,
	}

	switch notification.NotificationEventType {
	case "ITEM_SOLD":
		embed.Color = 0x00ff00
		embed.Title = "💰 Item Sold!"
		embed.Description = "You have a new sale!"
		
	case "ORDER_PAYMENT_RECEIVED":
		embed.Color = 0x00ff00
		embed.Title = "💵 Payment Received"
		embed.Description = "Payment has been received for an order"
		
	case "ORDER_SHIPPED":
		embed.Color = 0x3498db
		embed.Title = "📦 Order Shipped"
		embed.Description = "An order has been marked as shipped"
		
	case "OFFER_RECEIVED":
		embed.Color = 0xffaa00
		embed.Title = "💬 Offer Received"
		embed.Description = "You have received a new offer from a buyer"
		
	case "OFFER_COUNTERED":
		embed.Color = 0xffaa00
		embed.Title = "💬 Offer Countered"
		embed.Description = "A buyer has countered your offer"
		
	case "LISTING_ENDED":
		embed.Color = 0xff6600
		embed.Title = "⏰ Listing Ended"
		embed.Description = "One of your listings has ended"
		
	default:
		embed.Description = fmt.Sprintf("Event: %s", notification.NotificationEventType)
	}

	// Add notification details as fields
	if notification.Metadata != nil {
		for key, value := range notification.Metadata {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   key,
				Value:  fmt.Sprintf("%v", value),
				Inline: true,
			})
		}
	}

	return embed
}

// EbayNotification represents an eBay webhook notification
type EbayNotification struct {
	NotificationEventType string                 `json:"notificationEventType"`
	NotificationId        string                 `json:"notificationId"`
	PublishDate           string                 `json:"publishDate"`
	Metadata              map[string]interface{} `json:"metadata"`
}
