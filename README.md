# eBay Manager Discord Bot

A Discord bot written in Go that helps manage your eBay business, including listing items, tracking orders, handling offers, and receiving real-time notifications through webhooks.

## ⚡ Features

### ✅ Currently Implemented

- **🔐 OAuth 2.0 Authentication**
  - Full OAuth flow with authorization code exchange
  - Automatic token refresh (every 90 minutes)
  - Tokens saved to `.env` file
  - Manual code submission via Discord command

- **📦 Order Management**
  - View recent orders with `/get-orders`
  - Detailed order information (buyer, items, status, total)
  - Fallback to sample data in sandbox environment

- **💰 Offer Management**
  - View pending buyer offers with `/get-offers`
  - Accept offers with `/accept-offer`
  - Counter offers with `/counter-offer`
  - Decline offers with `/decline-offer`

- **💵 Financial Information**
  - View account balance with `/get-balance`
  - View recent payouts with `/get-payouts`
  - Transaction history

- **💬 Buyer Messages**
  - View buyer messages with `/get-messages`
  - Unread message indicators

- **🔔 Real-time Webhook Notifications**
  - Webhook server on port 8080
  - SHA-256 challenge verification
  - Automatic Discord notifications for:
    - New orders (`MARKETPLACE_ACCOUNT.ORDER.FULFILLED`)
    - Best offers (`MARKETPLACE_ACCOUNT.OFFER.UPDATED`)
    - Messages (`MARKETPLACE_ACCOUNT.QUESTION.CREATED`)
  - Subscribe/list webhooks via Discord commands

- **📊 Status Monitoring**
  - Check eBay API connection status
  - View environment (Sandbox/Production)
  - Token validation

### 🚧 Planned / In Development

- **📝 Listing Creation**
  - Create new eBay listings directly from Discord
  - Image upload support
  - Bulk listing management

- **🏷️ Shipping Integration**
  - Purchase shipping labels
  - Track shipments
  - Update tracking information on eBay

- **📈 Analytics & Reporting**
  - Sales statistics
  - Performance metrics
  - Revenue tracking

- **🤖 Automation**
  - Automated offer responses based on rules
  - Bulk operations
  - Scheduled tasks

## 📋 Prerequisites

- Go 1.21 or higher
- Discord Bot Token
- eBay Developer Account with API credentials
- Public URL for webhook endpoint (ngrok, hosted server, etc.)

## Setup

### 1. Get Discord Bot Token

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a new application
3. Go to the "Bot" section and create a bot
4. Copy the bot token
5. Under "OAuth2" → "URL Generator":
   - Select scopes: `bot`, `applications.commands`
   - Select permissions: `Send Messages`, `Use Slash Commands`
   - Copy the generated URL and invite the bot to your server

### 2. Get eBay API Credentials

1. Sign up at [eBay Developers Program](https://developer.ebay.com/)
2. Create a new application (use Sandbox for testing)
3. Get your credentials:
   - App ID (Client ID)
   - Cert ID (Client Secret)
   - Dev ID
4. Set up RuName (OAuth redirect URI):
   - Create a RuName in your eBay Developer account
   - Use `https://localhost:3000/callback` or your own URL, or leave the accept and deny fields empty.
5. **Note:** You'll get OAuth tokens through the bot using `/ebay-authorize` command

### 3. Configure Environment

1. Copy `.env.example` to `.env`:
   ```bash
   cp .env.example .env
   ```

2. Fill in your credentials in `.env`:
   ```env
   # Discord Configuration
   DISCORD_BOT_TOKEN=your_discord_bot_token
   
   # eBay API Configuration
   EBAY_APP_ID=your_ebay_app_id
   EBAY_CERT_ID=your_ebay_cert_id
   EBAY_DEV_ID=your_ebay_dev_id
   EBAY_REDIRECT_URI=your_redirect_uri
   
   # eBay OAuth Tokens (generated via /ebay-authorize command)
   EBAY_ACCESS_TOKEN=
   EBAY_REFRESH_TOKEN=
   
   # Environment (PRODUCTION or SANDBOX)
   EBAY_ENVIRONMENT=SANDBOX
   
   # Webhook Configuration (optional)
   WEBHOOK_PORT=8080
   WEBHOOK_VERIFY_TOKEN=your_secure_random_token
   NOTIFICATION_CHANNEL_ID=your_discord_channel_id
   ```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Bot

```bash
go run main.go
```

Or build and run:
```bash
go build -o ebaymanager-bot
./ebaymanager-bot  # On Linux/Mac
# or
ebaymanager-bot.exe  # On Windows
```

## 🎮 Discord Commands

### Authentication
- `/ebay-status` - Check eBay API connection status
- `/ebay-authorize` - Get authorization URL to connect your eBay account
- `/ebay-code` - Submit authorization code after eBay redirect

### Order Management
- `/get-orders` - View recent eBay orders with buyer details

### Offer Management
- `/get-offers` - View pending buyer offers (best offers)
- `/accept-offer` - Accept a buyer's offer
- `/counter-offer` - Counter an offer with a different price
- `/decline-offer` - Decline a buyer's offer

### Financial Information
- `/get-balance` - View your eBay account balance
- `/get-payouts` - View recent payouts and transactions

### Communication
- `/get-messages` - View buyer messages

### Webhook Management
- `/webhook-subscribe` - Subscribe to eBay real-time notifications
- `/webhook-list` - List active webhook subscriptions
- `/webhook-test` - Test webhook notification to current channel

## 📁 Project Structure

```
ebaymanager-bot/
├── main.go                 # Application entry point
├── internal/
│   ├── bot/
│   │   └── handler.go     # Discord bot command handlers
│   ├── config/
│   │   └── config.go      # Configuration management
│   ├── ebay/
│   │   ├── client.go      # eBay API client
│   │   ├── oauth.go       # OAuth 2.0 implementation
│   │   ├── oauth_flow.go  # OAuth flow handlers
│   │   └── types.go       # eBay data structures
│   └── webhook/
│       ├── server.go      # Webhook HTTP server
│       └── subscription.go # Webhook subscription management
├── go.mod                 # Go module definition
├── .env                   # Environment variables (DO NOT COMMIT)
├── .env.example          # Example environment variables
└── .gitignore            # Git ignore rules (includes .env)
```

## 🚀 Development Status

### ✅ Phase 1: Core Infrastructure (COMPLETE)
- [x] Discord bot setup with slash commands
- [x] eBay API client structure
- [x] Configuration management via environment variables

### ✅ Phase 2: eBay OAuth & Authentication (COMPLETE)
- [x] OAuth 2.0 authorization code flow
- [x] Automatic token refresh (every 90 minutes)
- [x] Secure token storage in .env file
- [x] Discord command integration for auth

### ✅ Phase 3: Basic eBay Operations (COMPLETE)
- [x] Fetch orders from Fulfillment API
- [x] Fetch offers from Negotiation API
- [x] Accept/decline/counter offers via Negotiation API
- [x] View account balance and payouts
- [x] View buyer messages

### ✅ Phase 4: Webhook Integration (COMPLETE)
- [x] Webhook HTTP server with challenge verification
- [x] Real-time Discord notifications for orders, offers, messages
- [x] Subscription management via Discord commands
- [x] Production-ready SHA-256 verification

### 🚧 Phase 5: Listing Management (PLANNED)
- [ ] Create listings using Inventory API
- [ ] Bulk listing operations
- [ ] Image upload support
- [ ] Inventory tracking

### 🚧 Phase 6: Shipping Integration (PLANNED)
- [ ] Purchase shipping labels via Buy API
- [ ] Track shipments
- [ ] Update tracking information on orders

### 🚧 Phase 7: Advanced Features (PLANNED)
- [ ] Analytics and sales reporting
- [ ] Automated offer response rules
- [ ] Inventory tracking and alerts
- [ ] Performance metrics dashboard

## 📊 eBay API Endpoints Currently Used

### ✅ Implemented
- **OAuth 2.0 API**: `/identity/v1/oauth2/token` - Token exchange & refresh
- **Fulfillment API**: `/sell/fulfillment/v1/order` - Fetch orders
- **Negotiation API**: `/sell/negotiation/v1/offer` - Manage buyer offers
- **Finances API**: `/sell/finances/v1/seller_funds_summary` - Account balance
- **Finances API**: `/sell/finances/v1/payout` - View payouts
- **Post-Order API**: `/post-order/v2/inquiry/search` - Buyer messages
- **Marketplace Notification API**: `/commerce/notification/v1/destination` - Webhooks

### 🚧 Planned
- **Inventory API**: `/sell/inventory/v1/inventory_item` - Create/manage listings
- **Buy API**: For purchasing shipping labels
- **Analytics API**: Sales and performance data

## 🔒 Security & Production Readiness

### ✅ Security Measures in Place
- `.env` file excluded from Git via `.gitignore`
- OAuth tokens never hardcoded in source
- Webhook signature verification (SHA-256)
- Automatic token refresh to prevent expiration
- Environment-based configuration (Sandbox/Production)

### ⚠️ Before Going to Production

1. **Switch to Production Environment**
   - Change `EBAY_ENVIRONMENT=PRODUCTION` in `.env`
   - Re-authenticate with production eBay account
   - Update eBay app credentials to production keys

2. **Webhook Setup Requirements**
   - Deploy to a server with a public IP/domain
   - Use HTTPS (required by eBay)
   - Update `WEBHOOK_VERIFY_TOKEN` to a strong random value
   - Test all webhook subscriptions in production

3. **Testing Requirements**
   - Test offer management with real sandbox offers
   - Verify webhook notifications are received
   - Test token refresh after 90 minutes
   - Validate all Discord commands work correctly

4. **Recommended Enhancements**
   - Implement database for order/offer history
   - Add logging to file or monitoring service
   - Set up error alerting (email/Discord)
   - Add rate limiting for API calls
   - Implement backup/restore for tokens

5. **Deployment Checklist**
   - [ ] Server has public HTTPS endpoint
   - [ ] All environment variables configured
   - [ ] Webhook endpoints accessible from internet
   - [ ] Discord bot invited to production server
   - [ ] eBay app has production credentials
   - [ ] Monitoring and logging configured
   - [ ] Backup strategy for .env file

See [PRODUCTION_READINESS.md](PRODUCTION_READINESS.md) for detailed testing and deployment guide.

## 📚 Additional Documentation

- [API_GUIDE.md](API_GUIDE.md) - Comprehensive eBay API reference
- [QUICKSTART.md](QUICKSTART.md) - Quick setup guide
- [WEBHOOK_SETUP.md](WEBHOOK_SETUP.md) - Webhook configuration details
- [PRODUCTION_READINESS.md](PRODUCTION_READINESS.md) - Production deployment checklist
- [TESTING_GUIDE.md](TESTING_GUIDE.md) - Testing procedures

## 🤝 Contributing

This is a personal project, but feel free to fork and modify for your own use. Pull requests are welcome!

## ⚠️ Important Security Notes

**NEVER COMMIT SENSITIVE DATA:**
- ✅ `.env` is in `.gitignore` (verified)
- ✅ No tokens or API keys in source code
- ⚠️ Always use `.env.example` as template, never copy actual values to commits
- ⚠️ Rotate tokens if accidentally committed
- ⚠️ Keep `WEBHOOK_VERIFY_TOKEN` secret and strong

**File Security:**
- `.env` - Contains all secrets (Discord token, eBay credentials, OAuth tokens)
- `.gitignore` - Configured to block .env, *.token, *.key, credentials.json, etc.
- `.env.example` - Safe template with placeholder values only

## 📄 License

MIT License - See LICENSE file for details

## 🔗 Resources

- [eBay API Documentation](https://developer.ebay.com/docs)
- [eBay Developer Program](https://developer.ebay.com/)
- [discordgo Documentation](https://github.com/bwmarrin/discordgo)
- [eBay API Explorer](https://developer.ebay.com/my/api_test_tool)
- [OAuth 2.0 Guide](https://developer.ebay.com/api-docs/static/oauth-tokens.html)

---

**Status:** 🟢 Sandbox-Ready | 🟡 Production requires testing & deployment

**Last Updated:** January 2026
