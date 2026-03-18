# eBay Manager Discord Bot

> A production-ready Discord bot for managing your eBay business operations—all from Discord.

[![Go Version](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows-lightgrey)]()

## 🎯 Overview

Manage your eBay seller account directly from Discord! View orders, respond to offers, and receive real-time notifications through webhooks—no need to constantly check eBay.com.

**Perfect for:** Small to medium eBay sellers who want to streamline their workflow and respond to customers faster.

---

## ✨ Features

### 📦 Order Management
- **View Recent Orders** - Get detailed order information with buyer details
- **Automatic Image Loading** - Uses eBay Browse API for product images
- **Rich Discord Embeds** - Beautiful, formatted order displays

### 💰 Offer Management
- **Accept Offers** - Instantly accept best offers
- **Counter Offers** - Send counteroffers with custom amounts
- **Decline Offers** - Politely decline with optional messages
- **Real-time Notifications** - Get notified when offers come in

### 🔔 Webhook Notifications (Production)
- **New Orders** - Instant notification when items sell
- **Best Offers** - Alert when buyers submit offers
- **Account Events** - Stay updated on account changes
- **SHA-256 Verification** - Secure webhook challenge validation

### 🔐 Security & Authentication
- **OAuth 2.0 Flow** - Fully automatic eBay authorization
- **Auto Token Refresh** - Tokens refresh every 90 minutes
- **Secure Storage** - Credentials stored in `.env` files (never committed)
- **Production Ready** - Follows security best practices

---

## 🚀 Quick Start

### 1. Prerequisites

- Go 1.21 or higher
- Discord bot token ([Get one here](https://discord.com/developers/applications))
- eBay Developer credentials ([Sign up here](https://developer.ebay.com/))

### 2. Install

```bash
git clone https://github.com/yourusername/ebay-manager-bot.git
cd ebay-manager-bot
go mod download
```

### 3. Configure

```bash
cp .env.example .env
# Edit .env with your credentials
```

### 4. Run

```bash
go run main.go
```

### 5. Authorize

In Discord, type `/ebay-authorize` and follow the link to connect your eBay account.

**📚 Full setup guide:** See [GETTING_STARTED.md](GETTING_STARTED.md)

---

## 💻 Commands

| Command | Description | Example |
|---------|-------------|---------|
| `/ebay-authorize` | Connect eBay account via OAuth | `/ebay-authorize` |
| `/ebay-status` | Check connection and token status | `/ebay-status` |
| `/ebay-scopes` | View current OAuth permissions | `/ebay-scopes` |
| `/get-orders` | View recent orders (last 10) | `/get-orders` |
| `/get-offers` | View pending best offers | `/get-offers` |
| `/accept-offer` | Accept a best offer | `/accept-offer offer_id:12345` |
| `/counter-offer` | Send a counteroffer | `/counter-offer offer_id:12345 amount:50.00` |
| `/decline-offer` | Decline an offer | `/decline-offer offer_id:12345` |
| `/webhook-subscribe` | Enable real-time notifications | `/webhook-subscribe url:https://...` |
| `/webhook-test` | Test webhook endpoint | `/webhook-test` |

---

## 🏗️ Architecture

```
┌─────────────────┐
│  Discord User   │
└────────┬────────┘
         │ Slash Commands
         ▼
┌─────────────────┐      OAuth 2.0      ┌──────────────┐
│  Discord Bot    │◄────────────────────►│  eBay API    │
│  (Go)           │                      └──────────────┘
└────────┬────────┘
         │ HTTP
         ▼
┌─────────────────┐      HTTPS/TLS      ┌──────────────┐
│ Webhook Server  │◄────────────────────►│ eBay Webhooks│
│ (Port 8081)     │                      └──────────────┘
└─────────────────┘
```

**Key Components:**
- **Discord Bot** - Handles slash commands and user interactions
- **eBay API Client** - Manages OAuth and API requests
- **Webhook Server** - Receives real-time notifications from eBay
- **Config Management** - Environment-based configuration

---

## 🌐 Production Deployment

Deploy to a Linux server for 24/7 operation and webhook support.

### Quick Deploy

```powershell
# 1. Configure your server details
cp deploy-config.env.example deploy-config.env
# Edit deploy-config.env

# 2. Deploy
.\scripts\deploy.ps1

# 3. Watch logs (optional)
.\scripts\deploy.ps1 -Watch
```

### Requirements

- **OS:** Ubuntu 22.04+ (or any Linux with systemd)
- **Domain:** Public domain with HTTPS (Let's Encrypt recommended)
- **Port:** 8081 open internally for webhook server
- **Reverse Proxy:** Nginx (configuration provided)

**📚 Deployment guide:** See [docs/DEPLOYMENT_SCRIPTS.md](docs/DEPLOYMENT_SCRIPTS.md)

---

## 📁 Project Structure

```
ebay-manager-bot/
├── main.go                      # Application entry point
├── go.mod                       # Go dependencies
├── .env.example                 # Environment template
├── deploy-config.env.example    # Deployment template
│
├── internal/
│   ├── bot/                     # Discord bot handlers
│   ├── config/                  # Configuration management
│   ├── ebay/                    # eBay API client
│   └── webhook/                 # Webhook server
│
├── config/
│   ├── .env.example            # Config examples
│   ├── ebay-bot.service.example # Systemd service template
│   └── webhook-domain.conf.example # Nginx template
│
├── scripts/
│   ├── deploy.ps1              # Main deployment script
│   ├── deploy-watch.ps1        # Deploy with log viewing
│   └── deploy-config.ps1       # Update .env on server
│
├── docs/
│   ├── DEPLOYMENT_SCRIPTS.md   # Deployment guide
│   ├── WEBHOOK_SETUP.md        # Webhook configuration
│   └── GET_PRODUCTION_CREDENTIALS.md # eBay credentials guide
│
├── GETTING_STARTED.md          # Setup guide
├── SECURITY.md                 # Security best practices
└── README.md                   # This file
```

---

## 🔧 Configuration

### Environment Variables

All configuration is done through `.env` files:

```env
# Discord
DISCORD_BOT_TOKEN=your_token
NOTIFICATION_CHANNEL_ID=channel_id

# eBay API
EBAY_APP_ID=your_app_id
EBAY_CERT_ID=your_cert_id
EBAY_DEV_ID=your_dev_id
EBAY_REDIRECT_URI=your_runame
EBAY_ENVIRONMENT=PRODUCTION # or SANDBOX

# OAuth (auto-generated)
EBAY_ACCESS_TOKEN=
EBAY_REFRESH_TOKEN=

# Webhooks
WEBHOOK_PORT=8081
WEBHOOK_VERIFY_TOKEN=random_secure_token
```

**🔐 Security:** Never commit `.env` files! Use the `.env.example` template.

---

## 🛡️ Security

- ✅ **OAuth 2.0** - Secure eBay authentication
- ✅ **Environment Variables** - Credentials never hardcoded
- ✅ **Gitignore Protection** - Sensitive files automatically excluded
- ✅ **Token Rotation** - Automatic refresh every 90 minutes
- ✅ **Webhook Verification** - SHA-256 challenge verification
- ✅ **HTTPS Only** - All production endpoints use TLS

**📚 Security guide:** See [SECURITY.md](SECURITY.md)

---

## 📊 Development

### Build

```bash
# Local development
go run main.go

# Production binary (Linux)
GOOS=linux GOARCH=amd64 go build -o bin/ebaymanager-bot-linux

# Windows binary
go build -o bin/ebaymanager-bot.exe
```

### Testing

```bash
# Run tests
go test ./...

# Check configuration
go run tools/check_config.go
```

### Tools

- `tools/check_config.go` - Validate environment configuration
- `tools/test_webhook_subscription.go` - Test eBay webhook subscriptions
- `tools/Test-Webhook-Simple.ps1` - Simple webhook endpoint tests

---

## 🐛 Troubleshooting

### Common Issues

**Bot won't start**
- Check `.env` file exists and contains valid credentials
- Verify Go dependencies are installed: `go mod download`

**Commands don't appear**
- Wait 10-15 seconds after bot starts
- Ensure bot has `applications.commands` permission
- Restart Discord if needed

**OAuth fails**
- Verify `EBAY_REDIRECT_URI` matches your eBay RuName exactly
- Check you're using correct environment (SANDBOX vs PRODUCTION)

**Webhook errors**
- Test endpoint accessibility: `curl https://yourdomain.com/webhook/health`
- Verify SSL certificate is valid
- Check webhook verify token is 32-80 characters

**📚 More help:** See [GETTING_STARTED.md](GETTING_STARTED.md#troubleshooting)

---

## 📝 API Endpoints Used

| eBay API | Purpose | Status |
|----------|---------|--------|
| OAuth 2.0 | Token exchange & refresh | ✅ Working |
| Fulfillment API | Order management | ✅ Working |
| Inventory API | Offer management | ✅ Working |
| Notification API | Webhook subscriptions | ✅ Working |
| Browse API | Product images | ✅ Working |

---

## 🚦 Status

- ✅ **Production Ready** - Fully tested and deployed
- ✅ **Active Development** - Regular updates and improvements
- ✅ **Well Documented** - Comprehensive guides and examples
- ✅ **Secure** - Follows security best practices

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **eBay Developer Program** - For comprehensive API documentation
- **Discord Developer Portal** - For bot development tools
- **bwmarrin/discordgo** - Excellent Discord library for Go
- **joho/godotenv** - Environment configuration management

---

## 📧 Support

- 📖 **Documentation:** [docs/](docs/)
- 🐛 **Issues:** [GitHub Issues](https://github.com/yourusername/ebay-manager-bot/issues)
- 💬 **Discussions:** [GitHub Discussions](https://github.com/yourusername/ebay-manager-bot/discussions)

---

<div align="center">

**Made with ❤️ for eBay sellers who want to work smarter, not harder**

[Get Started](GETTING_STARTED.md) • [View Docs](docs/) • [Report Bug](https://github.com/yourusername/ebay-manager-bot/issues)

</div>
