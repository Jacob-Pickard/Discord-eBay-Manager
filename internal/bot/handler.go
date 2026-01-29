package bot

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"ebaymanager-bot/internal/ebay"
)

// Handler manages Discord bot interactions
type Handler struct {
	discord *discordgo.Session
	ebay    *ebay.Client
}

// NewHandler creates a new bot handler
func NewHandler(discord *discordgo.Session, ebayClient *ebay.Client) *Handler {
	return &Handler{
		discord: discord,
		ebay:    ebayClient,
	}
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
					Description: "Your public webhook URL (e.g., https://yourdomain.com/webhook/ebay/notification)",
					Required:    true,
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
	case "get-balance":
		h.handleGetBalance(s, i)
	case "get-payouts":
		h.handleGetPayouts(s, i)
	case "get-messages":
		h.handleGetMessages(s, i)
	case "ebay-status":
		h.handleEbayStatus(s, i)
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

	msg := fmt.Sprintf("💰 **Your eBay Balance**\n\n**Total Balance:** $%.2f\n**Available:** $%.2f\n**Pending:** $%.2f\n\n📊 **Breakdown:**\n• Sales this month: $%.2f\n• Fees this month: $%.2f\n• Net income: $%.2f\n\n💡 Use `/get-payouts` to see recent transactions",
		balance["totalBalance"], balance["available"], balance["pending"],
		balance["salesThisMonth"], balance["feesThisMonth"], balance["netIncome"])

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
		msg := "📋 **Recent Orders**\n\n⚠️ No orders found in your eBay account.\n\n💡 **Sample Order Data:**\n\n**Order #19-12345-67890**\n• Buyer: tech_enthusiast_2026\n• Item: Vintage Gaming Console Bundle\n• Total: $299.99 USD\n• Status: AWAITING_SHIPMENT\n• Date: Jan 28, 2026\n\n**Order #19-54321-09876**\n• Buyer: collector_pro\n• Item: Limited Edition Trading Cards\n• Total: $149.50 USD\n• Status: PAID\n• Date: Jan 27, 2026\n\n**Order #19-11111-22222**\n• Buyer: gadget_lover\n• Item: Wireless Headphones - Premium\n• Total: $89.99 USD\n• Status: IN_TRANSIT\n• Date: Jan 26, 2026\n\n✨ *This is sample data. Real orders will appear here once you have sales!*"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	// Build response message
	response := "📋 **Recent Orders**\n\n"
	for i, order := range orders {
		if i >= 5 { // Limit to 5 orders in Discord message
			response += fmt.Sprintf("\n*...and %d more orders*", len(orders)-5)
			break
		}
		response += fmt.Sprintf("**Order #%s**\n", order.OrderID)
		response += fmt.Sprintf("• Buyer: %s\n", order.BuyerUsername)
		response += fmt.Sprintf("• Total: %.2f %s\n", order.TotalPrice, order.Currency)
		response += fmt.Sprintf("• Status: %s\n", order.FulfillmentStatus)
		response += fmt.Sprintf("• Date: %s\n\n", order.CreationDate.Format("Jan 02, 2006"))
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &response,
	})
}

func (h *Handler) handleGetOffers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	offers, err := h.ebay.GetOffers()
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to fetch offers: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	if len(offers) == 0 {
		msg := "💰 **Pending Buyer Offers**\n\n⚠️ No pending offers found.\n\n💡 **Sample Offer Data:**\n\n**Offer from vintage_collector**\n• Item: Vintage Camera Collection\n• Listed Price: $500.00 USD\n• Offer Amount: $425.00 USD (15% off)\n• Message: \"Very interested! Can you accept $425?\"\n• Status: PENDING\n• Expires: 48 hours\n\n**Offer from sneaker_enthusiast**\n• Item: Limited Edition Sneakers - Size 10\n• Listed Price: $250.00 USD\n• Offer Amount: $225.00 USD (10% off)\n• Message: \"Great condition! Would love to buy at $225\"\n• Status: PENDING\n• Expires: 36 hours\n\n**Offer from gaming_pro**\n• Item: Gaming PC - RTX 4090\n• Listed Price: $2,499.00 USD\n• Offer Amount: $2,200.00 USD (12% off)\n• Message: \"Cash ready! Can pick up today\"\n• Status: COUNTERED (Your counter: $2,350)\n• Expires: 24 hours\n\n✨ *This is sample data. Real offers will appear here once buyers make offers!*\n\n💡 **Actions you can take:**\n• Accept offer - Automatically creates sale\n• Counter offer - Negotiate price\n• Decline - Reject the offer\n\n*Note: Full offer management will be available with eBay Trading API integration.*"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	response := "💰 **Pending Offers**\n\n"
	for i, offer := range offers {
		if i >= 5 {
			response += fmt.Sprintf("\n*...and %d more offers*", len(offers)-5)
			break
		}
		response += fmt.Sprintf("**Offer #%s**\n", offer.OfferID)
		response += fmt.Sprintf("• Item: %s\n", offer.ItemTitle)
		response += fmt.Sprintf("• Buyer: %s\n", offer.BuyerUsername)
		response += fmt.Sprintf("• Offer: %.2f %s (List: %.2f)\n", offer.OfferPrice, offer.Currency, offer.ListPrice)
		response += fmt.Sprintf("• Status: %s\n\n", offer.Status)
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &response,
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

func (h *Handler) handleEbayAuthorize(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Generate authorization URL
	authURL := h.ebay.GetUserAuthorizationURL("")
	
	msg := fmt.Sprintf("🔐 **eBay Authorization**\n\n**Step 1:** Open this URL in your browser:\n%s\n\n**Step 2:** Sign in with your eBay account and click \"Agree\"\n\n**Step 3:** After you authorize, eBay will redirect you to a page. Look at the URL in your browser - it will contain `code=...`\n\n**Step 4:** Copy the code value from the URL (everything after `code=` and before the next `&` if there is one)\n\n**Step 5:** Use the command `/ebay-code` and paste the code\n\n💡 Example URL after redirect:\n`https://signin.ebay.com/ws/eBayISAPI.dll?...&code=v%5E1.1%23i%5E1...&expires_in=299`\n\nCopy only the code part!", authURL)
	
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
	webhookURL := options[0].StringValue()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Topics that will be available for subscription
	_ = []string{
		"marketplace.order.placed",
		"marketplace.order.paid",
		"marketplace.offer.created",
		"marketplace.offer.updated",
	}

	msg := fmt.Sprintf("🎣 **Webhook Subscription**\n\n📍 **Endpoint**: %s\n\n⚠️ **Important Setup Steps:**\n\n1. Make sure your webhook URL is **publicly accessible** (use ngrok, a VPS, or cloud hosting)\n2. The bot's webhook server must be running on port 8080\n3. eBay will send a challenge request to verify your endpoint\n4. Once verified, you'll receive notifications for:\n   • New orders\n   • Payment received\n   • Offers from buyers\n   • Offer updates\n\n💡 **Testing locally?** Use ngrok:\n```\nngrok http 8080\n```\nThen use the ngrok URL (e.g., https://abc123.ngrok.io/webhook/ebay/notification)\n\n✅ Run `/webhook-list` to see your active subscriptions.", webhookURL)
	
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
