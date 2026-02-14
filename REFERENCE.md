# 🎮 eBay Discord Bot - Quick Reference

## 🚀 Deployment (Passwordless!)

```powershell
# Deploy to production (one command!)
.\scripts\deploy.ps1
```

**What it does:**
- Builds Linux binary to `bin/`
- Uploads to jacob.it.com
- Restarts bot service
- Shows logs

**No passwords needed!** SSH keys are configured.

---

## 🗂️ Project Structure

```
├── bin/            # Compiled binaries
├── config/         # Service configs (.service, nginx)
├── docs/           # Documentation
├── scripts/        # Deployment scripts
├── logs/           # Log files (if running locally)
├── internal/       # Go source code
│   ├── bot/       # Discord bot handlers
│   ├── config/    # Configuration loading
│   ├── ebay/      # eBay API client
│   └── webhook/   # Webhook server
├── tools/          # Utility tools
├── main.go         # Entry point
└── .env            # Environment variables
```

---

## 🔐 OAuth Authorization (Now AUTOMATIC!)

### New Way (Automatic):
```
/ebay-authorize
```
1. Click the link
2. Sign in to eBay
3. Click "Agree"
4. **Done!** The bot automatically detects your authorization

No code copying needed! The bot runs a local server that captures the OAuth callback.

### Old Way (Manual - still works):
```
/ebay-code code:YOUR_VERY_LONG_CODE
```

---

## 📊 Available Discord Commands

### eBay Data Commands:
- `/get-orders` - View recent orders (now shows buyer, price, status!)
- `/get-balance` - Check account balance
- `/get-payouts` - View payout history
- `/get-messages` - View buyer messages
- `/get-offers` - Check for offers (webhook-based)

### Authorization:
- `/ebay-status` - Check API connection
- `/ebay-authorize` - Authorize bot (automatic OAuth!)
- `/ebay-code` - Manual code entry

### Webhooks:
- `/webhook-subscribe` - Enable real-time notifications
- `/webhook-list` - List active subscriptions
- `/webhook-test` - Test webhook endpoint

### Offer Management:
- `/accept-offer` - Accept a buyer offer
- `/counter-offer` - Counter with different price
- `/decline-offer` - Decline an offer

---

## 🔧 What Was Fixed

### ✅ File Organization
Before: 40+ files in root directory 😱
After: Clean structure with organized folders 🎯

### ✅ Passwordless Deployment
Before: Password prompts every 5 seconds 🔑🔑🔑
After: Zero passwords, instant deployment! ⚡

### ✅ Automatic OAuth
Before: Manual code copy/paste
After: Click link → Sign in → Done!  

### ✅ API Fixes
- **Orders**: Now parses buyer, price, and status correctly
- **Balance/Payouts**: Proper error handling (404 = no data yet)
- **Offers**: Explained that offers come via webhooks, not API endpoint

---

## 🌐 Production URLs

- **Bot Server**: https://jacob.it.com
- **Webhook Health**: https://jacob.it.com/webhook/health
- **Webhook Endpoint**: https://jacob.it.com/webhook/ebay/notification

---

## 🐛 Troubleshooting

### Bot not responding?
```powershell
# Check status
ssh jacob@192.168.0.12 "systemctl status ebay-bot"

# View logs
ssh jacob@192.168.0.12 "tail -f /home/jacob/ebay-bot/bot-error.log"

# Restart
ssh jacob@192.168.0.12 "sudo systemctl restart ebay-bot"
```

### Webhook not working?
```powershell
# Test endpoint
Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health"

# Should return: "OK"
```

### Need to rebuild?
```powershell
# Full deployment
.\scripts\deploy.ps1

# Just update config
.\scripts\deploy-config.ps1  

# Deploy and watch logs
.\scripts\deploy-watch.ps1
```

---

## 📝 Notes

### About API 404 Errors:
- **Balance/Payouts 404**: Normal if account has no transactions yet
- **Offers 404**: Correct! Use webhooks for offer notifications

### OAuth Scopes Enabled:
- `sell.inventory` - Listings and inventory
- `sell.fulfillment` - Orders and shipping
- `sell.account` - Account information
- `sell.finances` - Balance and payouts
- `sell.marketing` - Marketing campaigns

### Token Auto-Refresh:
Tokens automatically refresh every 90 minutes. You never need to re-authorize unless you revoke access.

---

## 🎉 That's It!

Your bot is production-ready with:
- ✅ Clean, organized codebase
- ✅ Passwordless deployment
- ✅ Automatic OAuth
- ✅ Fixed API parsing
- ✅ Real-time webhooks
- ✅ Running 24/7 on jacob.it.com

Happy selling! 🚀
