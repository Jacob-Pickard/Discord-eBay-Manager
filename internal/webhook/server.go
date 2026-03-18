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
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Server handles incoming eBay webhook notifications
type Server struct {
	discord     *discordgo.Session
	channelID   string
	verifyToken string
	port        string
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

	// Setup OAuth callback handlers
	s.SetupOAuthHandlers()

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
	log.Printf("[DEBUG] Host header: %s", r.Host)
	log.Printf("[DEBUG] Full request URL: %s", r.URL.String())

	// Construct full endpoint URL for hash calculation
	// eBay expects: SHA256(challengeCode + verificationToken + endpointUrl)
	scheme := "https" // Always use https for eBay challenge verification
	host := r.Host
	if host == "" {
		host = "jacob.it.com"
	}
	endpointURL := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path)

	log.Printf("🔐 Computing challenge response for endpoint: %s", endpointURL)
	log.Printf("[DEBUG] Webhook Verification Inputs:\n  challengeCode: %s\n  verifyToken: %s\n  endpointURL: %s", challengeCode, s.verifyToken, endpointURL)

	// Create challenge response with SHA-256 hash
	hash := sha256.New()
	hash.Write([]byte(challengeCode))
	hash.Write([]byte(s.verifyToken))
	hash.Write([]byte(endpointURL))

	challengeResponse := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	log.Printf("[DEBUG] Computed challengeResponse (base64 SHA256): %s", challengeResponse)

	// Return JSON response
	response := map[string]string{
		"challengeResponse": challengeResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ Challenge response sent successfully")
}

// handleNotification processes incoming eBay notifications and challenges
func (s *Server) handleNotification(w http.ResponseWriter, r *http.Request) {
	// Handle GET requests as eBay verification challenges
	if r.Method == http.MethodGet {
		challengeCode := r.URL.Query().Get("challenge_code")
		if challengeCode == "" {
			// eBay might check the endpoint without challenge_code first
			// Return 200 OK to indicate the endpoint is live
			log.Printf("📨 Received GET request without challenge_code (possibly eBay preliminary check)")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "eBay Webhook Endpoint - Ready")
			return
		}

		log.Printf("📨 Received eBay challenge on notification endpoint: %s", challengeCode)

		// Construct full endpoint URL for hash calculation
		// eBay expects: SHA256(challengeCode + verificationToken + endpointUrl)
		scheme := "https" // Always use https for eBay challenge verification
		host := r.Host
		if host == "" {
			host = "jacob.it.com"
		}
		endpointURL := fmt.Sprintf("%s://%s%s", scheme, host, r.URL.Path)

		log.Printf("🔐 Computing challenge response for endpoint: %s", endpointURL)
		log.Printf("[DEBUG] Webhook Verification Inputs:\n  challengeCode: %s\n  verifyToken: %s\n  endpointURL: %s", challengeCode, s.verifyToken, endpointURL)

		// Create challenge response with SHA-256 hash
		hash := sha256.New()
		hash.Write([]byte(challengeCode))
		hash.Write([]byte(s.verifyToken))
		hash.Write([]byte(endpointURL))

		challengeResponse := base64.StdEncoding.EncodeToString(hash.Sum(nil))
		log.Printf("[DEBUG] Computed challengeResponse (base64 SHA256): %s", challengeResponse)

		// Return JSON response
		response := map[string]string{
			"challengeResponse": challengeResponse,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)

		log.Printf("✅ Challenge response sent successfully from notification endpoint")
		return
	}

	// Handle POST requests as actual notifications
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

	// Parse notification type (eBay uses format: MARKETPLACE_ACCOUNT.ORDER.FULFILLED)
	notifType := notification.NotificationEventType

	// Handle ORDER notifications
	if contains(notifType, "ORDER") {
		s.handleOrderNotification(notification, embed)
	} else if contains(notifType, "OFFER") {
		s.handleOfferNotification(notification, embed)
	} else {
		// Generic notification
		s.handleGenericNotification(notification, embed)
	}

	return embed
}

// handleOrderNotification processes order-related notifications
func (s *Server) handleOrderNotification(notification *EbayNotification, embed *discordgo.MessageEmbed) {
	notifType := notification.NotificationEventType

	// Determine specific order event
	if contains(notifType, "FULFILLED") || contains(notifType, "PLACED") {
		embed.Color = 0x00ff00
		embed.Title = "🎉 New Order Received!"
		embed.Description = "You have a new sale!"
	} else if contains(notifType, "PAID") || contains(notifType, "PAYMENT") {
		embed.Color = 0x00ff00
		embed.Title = "💵 Payment Received"
		embed.Description = "Payment has been confirmed for an order"
	} else if contains(notifType, "SHIPPED") {
		embed.Color = 0x3498db
		embed.Title = "📦 Order Shipped"
		embed.Description = "An order has been marked as shipped"
	} else {
		embed.Color = 0x0099ff
		embed.Title = "📋 Order Update"
		embed.Description = "Order status has changed"
	}

	// Extract order details from metadata
	if notification.Metadata != nil {
		if orderId, ok := notification.Metadata["orderId"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "📦 Order ID",
				Value:  orderId,
				Inline: true,
			})
		}
		if buyer, ok := notification.Metadata["buyerUsername"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "👤 Buyer",
				Value:  buyer,
				Inline: true,
			})
		}
		if price, ok := notification.Metadata["totalPrice"].(float64); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "💰 Total",
				Value:  fmt.Sprintf("$%.2f", price),
				Inline: true,
			})
		} else if priceStr, ok := notification.Metadata["totalPrice"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "💰 Total",
				Value:  priceStr,
				Inline: true,
			})
		}
		if itemTitle, ok := notification.Metadata["itemTitle"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "📦 Item",
				Value:  itemTitle,
				Inline: false,
			})
		}
	}

	embed.Footer = &discordgo.MessageEmbedFooter{
		Text: "💡 Use /get-orders to view all orders",
	}
}

// handleOfferNotification processes offer-related notifications
func (s *Server) handleOfferNotification(notification *EbayNotification, embed *discordgo.MessageEmbed) {
	notifType := notification.NotificationEventType

	// Determine specific offer event
	if contains(notifType, "CREATED") {
		embed.Color = 0xffaa00
		embed.Title = "💬 New Offer Received!"
		embed.Description = "A buyer has submitted an offer on your listing"
	} else if contains(notifType, "ACCEPTED") {
		embed.Color = 0x00ff00
		embed.Title = "✅ Offer Accepted"
		embed.Description = "An offer has been accepted"
	} else if contains(notifType, "DECLINED") {
		embed.Color = 0xff0000
		embed.Title = "❌ Offer Declined"
		embed.Description = "An offer has been declined"
	} else if contains(notifType, "COUNTERED") || contains(notifType, "UPDATED") {
		embed.Color = 0xffaa00
		embed.Title = "💬 Offer Updated"
		embed.Description = "An offer has been countered or updated"
	} else if contains(notifType, "EXPIRED") {
		embed.Color = 0x808080
		embed.Title = "⏰ Offer Expired"
		embed.Description = "An offer has expired"
	} else {
		embed.Color = 0xffaa00
		embed.Title = "💬 Offer Update"
		embed.Description = "Offer status has changed"
	}

	// Extract offer details from metadata
	if notification.Metadata != nil {
		if offerId, ok := notification.Metadata["offerId"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "🆔 Offer ID",
				Value:  fmt.Sprintf("`%s`", offerId),
				Inline: false,
			})
		}
		if buyer, ok := notification.Metadata["buyerUsername"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "👤 Buyer",
				Value:  buyer,
				Inline: true,
			})
		}
		if offerPrice, ok := notification.Metadata["offerPrice"].(float64); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "💰 Offer Amount",
				Value:  fmt.Sprintf("$%.2f", offerPrice),
				Inline: true,
			})
		} else if offerPriceStr, ok := notification.Metadata["offerPrice"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "💰 Offer Amount",
				Value:  offerPriceStr,
				Inline: true,
			})
		}
		if listPrice, ok := notification.Metadata["listPrice"].(float64); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "🏷️ List Price",
				Value:  fmt.Sprintf("$%.2f", listPrice),
				Inline: true,
			})
		}
		if itemTitle, ok := notification.Metadata["itemTitle"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "📦 Item",
				Value:  itemTitle,
				Inline: false,
			})
		}
		if itemId, ok := notification.Metadata["itemId"].(string); ok {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   "🔗 Item ID",
				Value:  itemId,
				Inline: true,
			})
		}
	}

	// Add action buttons guidance if it's a new/updated offer
	if contains(notifType, "CREATED") || contains(notifType, "UPDATED") {
		if offerId, ok := notification.Metadata["offerId"].(string); ok {
			embed.Footer = &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("💡 Respond with: /accept-offer, /counter-offer, or /decline-offer using offer ID: %s", offerId),
			}
		}
	}
}

// handleGenericNotification processes other notification types
func (s *Server) handleGenericNotification(notification *EbayNotification, embed *discordgo.MessageEmbed) {
	notifType := notification.NotificationEventType

	if contains(notifType, "ITEM_SOLD") {
		embed.Color = 0x00ff00
		embed.Title = "💰 Item Sold!"
		embed.Description = "You have a new sale!"
	} else if contains(notifType, "LISTING_ENDED") {
		embed.Color = 0xff6600
		embed.Title = "⏰ Listing Ended"
		embed.Description = "One of your listings has ended"

	} else {
		embed.Description = fmt.Sprintf("Event: %s", notification.NotificationEventType)
	}

	// Add any remaining metadata as fields
	if notification.Metadata != nil {
		for key, value := range notification.Metadata {
			// Skip already displayed fields
			if key == "orderId" || key == "offerId" || key == "buyerUsername" ||
				key == "totalPrice" || key == "offerPrice" || key == "listPrice" ||
				key == "itemTitle" || key == "itemId" {
				continue
			}
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:   key,
				Value:  fmt.Sprintf("%v", value),
				Inline: true,
			})
		}
	}
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(strings.ToUpper(s), strings.ToUpper(substr))
}

// EbayNotification represents an eBay webhook notification
type EbayNotification struct {
	NotificationEventType string                 `json:"notificationEventType"`
	NotificationId        string                 `json:"notificationId"`
	PublishDate           string                 `json:"publishDate"`
	Metadata              map[string]interface{} `json:"metadata"`
	Notification          map[string]interface{} `json:"notification"` // Alternative structure
}
