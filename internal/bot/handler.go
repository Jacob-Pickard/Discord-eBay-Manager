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

	msg := fmt.Sprintf("💰 **Your eBay Balance**\n\n**Available for Next Payout:** $%.2f\n**Total Balance:** $%.2f\n\n💡 *Available funds will be included in your next scheduled payout. Use `/get-payouts` to see completed payouts.*",
		balance["available"], balance["total"])

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
		if statusStr, ok := payout["status"].(string); ok {
			if statusStr == "PENDING" {
				status = "⏳"
			} else if statusStr == "FAILED" {
				status = "❌"
			}
		}

		msg += fmt.Sprintf("%d. %s **$%.2f** - %s\n", i+1, status, payout["amount"], payout["type"])
		msg += fmt.Sprintf("   📅 %s | ID: `%s`\n\n", payout["date"], payout["id"])
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
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Try to fetch offers from eBay API
	offers, err := h.ebay.GetOffers()
	if err != nil {
		// If API call fails, show setup instructions
		msg := "💰 **Buyer Offers**\n\n" +
			"⚠️ Unable to fetch offers directly from API.\n\n" +
			"🔔 **Set up offer notifications to get alerts in Discord when buyers make offers!**\n\n" +
			"**How to enable:**\n" +
			"1. Run `/webhook-subscribe` to set up webhook notifications\n" +
			"2. When a buyer makes an offer, you'll get an instant Discord notification\n" +
			"3. Use `/accept-offer`, `/counter-offer`, or `/decline-offer` to respond\n\n" +
			"💡 **Tip:** You can also view offers at:\n" +
			"🔗 https://offer.ebay.com/ws/eBayISAPI.dll?BestOfferList"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	if len(offers) == 0 {
		msg := "💰 **Buyer Offers**\n\n" +
			"📭 No pending offers found.\n\n" +
			"💡 **Enable real-time notifications:**\n" +
			"Run `/webhook-subscribe` to get instant alerts when buyers make offers!\n\n" +
			"🔗 You can also check: https://offer.ebay.com/ws/eBayISAPI.dll?BestOfferList"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	// Display offers
	msg := fmt.Sprintf("💰 **Pending Offers** (%d found):\n\n", len(offers))
	for i, offer := range offers {
		if i >= 10 { // Limit to 10 offers
			msg += fmt.Sprintf("\n*...and %d more offers*", len(offers)-10)
			break
		}
		status := "⏳"
		if offer.Status == "PENDING" {
			status = "⏳"
		} else if offer.Status == "ACCEPTED" {
			status = "✅"
		} else if offer.Status == "DECLINED" {
			status = "❌"
		}

		msg += fmt.Sprintf("%d. %s **$%.2f** (List: $%.2f)\n", i+1, status, offer.OfferPrice, offer.ListPrice)
		msg += fmt.Sprintf("   👤 %s | 📦 %s\n", offer.BuyerUsername, offer.ItemTitle)
		msg += fmt.Sprintf("   🆔 `%s`\n", offer.OfferID)
		msg += fmt.Sprintf("   📅 %s\n\n", offer.CreatedDate.Format("Jan 02, 2006"))
	}

	msg += "\n💡 Respond with:\n" +
		"• `/accept-offer offer-id:<ID>`\n" +
		"• `/counter-offer offer-id:<ID> price:<AMOUNT>`\n" +
		"• `/decline-offer offer-id:<ID>`"

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleGetListings(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	limit := 10
	if len(options) > 0 {
		limit = int(options[0].IntValue())
	}

	log.Printf("[listings] Command triggered, limit=%d", limit)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	listings, err := h.ebay.GetListings(limit)
	if err != nil {
		log.Printf("[listings] ERROR: %v", err)
		errMsg := fmt.Sprintf("❌ Failed to fetch listings: %v", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &errMsg})
		return
	}

	log.Printf("[listings] Got %d listings", len(listings))

	if len(listings) == 0 {
		msg := "📦 **Active Listings**\n\n" +
			"⚠️ No active listings found for your eBay account.\n\n" +
			"💡 View your listings directly at:\n" +
			"🔗 https://www.ebay.com/sh/lst/active"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{Content: &msg})
		return
	}

	header := fmt.Sprintf("📦 **Active Listings** (%d found)\n\u200b", len(listings))
	embeds := []*discordgo.MessageEmbed{}

	for idx, listing := range listings {
		if idx >= 5 {
			header += fmt.Sprintf("\n*...and %d more. Visit https://www.ebay.com/sh/lst/active for all*", len(listings)-5)
			break
		}

		priceStr := fmt.Sprintf("$%.2f %s", listing.Price, listing.Currency)
		if listing.Price == 0 {
			priceStr = "See listing"
		}

		condition := listing.Condition
		if condition == "" {
			condition = "Not specified"
		}

		title := listing.Title
		if title == "" {
			title = listing.SKU
		}

		fields := []*discordgo.MessageEmbedField{
			{Name: "💰 Price", Value: priceStr, Inline: true},
			{Name: "🚚 Shipping", Value: listing.Shipping, Inline: true},
			{Name: "📦 Qty", Value: fmt.Sprintf("%d", listing.Quantity), Inline: true},
			{Name: "🏷️ Condition", Value: condition, Inline: true},
		}
		if listing.SKU != "" {
			fields = append(fields, &discordgo.MessageEmbedField{Name: "🔑 SKU", Value: listing.SKU, Inline: true})
		}

		embed := &discordgo.MessageEmbed{
			Title:  title,
			Color:  0x0064d2, // eBay blue
			Fields: fields,
		}

		if listing.ListingURL != "" {
			embed.URL = listing.ListingURL
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: "🔗 Link", Value: listing.ListingURL, Inline: false,
			})
		}

		if listing.ImageURL != "" {
			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: listing.ImageURL}
		}

		embeds = append(embeds, embed)
	}

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &header,
		Embeds:  &embeds,
	})
}

func (h *Handler) handleEbayStatus(s *discordgo.Session, i *discordgo.InteractionCreate) {
	status := h.ebay.CheckConnection()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: status,
		},
	})
}

func (h *Handler) handleEbayScopes(s *discordgo.Session, i *discordgo.InteractionCreate) {
	scopeInfo := h.ebay.GetTokenScopes()

	hasToken := scopeInfo["hasToken"].(bool)
	env := scopeInfo["environment"].(string)
	requested := scopeInfo["requestedScopes"].([]string)

	response := "🔐 **OAuth Scopes Configuration**\n\n"

	if !hasToken {
		response += "❌ **No token found**\n\nRun `/ebay-authorize` to connect your eBay account.\n\n"
	} else {
		response += fmt.Sprintf("✅ **Token Status:** Active\n**Environment:** `%s`\n\n", env)
	}

	response += "**Requested Scopes:**\n"
	for _, scope := range requested {
		// Extract just the last part of the scope URL for readability
		parts := strings.Split(scope, "/")
		scopeName := parts[len(parts)-1]
		if strings.Contains(scopeName, ".") {
			response += fmt.Sprintf("• `%s`\n", scopeName)
		} else {
			response += fmt.Sprintf("• `%s` (Basic API access)\n", scopeName)
		}
	}

	response += "\n**What each scope does:**\n"
	response += "• **sell.inventory** - Manage your listings\n"
	response += "• **sell.fulfillment** - View & manage orders\n"
	response += "• **sell.account** - Account settings & policies\n"
	response += "• **sell.finances** - Balance, payouts & transactions\n"

	if !hasToken {
		response += "\n💡 Authorize now to enable all bot features!"
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

	log.Printf("🔔 webhook-subscribe command called for URL: %s", webhookURL)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Check for existing subscriptions first and delete any for the same endpoint
	log.Printf("📋 Checking for existing subscriptions...")
	existingSubscriptions, err := h.ebay.ListWebhookSubscriptions()
	if err != nil {
		log.Printf("⚠️ Warning: Failed to list existing subscriptions: %v", err)
		log.Printf("Proceeding anyway, but might encounter conflicts...")
	} else {
		log.Printf("✅ Found %d existing subscriptions", len(existingSubscriptions))
		foundMatch := false
		for _, sub := range existingSubscriptions {
			if endpoint, ok := sub.DeliveryConfig["endpoint"].(string); ok {
				log.Printf("  - Subscription %s: %s (status: %s)", sub.DestinationID, endpoint, sub.Status)
				if endpoint == webhookURL {
					foundMatch = true
					log.Printf("🗑️ Found existing subscription for %s, deleting it first (ID: %s)", webhookURL, sub.DestinationID)
					if delErr := h.ebay.DeleteWebhookSubscription(sub.DestinationID); delErr != nil {
						log.Printf("❌ Failed to delete existing subscription: %v", delErr)
						errMsg := fmt.Sprintf("❌ **Cannot create webhook subscription**\n\nThere's already a subscription for this endpoint, but I couldn't delete it.\n\nError: %v\n\n**Manual fix required:**\n1. Run `/webhook-list` to see the subscription ID\n2. Ask eBay support to delete it, or\n3. Try using a different webhook URL\n\n**Existing endpoint:** `%s`", delErr, endpoint)
						s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
							Content: &errMsg,
						})
						return
					} else {
						log.Printf("✅ Successfully deleted existing subscription for %s", webhookURL)
					}
				}
			}
		}
		if !foundMatch {
			log.Printf("ℹ️ No existing subscription found for %s", webhookURL)
		}
	}

	// Actually create the subscription with eBay
	log.Printf("🔨 Creating new webhook subscription for %s...", webhookURL)
	err = h.ebay.CreateWebhookSubscription(webhookURL)
	if err != nil {
		log.Printf("❌ Failed to create subscription: %v", err)
		errMsg := fmt.Sprintf("❌ **Failed to create webhook subscription**\n\nError: %v\n\n**Troubleshooting:**\n• Make sure you're authorized: `/ebay-authorize`\n• Check if subscription already exists: `/webhook-list`\n• Verify your webhook URL is accessible from the internet\n• URL must use HTTPS (not HTTP)\n• Make sure your webhook server is running and responding to challenges\n\n**Your webhook URL:** `%s`\n\n**Debug Info:**\nTo test if your webhook is reachable, visit:\n`%s?challenge_code=test`", err, webhookURL, webhookURL)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	log.Printf("✅ Successfully created webhook subscription for %s", webhookURL)

	msg := fmt.Sprintf("✅ **Webhook Subscription Created!**\n\n🎣 **Your Webhook URL:**\n`%s`\n\n**📋 Active Subscriptions:**\n• 🛒 **Order notifications** - New orders, payments, shipments\n• 💰 **Offer notifications** - New offers, counters, acceptances\n• 📦 **Inventory updates** - Listing changes\n• 🔔 **Account events** - Important account notifications\n\n**✨ What happens now:**\nWhen eBay sends notifications for these events, they'll appear automatically in this Discord channel!\n\n**🧪 Test it:**\n1. Have someone make an offer on one of your listings\n2. The notification will appear here within seconds!\n3. Use `/accept-offer`, `/counter-offer`, or `/decline-offer` to respond\n\n**📊 View subscriptions:** `/webhook-list`\n\n💡 Your webhook server is running at jacob.it.com and ready to receive notifications!", webhookURL)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &msg,
	})
}

func (h *Handler) handleWebhookList(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	// Fetch actual subscriptions from eBay
	subscriptions, err := h.ebay.ListWebhookSubscriptions()
	if err != nil {
		errMsg := fmt.Sprintf("❌ Failed to list subscriptions: %v\n\n💡 Make sure you're authorized with `/ebay-authorize`", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &errMsg,
		})
		return
	}

	if len(subscriptions) == 0 {
		msg := "📋 **Active Webhook Subscriptions**\n\n📭 No active subscriptions found.\n\n**Get started:**\nRun `/webhook-subscribe` to create a webhook subscription and start receiving real-time notifications for:\n• 🛒 Orders (placed, paid, shipped)\n• 💰 Offers (created, accepted, declined, countered)\n• 📦 Inventory updates\n• 🔔 Account events"
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &msg,
		})
		return
	}

	msg := fmt.Sprintf("📋 **Active Webhook Subscriptions** (%d found)\n\n", len(subscriptions))

	for i, sub := range subscriptions {
		status := "✅"
		if sub.Status != "ENABLED" {
			status = "⚠️"
		}

		msg += fmt.Sprintf("%d. %s **%s** (%s)\n", i+1, status, sub.Name, sub.Status)
		if endpoint, ok := sub.DeliveryConfig["endpoint"].(string); ok {
			msg += fmt.Sprintf("   📍 Endpoint: `%s`\n", endpoint)
		}
		msg += fmt.Sprintf("   🆔 ID: `%s`\n", sub.DestinationID)

		if len(sub.Topics) > 0 {
			msg += "   📢 Topics:\n"
			for _, topic := range sub.Topics {
				if topicName, ok := topic["topicName"]; ok {
					msg += fmt.Sprintf("      • %s\n", topicName)
				}
			}
		}
		msg += "\n"
	}

	msg += "\n**Webhook Server Status:**\n"
	msg += "✅ Server running on port 8081\n"
	msg += "📍 Health: https://jacob.it.com/webhook/health\n"
	msg += "📍 Notification endpoint: /webhook/ebay/notification\n"
	msg += "\n💡 To create a new subscription, run `/webhook-subscribe`"

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
