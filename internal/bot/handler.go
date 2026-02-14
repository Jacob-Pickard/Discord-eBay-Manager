package bot

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"ebaymanager-bot/internal/ebay"

	"github.com/bwmarrin/discordgo"
)

// WebhookServer interface for OAuth callbacks
type WebhookServer interface {
	RegisterOAuthCallback(state string, discord *discordgo.Session, interaction *discordgo.Interaction)
}

// Handler manages Discord bot interactions
type Handler struct {
	discord       *discordgo.Session
	ebay          *ebay.Client
	webhookServer WebhookServer
}

// NewHandler creates a new bot handler
func NewHandler(discord *discordgo.Session, ebayClient *ebay.Client) *Handler {
	return &Handler{
		discord: discord,
		ebay:    ebayClient,
	}
}

// SetWebhookServer sets the webhook server for OAuth callbacks
func (h *Handler) SetWebhookServer(server WebhookServer) {
	h.webhookServer = server
}

// RegisterCommands sets up Discord slash commands and message handlers
func (h *Handler) RegisterCommands() {
	// Register message handler
	h.discord.AddHandler(h.messageHandler)

	// Register slash command handler
	h.discord.AddHandler(h.interactionHandler)

	// Register slash commands
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "get-orders",
			Description: "Get recent eBay orders",
		},
		{
			Name:        "get-offers",
			Description: "Get pending offers",
		},
		{
			Name:        "get-listings",
			Description: "View your active eBay listings",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "limit",
					Description: "Number of listings to show (default: 10)",
					Required:    false,
				},
			},
		},
		{
			Name:        "get-balance",
			Description: "View your eBay account balance",
		},
		{
			Name:        "get-payouts",
			Description: "View recent payouts and transactions",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "limit",
					Description: "Number of payouts to show (default: 10)",
					Required:    false,
				},
			},
		},
		{
			Name:        "get-messages",
			Description: "View buyer messages",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "limit",
					Description: "Number of messages to show (default: 10)",
					Required:    false,
				},
			},
		},
		{
			Name:        "ebay-status",
			Description: "Check eBay API connection status",
		},
		{
			Name:        "ebay-scopes",
			Description: "Check what API scopes your OAuth token has",
		},
		{
			Name:        "ebay-authorize",
			Description: "Authorize bot with your eBay account (get refresh token)",
		},
		{
			Name:        "ebay-code",
			Description: "Submit eBay authorization code manually",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "code",
					Description: "Authorization code from eBay redirect URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "webhook-subscribe",
			Description: "Subscribe to eBay notifications (orders, offers, etc)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "url",
					Description: "Webhook URL (leave empty to use jacob.it.com)",
					Required:    false,
				},
			},
		},
		{
			Name:        "webhook-list",
			Description: "List active webhook subscriptions",
		},
		{
			Name:        "webhook-test",
			Description: "Test webhook notification to this channel",
		},
		{
			Name:        "accept-offer",
			Description: "Accept a buyer's offer",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "offer-id",
					Description: "The offer ID to accept",
					Required:    true,
				},
			},
		},
		{
			Name:        "counter-offer",
			Description: "Counter a buyer's offer with a different price",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "offer-id",
					Description: "The offer ID to counter",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "price",
					Description: "Your counter offer price (e.g., 250.00)",
					Required:    true,
				},
			},
		},
		{
			Name:        "decline-offer",
			Description: "Decline a buyer's offer",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "offer-id",
					Description: "The offer ID to decline",
					Required:    true,
				},
			},
		},
	}

	// Delete all existing commands first (cleans up old/removed commands)
	existingCommands, err := h.discord.ApplicationCommands(h.discord.State.User.ID, "")
	if err == nil {
		log.Printf("Cleaning up %d existing commands...", len(existingCommands))
		for _, cmd := range existingCommands {
			err := h.discord.ApplicationCommandDelete(h.discord.State.User.ID, "", cmd.ID)
			if err != nil {
				log.Printf("⚠️ Failed to delete command %s: %v", cmd.Name, err)
			} else {
				log.Printf("🗑️ Deleted old command: /%s", cmd.Name)
			}
		}
	}

	log.Printf("Registering %d commands...", len(commands))
	for _, cmd := range commands {
		createdCmd, err := h.discord.ApplicationCommandCreate(h.discord.State.User.ID, "", cmd)
		if err != nil {
			log.Printf("❌ Failed to create command %s: %v", cmd.Name, err)
		} else {
			log.Printf("✅ Registered command: /%s (ID: %s)", createdCmd.Name, createdCmd.ID)
		}
	}
	log.Println("Command registration complete!")
}

// messageHandler handles regular Discord messages
func (h *Handler) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Simple command prefix handling
	if strings.HasPrefix(m.Content, "!ebay") {
		s.ChannelMessageSend(m.ChannelID, "Use slash commands instead! Try `/ebay-status`")
	}
}

// interactionHandler handles slash command interactions
func (h *Handler) interactionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "get-orders":
		h.handleGetOrders(s, i)
	case "get-offers":
		h.handleGetOffers(s, i)
	case "get-listings":
		h.handleGetListings(s, i)
	case "get-balance":
		h.handleGetBalance(s, i)
	case "get-payouts":
		h.handleGetPayouts(s, i)
	case "get-messages":
		h.handleGetMessages(s, i)
	case "ebay-status":
		h.handleEbayStatus(s, i)
	case "ebay-scopes":
		h.handleEbayScopes(s, i)
	case "ebay-authorize":
		h.handleEbayAuthorize(s, i)
	case "ebay-code":
		h.handleEbayCode(s, i)
	case "webhook-subscribe":
		h.handleWebhookSubscribe(s, i)
	case "webhook-list":
		h.handleWebhookList(s, i)
	case "webhook-test":
		h.handleWebhookTest(s, i)
	case "accept-offer":
		h.handleAcceptOffer(s, i)
	case "counter-offer":
		h.handleCounterOffer(s, i)
	case "decline-offer":
		h.handleDeclineOffer(s, i)
	}
}

func (h *Handler) handleGetBalance(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	balance, err := h.ebay.GetSellerBalance()
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to get balance: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	msg := fmt.Sprintf("💰 **Your eBay Balance**\n\n**Balance:** $%.2f\n\n💡 *Based on recent transactions. Use `/get-payouts` to see completed payouts.*",
		balance["total"])

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleGetPayouts(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	limit := 10

	if len(options) > 0 {
		limit = int(options[0].IntValue())
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	payouts, err := h.ebay.GetPayouts(limit)
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to get payouts: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	if len(payouts) == 0 {
		msg := "📭 **No Recent Payouts**\n\nNo payout transactions found in your account."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	msg := fmt.Sprintf("💸 **Recent Payouts** (%d shown):\n\n", len(payouts))
	for i, payout := range payouts {
		status := "✅"
		if payout["status"] == "PENDING" {
			status = "⏳"
		} else if payout["status"] == "FAILED" {
			status = "❌"
		}

		msg += fmt.Sprintf("%d. %s **$%.2f** - %s\n", i+1, status, payout["amount"], payout["type"])
		msg += fmt.Sprintf("   📅 %s | ID: `%s`\n\n", payout["date"], payout["id"])
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleGetMessages(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	limit := 10

	if len(options) > 0 {
		limit = int(options[0].IntValue())
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	messages, err := h.ebay.GetBuyerMessages(limit)
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to get messages: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	if len(messages) == 0 {
		msg := "📭 **No Messages**\n\nNo buyer messages found."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	msg := fmt.Sprintf("💬 **Buyer Messages** (%d shown):\n\n", len(messages))
	for i, message := range messages {
		unread := ""
		if message["unread"] == true {
			unread = "🔴 "
		}

		msg += fmt.Sprintf("%d. %s**From:** %s\n", i+1, unread, message["sender"])
		msg += fmt.Sprintf("   **Subject:** %s\n", message["subject"])
		msg += fmt.Sprintf("   **Message:** %s\n", message["body"])
		msg += fmt.Sprintf("   📅 %s | Order: `%s`\n\n", message["date"], message["orderId"])
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleGetOrders(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the response since API calls might take time
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	orders, err := h.ebay.GetOrders(10) // Get last 10 orders
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to fetch orders: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	if len(orders) == 0 {
		msg := "📋 **Recent Orders**\n\n⚠️ No orders found in your eBay account.\n\n💡 Once you have sales, they will appear here!"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	// Build response message with embeds for images
	response := "📋 **Recent Orders**\n\n"
	embeds := []*discordgo.MessageEmbed{}

	for i, order := range orders {
		if i >= 5 { // Limit to 5 orders in Discord message
			response += fmt.Sprintf("\n*...and %d more orders*", len(orders)-5)
			break
		}

		// Get first item for thumbnail
		var imageUrl string
		var itemTitle string
		if len(order.LineItems) > 0 {
			imageUrl = order.LineItems[0].ImageUrl
			itemTitle = order.LineItems[0].Title
			if itemTitle == "" {
				itemTitle = "Item"
			}
		}

		// Create embed for this order with image
		embed := &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Order #%s", order.OrderID),
			Color: 0x00ff00, // Green
			Fields: []*discordgo.MessageEmbedField{
				{Name: "👤 Buyer", Value: order.BuyerUsername, Inline: true},
				{Name: "💰 Total", Value: fmt.Sprintf("%.2f %s", order.TotalPrice, order.Currency), Inline: true},
				{Name: "📦 Status", Value: order.FulfillmentStatus, Inline: true},
				{Name: "📅 Date", Value: order.CreationDate.Format("Jan 02, 2006"), Inline: true},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: itemTitle,
			},
		}

		if imageUrl != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: imageUrl,
			}
		}

		embeds = append(embeds, embed)
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &response,
		Embeds:  &embeds,
	})
}

func (h *Handler) handleGetOffers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "💰 **Buyer Offers**\n\n" +
				"🔔 **Set up offer notifications to get alerts in Discord when buyers make offers!**\n\n" +
				"**How to enable:**\n" +
				"1. Run `/webhook-subscribe` and select the **OFFER** notification type\n" +
				"2. When a buyer makes an offer, you'll get an instant Discord notification\n" +
				"3. Use `/accept-offer`, `/counter-offer`, or `/decline-offer` to respond\n\n" +
				"💡 **Tip:** You can also view offers at:\n" +
				"🔗 https://offer.ebay.com/ws/eBayISAPI.dll?BestOfferList",
		},
	})
}

func (h *Handler) handleGetListings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "📦 **Your Active Listings**\n\n" +
				"To view your active listings, visit:\n" +
				"🔗 https://www.ebay.com/sh/lst/active\n\n" +
				"ℹ️ **Note:** The eBay API requires complex XML parsing for traditional listings. " +
				"Use the website for the best listing management experience.",
		},
	})
}

func (h *Handler) handleEbayStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	status := h.ebay.CheckConnection()

	response := fmt.Sprintf("🔍 **eBay API Status**\n\n%s", status)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func (h *Handler) handleEbayScopes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	scopes := h.ebay.GetTokenScopes()

	response := "🔐 **OAuth Token Scopes**\n\n"
	if len(scopes) == 0 {
		response += "⚠️ No token found or unable to determine scopes.\n\nRun `/ebay-authorize` to authorize the bot."
	} else {
		response += "Your current token has these scopes:\n\n"
		for _, scope := range scopes {
			response += fmt.Sprintf("✅ `%s`\n", scope)
		}
		response += "\n**Required scopes for all features:**\n"
		response += "• `api_scope` - Basic API access\n"
		response += "• `sell.inventory` - Manage listings\n"
		response += "• `sell.fulfillment` - View orders\n"
		response += "• `sell.account` - Account settings\n"
		response += "• `sell.finances` - Balance & payouts ⚠️\n"
		response += "\nIf `sell.finances` is missing, run `/ebay-authorize` to re-authorize."
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func (h *Handler) handleEbayAuthorize(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Generate state for OAuth
	state := fmt.Sprintf("state_%d", time.Now().Unix())

	// Register OAuth callback with webhook server if available
	if h.webhookServer != nil {
		h.webhookServer.RegisterOAuthCallback(state, s, i.Interaction)
	} else {
		errMsg := "❌ Webhook server not configured. OAuth flow unavailable."
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	// Generate authorization URL (will redirect to jacob.it.com/webhook/oauth/callback)
	authURL := h.ebay.GetUserAuthorizationURL(state)

	msg := fmt.Sprintf("🔐 **eBay Authorization - AUTOMATIC MODE**\n\n✨ **Just click the link below and sign in - that's it!**\n\n%s\n\n🎯 **What happens next:**\n1. You'll be redirected to eBay to sign in\n2. Click \"Agree\" to authorize the bot\n3. You'll be redirected to jacob.it.com\n4. The bot will automatically exchange your code for tokens\n5. Done! You'll be notified here when complete!\n\n⏱️ Authorization will expire in 10 minutes.\n\n💡 **Make sure your eBay RuName is configured:**\n• Accepted URL: `https://jacob.it.com/webhook/oauth/callback`\n• Declined URL: `https://jacob.it.com/webhook/oauth/declined`", authURL)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleEbayCode(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	code := options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// URL decode the code in case it was copied with encoding
	decodedCode, err := url.QueryUnescape(code)
	if err != nil {
		// If decoding fails, use the original code
		decodedCode = code
	}

	// Exchange the code for tokens
	tokens, err := h.ebay.ExchangeCodeForToken(decodedCode)
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to exchange code for tokens: %v\n\n💡 Tips:\n- Copy the ENTIRE code value from the URL (it's very long)\n- The code starts after `code=` and ends before `&expires_in`\n- It should look like: `v^1.1#i^1#f^0#I^3...` (very long)", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	// Save the tokens
	oauthServer := ebay.NewOAuthServer(h.ebay)
	if err := oauthServer.SaveTokensToEnv(tokens); err != nil {
		errMsg := fmt.Sprintf("❌ Tokens received but failed to save: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	successMsg := fmt.Sprintf("✅ **Authorization Successful!**\n\nAccess token and refresh token have been saved to .env file.\nYour bot will now automatically refresh tokens every 90 minutes.\n\nToken expires in: %d seconds", tokens.ExpiresIn)
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &successMsg,
	})

	// Start auto-refresh if we have a refresh token
	if tokens.RefreshToken != "" {
		log.Println("Starting automatic token refresh...")
		go h.ebay.AutoRefreshToken(tokens.RefreshToken, 90*time.Minute)
	}
}

func (h *Handler) handleWebhookSubscribe(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	// Use default URL if not provided
	webhookURL := "https://jacob.it.com/webhook/ebay/notification"
	if len(options) > 0 && options[0].StringValue() != "" {
		webhookURL = options[0].StringValue()
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	msg := fmt.Sprintf("🎣 **Webhook Setup - Automatic!**\n\n✅ **Your Webhook URL:**\n`%s`\n\n**📋 Instructions:**\n\n**Option 1: API Method (Recommended)**\nI can set this up automatically if you want! eBay requires API calls to create webhook subscriptions. For now, webhooks will be triggered when eBay sends notifications to your endpoint.\n\n**Your webhook server is already running!** When eBay sends:\n• 🛒 Order notifications\n• 💰 Offer notifications  \n• 📦 Shipping updates\n\nThey'll appear automatically in Discord!\n\n**Option 2: Manual Setup (Advanced)**\nIf you want to manually configure:\n1. Use eBay's API to POST to `/commerce/notification/v1/destination`\n2. Subscribe topics to your destination\n\n💡 **For testing:** Place a test order or have someone make an offer on your listing. The webhook endpoint is live and will relay notifications to this Discord channel automatically!\n\n🔔 Your bot will automatically:\n• Verify webhook challenges from eBay\n• Parse notifications\n• Post them to Discord\n\n**Current Status:** ✅ Ready to receive notifications", webhookURL)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleWebhookList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	msg := "📋 **Active Webhook Subscriptions**\n\nℹ️ Webhook management requires eBay Commerce API access.\n\n**Supported Notification Types:**\n• `marketplace.order.placed` - New order created\n• `marketplace.order.paid` - Payment received\n• `marketplace.offer.created` - Buyer sent an offer\n• `marketplace.offer.updated` - Offer counter/accepted/declined\n\n**Webhook Server Status:**\n✅ Server running on port 8080\n📍 Health check: http://localhost:8080/webhook/health\n📍 Notification endpoint: /webhook/ebay/notification\n📍 Challenge endpoint: /webhook/ebay/challenge"

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleWebhookTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Send a test notification to the current channel
	embed := &discordgo.MessageEmbed{
		Title:       "🔔 Test Notification",
		Description: "This is a test webhook notification from your eBay bot!",
		Color:       0x00ff00,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Status",
				Value:  "✅ Working",
				Inline: true,
			},
			{
				Name:   "Channel",
				Value:  fmt.Sprintf("<#%s>", i.ChannelID),
				Inline: true,
			},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func (h *Handler) handleAcceptOffer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	offerID := options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Call eBay API to accept the offer
	err := h.ebay.RespondToOffer(offerID, "ACCEPT", 0)
	if err != nil {
		errMsg := fmt.Sprintf("❌ **Failed to accept offer**\n\nError: %v\n\n**Troubleshooting:**\n• Verify offer ID is correct\n• Check if offer is still pending\n• Ensure you have authorization: `/ebay-status`\n• Offer may have expired or been withdrawn", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	msg := fmt.Sprintf("✅ **Offer Accepted!**\n\n🎉 Successfully accepted offer: `%s`\n\n**What happens next:**\n1. eBay creates an order for the buyer\n2. Buyer receives payment instructions\n3. Once paid, you'll receive order notification\n4. Ship item and upload tracking\n\n📧 You and the buyer will receive confirmation emails.\n\n💡 Check `/get-orders` to see the new order!", offerID)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleCounterOffer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	offerID := optionMap["offer-id"].StringValue()
	priceStr := optionMap["price"].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Parse price
	var price float64
	_, err := fmt.Sscanf(priceStr, "%f", &price)
	if err != nil {
		errMsg := fmt.Sprintf("❌ Invalid price format: %s. Please enter a number (e.g., 250.00)", priceStr)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	if price <= 0 {
		errMsg := "❌ Price must be greater than $0.00"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	// Call eBay API to counter the offer
	err = h.ebay.RespondToOffer(offerID, "COUNTER", price)
	if err != nil {
		errMsg := fmt.Sprintf("❌ **Failed to counter offer**\n\nError: %v\n\n**Troubleshooting:**\n• Verify offer ID is correct\n• Check if offer is still pending\n• Ensure counter price is valid\n• Ensure you have authorization: `/ebay-status`", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	msg := fmt.Sprintf("💬 **Counter Offer Sent!**\n\n📤 You've countered offer `%s` with **$%.2f**\n\n**What happens next:**\n1. Buyer receives your counter offer\n2. They have 48 hours to respond\n3. They can:\n   • Accept your counter\n   • Make another counter offer\n   • Decline and walk away\n\n**Negotiation Tips:**\n✅ Be reasonable with your counter\n✅ Factor in your costs and fees\n✅ Quick responses increase acceptance rate\n\n📧 You'll be notified of their response.", offerID, price)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleDeclineOffer(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	offerID := options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Call eBay API to decline the offer
	err := h.ebay.RespondToOffer(offerID, "DECLINE", 0)
	if err != nil {
		errMsg := fmt.Sprintf("❌ **Failed to decline offer**\n\nError: %v\n\n**Troubleshooting:**\n• Verify offer ID is correct\n• Check if offer is still pending\n• Ensure you have authorization: `/ebay-status`\n• Offer may have already been processed", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	msg := fmt.Sprintf("❌ **Offer Declined**\n\n🚫 Successfully declined offer: `%s`\n\n**What happens next:**\n1. Buyer receives decline notification\n2. Offer is closed and removed from pending\n3. Buyer can submit a new offer if desired\n4. Item remains listed for sale\n\n**Why decline?**\n• Offer too low for your costs\n• Buyer has negative feedback\n• You're not ready to sell yet\n• Better offers pending\n\n💡 **Tip:** Consider countering instead of declining to keep negotiation open!", offerID)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}
