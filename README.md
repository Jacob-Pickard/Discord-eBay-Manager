# eBay Manager Discord Bot

A production-ready Discord bot written in Go that helps you manage your eBay business operations. View orders, respond to offers, and receive real-time notifications through webhooks - all from Discord.

**Current Status:** ✅ **Deployed to Production** | Running 24/7 on Ubuntu Server 22.04

## ⚡ Features

### ✅ Fully Implemented & Working

- **🔐 Automatic OAuth 2.0 Authentication**
  - **NEW:** Fully automatic OAuth flow through webhook server
  - User clicks authorization link → eBay redirects → Bot automatically exchanges code
  - No manual code copying required!
  - Automatic token refresh (every 90 minutes)
  - Tokens saved to `.env` file
  - Production webhook callbacks: `https://jacob.it.com/webhook/oauth/callback`

- **📦 Order Management**
  - View recent orders with `/get-orders`
  - **NEW:** Image loading with Browse API integration
  - Detailed order information (buyer, items, status, total)
  - Fallback image URL generation for legacy items
  - Rich Discord embeds with buyer information

- **🔔 Real-time Webhook Notifications**
  - Production webhook server on port 8081
  - Public endpoint: `https://jacob.it.com/webhook/`
  - SHA-256 challenge verification
  - Automatic Discord notifications for:
    - New orders (`MARKETPLACE_ACCOUNT.ORDER.FULFILLED`)
    - Best offers (`MARKETPLACE_ACCOUNT.OFFER.UPDATED`)
  - Subscribe via `/webhook-subscribe` command
  - **NOTE:** Offer notifications require setup - see "How to Enable Offer Notifications" below

- **📊 Status Monitoring**
  - Check eBay API connection with `/ebay-status`
  - View OAuth scopes with `/ebay-scopes`
  - Token validation and environment info

### ⚠️ Limited / Requires Setup

- **💰 Offer Management**
  - Accept offers with `/accept-offer`
  - Counter offers with `/counter-offer`
  - Decline offers with `/decline-offer`
  - **NOTE:** Viewing offers requires webhook subscription
  - Run `/get-offers` for setup instructions

- **📦 Listing Viewer**
  - `/get-listings` - Links to eBay's active listings page
  - **LIMITATION:** eBay API requires complex XML parsing for traditional listings
  - Best managed through eBay website

### ❌ Currently Not Available

- **💵 Financial Information** (Requires eBay Developer Approval)
  - `/get-balance` - Returns 404
  - `/get-payouts` - Returns 404
  - **REASON:** Finances API scope not granted by eBay
  - **SOLUTION:** Contact eBay Developer Support to enable Finances API for your keyset
  - Your account DOES have Managed Payments enabled
  - Your keyset DOES have Finances API enabled in Developer Portal
  - Issue: OAuth token doesn't include `sell.finances` scope

- **💬 Buyer Messages** (API Limitation)
  - **REASON:** Post-Order API doesn't support OAuth Bearer tokens
  - Requires legacy IAF authentication (being phased out by eBay)
  - **WORKAROUND:** Check messages on eBay.com

### ❌ Out of Scope

- **Listing Creation**: Not feasible through Discord due to eBay's complex listing requirements (item specifics, categories, shipping policies) and Discord's image compression. Use eBay's web interface or dedicated listing tools instead.

## 📋 Prerequisites

- Go 1.25.6 or higher
- Discord Bot Token
- eBay Developer Account with **Production** API credentials
- Public domain with HTTPS for webhook endpoint (e.g., jacob.it.com)
- Ubuntu/Linux server for deployment (Windows cross-compilation supported)

## 🚀 Production Deployment

**Current Server Configuration:**
- **OS:** Ubuntu Server 22.04.5 LTS
- **Server:** jacob.it.com (192.168.0.12)
- **Service:** systemd service `ebay-bot.service`
- **Binary:** Cross-compiled from Windows using `GOOS=linux GOARCH=amd64`
- **Logs:** `/home/jacob/ebay-bot/bot.log` and `/home/jacob/ebay-bot/bot-error.log`
- **Webhook:** Nginx reverse proxy to port 8081 with Let's Encrypt SSL

### Deployment Script

Use the automated deployment script:
```powershell
.\scripts\deploy.ps1
```

Or manually:
```powershell
# Build for Linux
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o bin/ebaymanager-bot-linux

# Deploy to server
scp bin/ebaymanager-bot-linux jacob@192.168.0.12:/home/jacob/ebay-bot/
ssh jacob@192.168.0.12 "chmod +x /home/jacob/ebay-bot/ebaymanager-bot-linux && sudo systemctl restart ebay-bot"

# Check status
ssh jacob@192.168.0.12 "sudo systemctl status ebay-bot"
```

### Service Management

```bash
# Start the bot
sudo systemctl start ebay-bot

# Stop the bot
sudo systemctl stop ebay-bot

# Restart the bot
sudo systemctl restart ebay-bot

# Check status
sudo systemctl status ebay-bot

# View logs
tail -f /home/jacob/ebay-bot/bot-error.log
```

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
- `/ebay-scopes` - View current OAuth scopes granted by eBay
- `/ebay-authorize` - **Automatic OAuth!** Get authorization URL - just click and sign in
- `/ebay-code` - Manual authorization code submission (backup method)

### Order Management
- `/get-orders [limit]` - ✅ View recent eBay orders with buyer details and images

### Offer Management
- `/get-offers` - View instructions for setting up webhook offer notifications
- `/accept-offer <offer_id>` - Accept a buyer's offer
- `/counter-offer <offer_id> <price>` - Counter an offer with a different price
- `/decline-offer <offer_id>` - Decline a buyer's offer

### Listing Management
- `/get-listings` - Link to eBay's active listings page (API limitations prevent direct listing)

### Financial Information (⚠️ Requires eBay Developer Approval)
- `/get-balance` - ❌ Currently returns 404 - Awaiting Finances API scope approval
- `/get-payouts` - ❌ Currently returns 404 - Awaiting Finances API scope approval

### Communication
- `/get-messages` - ❌ Not supported (Post-Order API doesn't support OAuth)

### Webhook Management  
- `/webhook-subscribe` - ✅ Subscribe to eBay real-time notifications (ORDER, OFFER)
- `/webhook-list` - ✅ List active webhook subscriptions
- `/webhook-test` - ✅ Test webhook notification to current channel

## 🔧 How to Enable Offer Notifications

Since eBay's API doesn't provide a direct endpoint to list offers, you'll receive them via webhooks:

1. Run `/webhook-subscribe` in Discord
2. Select **OFFER** from the notification types
3. When a buyer makes an offer, you'll get an instant Discord notification with:
   - Offer ID
   - Item details
   - Buyer's offer amount
   - Your list price
4. Use `/accept-offer`, `/counter-offer`, or `/decline-offer` to respond

## ⚠️ Known Issues & Limitations

### Requires eBay Developer Support

**Financial APIs (404 errors):**
- Issue: `/get-balance` and `/get-payouts` return 404
- Cause: OAuth token doesn't include `sell.finances` scope
- Evidence: Your eBay keyset HAS Finances API enabled in Developer Portal
- Evidence: Your account HAS Managed Payments enabled
- Evidence: OAuth token response from eBay doesn't include any scope field
- **Solution:** Contact eBay Developer Support to request Finances API scope for your RuName
- Until resolved, check finances on eBay.com

**Messages API (401 errors):**
- Issue: `/get-messages` returns 401 "Bad scheme: Bearer"
- Cause: Post-Order API requires legacy IAF authentication, not OAuth Bearer tokens
- IAF authentication is being phased out by eBay
- **Solution:** None available - check messages on eBay.com

### API Design Limitations

**Listings:**
- Viewing listings requires Trading API with complex XML parsing
- Current implementation provides link to eBay's active listings page
- Traditional eBay listings (non-inventory) don't work well with REST APIs

**Offers:**
- No direct API endpoint to list all offers
- Requires webhook subscription for real-time notifications
- `/get-offers` provides setup instructions

## 🐛 Troubleshooting

### OAuth Issues

**"Authorization Successful" but APIs still fail:**
- Check granted scopes with `/ebay-scopes`
- If `sell.finances` is missing, contact eBay Developer Support
- RuName must be configured with webhook callback URLs:
  - Auth accepted: `https://jacob.it.com/webhook/oauth/callback`
  - Auth declined: `https://jacob.it.com/webhook/oauth/declined`

**Bot doesn't respond to OAuth callback:**
- Ensure RuName is configured correctly in eBay Developer Portal
- Verify webhook server is running: `sudo systemctl status ebay-bot`
- Check logs: `tail -f /home/jacob/ebay-bot/bot-error.log`
- Look for "OAuth callback received" message

### Webhook Issues

**Not receiving offer notifications:**
1. Check subscription status with `/webhook-list`
2. Verify webhook endpoint is accessible: `https://jacob.it.com/webhook/health`
3. Re-subscribe with `/webhook-subscribe` if needed
4. Check webhook logs for incoming notifications

### Financial API 404 Errors

**Even after re-authorization:**
- This is an eBay permission issue, not a bot issue
- Your keyset needs Finances API scope approval from eBay
- Use `/ebay-scopes` to confirm which scopes you have
- Contact api@ebay.com with your:
  - App ID (Client ID)
  - RuName
  - Request for Finances API scope approval

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
- [x] Production deployment on Ubuntu Server
- [x] Systemd service configuration
- [x] Cross-compilation from Windows to Linux

### ✅ Phase 2: eBay OAuth & Authentication (COMPLETE)
- [x] OAuth 2.0 authorization code flow
- [x] **Automatic token exchange via webhook server**
- [x] Automatic token refresh (every 90 minutes)
- [x] Secure token storage in .env file
- [x] Discord command integration for auth
- [x] OAuth scope diagnostics (`/ebay-scopes`)

### ✅ Phase 3: Basic eBay Operations (COMPLETE)
- [x] Fetch orders from Fulfillment API with image support
- [x] Browse API integration for item images
- [x] Image URL fallback generation
- [x] Fetch offers guidance (webhook-based)
- [x] Accept/decline/counter offers via Negotiation API

### ✅ Phase 4: Webhook Integration (COMPLETE)
- [x] Webhook HTTP server with challenge verification
- [x] Production HTTPS endpoint with Nginx reverse proxy
- [x] Real-time Discord notifications for orders and offers
- [x] Subscription management via Discord commands
- [x] OAuth callback integration in webhook server
- [x] SHA-256 verification for eBay notifications

### ⏸️ Phase 5: Financial APIs (BLOCKED - Awaiting eBay Approval)
- [ ] Account balance (API: 404 - scope not granted)
- [ ] Payout history (API: 404 - scope not granted)
- [ ] Transaction details (API: 404 - scope not granted)
- **Blocker:** Requires `sell.finances` scope approval from eBay Developer Support

### ❌ Phase 6: Messages API (BLOCKED - API Limitation)
- [ ] Buyer messages (API: 401 - OAuth not supported)
- **Blocker:** Post-Order API requires legacy IAF authentication

### 🚧 Phase 7: Advanced Features (PLANNED)
- [ ] Shipping label integration
- [ ] Analytics and sales reporting
- [ ] Automated offer response rules
- [ ] Performance metrics dashboard

## 📊 eBay API Endpoints

### ✅ Working in Production
- **OAuth 2.0 API**: `/identity/v1/oauth2/token` - Token exchange & refresh ✅
- **Fulfillment API**: `/sell/fulfillment/v1/order` - Fetch orders ✅
- **Browse API**: `/buy/browse/v1/item/get_item_by_legacy_id` - Item images ✅
- **Negotiation API**: `/sell/negotiation/v1/offer/{id}/respond` - Manage offers ✅
- **Marketplace Notification API**: `/commerce/notification/v1/destination` - Webhooks ✅

### ❌ Blocked - Requires eBay Developer Approval
- **Finances API**: `/sell/finances/v1/seller_funds_summary` - 404 (scope not granted) ❌
- **Finances API**: `/sell/finances/v1/payout` - 404 (scope not granted) ❌
- **Finances API**: `/sell/finances/v1/transaction` - 404 (scope not granted) ❌

### ❌ Not Supported - API Limitation
- **Post-Order API**: `/post-order/v2/inquiry/search` - 401 (OAuth not supported) ❌

### 🚧 Not Implemented
- **Buy API**: Shipping labels
- **Analytics API**: Performance metrics
- **Inventory API**: Listing management (complex XML required)

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
QUICKSTART.md](QUICKSTART.md) - Quick setup guide
- [WEBHOOK_SETUP.md](WEBHOOK_SETUP.md) - Webhook configuration details
- [PRODUCTION_READINESS.md](PRODUCTION_READINESS.md) - Production deployment checklist
- [QUICK_TEST_GUIDE.md](QUICK_TEST
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
