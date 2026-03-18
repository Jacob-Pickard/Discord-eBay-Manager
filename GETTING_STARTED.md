# Getting Started

Complete guide to set up and run the eBay Manager Discord Bot.

## 📋 Prerequisites

- **Go 1.25.6+** - [Download here](https://go.dev/dl/)
- **Discord Account** - For creating the bot
- **eBay Developer Account** - [Sign up here](https://developer.ebay.com/)
- **(Optional) Linux Server** - For production deployment with webhooks

---

## 🚀 Quick Start (Local Development)

### Step 1: Get Discord Bot Token (5 minutes)

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click **New Application** → Name it "eBay Manager Bot"
3. Go to **Bot** tab → Click **Add Bot**
4. Copy the bot token (you'll need this for `.env`)
5. Enable **MESSAGE CONTENT INTENT** and **SERVER MEMBERS INTENT**
6. Go to **OAuth2 → URL Generator**:
   - **Scopes**: Select `bot` and `applications.commands`
   - **Permissions**: Select `Send Messages`, `Embed Links`, `Use Slash Commands`
7. Copy the generated URL and open it to invite the bot to your server

### Step 2: Get eBay API Credentials (10 minutes)

1. Go to [eBay Developer Portal](https://developer.ebay.com/my/keys)
2. Click **Create Application Keys** (if you haven't already)
3. For testing, use the **SANDBOX** tab
4. Note down your:
   - **App ID (Client ID)**
   - **Cert ID (Client Secret)**
   - **Dev ID**

**For Production:**
- Switch to the **PRODUCTION** tab and get production keys
- You'll also need a **RuName (Redirect URI)** - create one in the Application Keys section

### Step 3: Configure Environment

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your credentials:**
   ```env
   # Discord Configuration
   DISCORD_BOT_TOKEN=your_discord_bot_token_here
   NOTIFICATION_CHANNEL_ID=your_discord_channel_id_here
   
   # eBay API Configuration
   EBAY_APP_ID=your_ebay_app_id
   EBAY_CERT_ID=your_ebay_cert_id
   EBAY_DEV_ID=your_ebay_dev_id
   EBAY_REDIRECT_URI=your_ebay_runame
   
   # OAuth Tokens (leave blank - will be generated via /ebay-authorize)
   EBAY_ACCESS_TOKEN=
   EBAY_REFRESH_TOKEN=
   
   # Environment (SANDBOX for testing, PRODUCTION for live)
   EBAY_ENVIRONMENT=SANDBOX
   
   # Webhook Configuration
   WEBHOOK_PORT=8081
   WEBHOOK_VERIFY_TOKEN=generate_a_secure_random_token_here
   ```

3. **Get your Discord Channel ID:**
   - Enable Developer Mode in Discord (User Settings → Advanced → Developer Mode)
   - Right-click your channel → Copy ID
   - Paste into `NOTIFICATION_CHANNEL_ID`

### Step 4: Install Dependencies

```bash
go mod download
```

### Step 5: Run the Bot

```bash
go run main.go
```

You should see:
```
Connected as: eBay Manager Bot#1234 (ID: ...)
Webhook server listening on :8081
Command registration complete!
```

### Step 6: Authorize eBay Account

1. In Discord, type `/ebay-authorize`
2. Click the authorization link
3. Sign in to your eBay account (sandbox or production)
4. Grant permissions
5. The bot will automatically exchange the code for tokens

### Step 7: Test the Bot

Try these commands in Discord:
- `/ebay-status` - Check connection and token status
- `/get-orders` - View recent orders
- `/ebay-scopes` - See what permissions your token has

---

## 🌐 Production Deployment (Optional)

For production use with real-time webhook notifications, you'll need a public server.

### Prerequisites for Production

- Linux server (Ubuntu 22.04+ recommended)
- Domain name with HTTPS (Let's Encrypt recommended)
- SSH access to your server

### Step 1: Configure Deployment

1. **Copy deployment configuration:**
   ```bash
   cp deploy-config.env.example deploy-config.env
   ```

2. **Edit `deploy-config.env`:**
   ```env
   DEPLOY_SERVER_IP=your.server.ip.address
   DEPLOY_SERVER_USER=your_ssh_username
   DEPLOY_SERVER_PATH=/path/to/deployment/directory
   DEPLOY_DOMAIN=yourdomain.com
   DEPLOY_WEBHOOK_URL=https://yourdomain.com/webhook/ebay/notification
   ```

### Step 2: Set Up Server

1. **Create systemd service:**
   ```bash
   # Copy example to server
   scp config/ebay-bot.service.example user@server:/tmp/
   
   # SSH to server and edit
   ssh user@server
   sudo nano /tmp/ebay-bot.service.example
   # Update User, WorkingDirectory, and paths
   sudo mv /tmp/ebay-bot.service.example /etc/systemd/system/ebay-bot.service
   sudo systemctl daemon-reload
   sudo systemctl enable ebay-bot
   ```

2. **Configure Nginx:**
   ```bash
   # Copy and edit nginx config
   scp config/webhook-domain.conf.example user@server:/tmp/
   ssh user@server
   sudo nano /tmp/webhook-domain.conf.example
   # Replace yourdomain.com with your actual domain
   sudo mv /tmp/webhook-domain.conf.example /etc/nginx/sites-available/yourdomain.com
   sudo ln -s /etc/nginx/sites-available/yourdomain.com /etc/nginx/sites-enabled/
   sudo nginx -t
   sudo systemctl reload nginx
   ```

3. **Set up SSL with Let's Encrypt:**
   ```bash
   ssh user@server
   sudo apt install certbot python3-certbot-nginx
   sudo certbot --nginx -d yourdomain.com
   ```

### Step 3: Deploy

```powershell
# Build and deploy to server
.\scripts\deploy.ps1
```

The script will:
1. Build Linux binary
2. Upload to your server  
3. Restart the bot service
4. Test the webhook endpoint
5. Show deployment status

### Step 4: Set Up Webhooks

1. In Discord: `/webhook-subscribe url:https://yourdomain.com/webhook/ebay/notification`
2. The bot will create the eBay webhook subscription
3. You'll now receive real-time notifications for:
   - New orders
   - Best offers
   - Account changes

---

## 📁 Configuration Files Reference

| File | Purpose | Committed to Git? |
|------|---------|-------------------|
| `.env` | Your actual secrets | ❌ No (gitignored) |
| `.env.example` | Template with placeholders | ✅ Yes |
| `deploy-config.env` | Your server details | ❌ No (gitignored) |
| `deploy-config.env.example` | Deployment template | ✅ Yes |
| `config/ebay-bot.service` | Your service file | ❌ No (gitignored) |
| `config/ebay-bot.service.example` | Service template | ✅ Yes |

---

## 🐛 Troubleshooting

### Bot won't start
- ✅ Check `.env` file exists in project root
- ✅ Verify `DISCORD_BOT_TOKEN` is set correctly
- ✅ Make sure Go dependencies are installed: `go mod download`

### Commands don't appear in Discord
- ✅ Wait 10-15 seconds after bot starts
- ✅ Make sure bot has `applications.commands` permission
- ✅ Try typing `/` in Discord to see registered commands
- ✅ Restart Discord if needed

### "DISCORD_BOT_TOKEN not set"
- ✅ Create `.env` from `.env.example`
- ✅ Add your Discord bot token to `.env`

### OAuth authorization fails
- ✅ Check `EBAY_REDIRECT_URI` matches your eBay RuName exactly
- ✅ For production, ensure your domain is accessible
- ✅ Try sandbox environment first for testing

### Webhook errors
- ✅ Check domain is accessible: `curl https://yourdomain.com/webhook/health`
- ✅ Verify SSL certificate is valid
- ✅ Check bot logs: `ssh user@server "journalctl -u ebay-bot -f"`
- ✅ Ensure webhook verify token is 32-80 characters

### Deployment fails
- ✅ Verify `deploy-config.env` is configured correctly
- ✅ Test SSH connection: `ssh user@server`
- ✅ Check server has enough disk space
- ✅ Verify systemd service is configured correctly

---

## 🎯 What to Do Next

**After Local Setup:**
1. Run `/ebay-authorize` to connect your eBay account
2. Test with `/ebay-status` and `/get-orders`
3. Try offer management commands with test listings

**For Production:**
1. Get production eBay credentials
2. Set up your server and domain
3. Deploy with `.\scripts\deploy.ps1`
4. Subscribe to webhooks with `/webhook-subscribe`

---

## 📚 Additional Resources

- [README.md](README.md) - Complete feature documentation
- [SECURITY.md](SECURITY.md) - Security best practices
- [docs/](docs/) - Detailed deployment guides
- [eBay API Documentation](https://developer.ebay.com/docs)
- [Discord.js Guide](https://discord.com/developers/docs)

---

## 💡 Pro Tips

- **Use SANDBOX first** - Test with eBay's sandbox environment before going live
- **Secure tokens** - Generate strong webhook tokens: `openssl rand -base64 48 | tr -d "=+/" | cut -c1-60`
- **Monitor logs** - Use `.\scripts\deploy-watch.ps1` to deploy and watch logs in real-time
- **Backup configs** - Keep your `.env` and `deploy-config.env` backed up securely (encrypted)
- **Update regularly** - Pull latest changes and redeploy: `git pull && .\scripts\deploy.ps1`

---

Need help? Check the troubleshooting section or review the documentation in the `docs/` folder! 🚀
