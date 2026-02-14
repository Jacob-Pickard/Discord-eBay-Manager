# eBay Bot Deployment Scripts

## Quick Start

### Deploy Changes to Server
After making code changes, just run:
```powershell
.\deploy.ps1
```

This will:
- Build Linux binary
- Upload to jacob.it.com
- Restart bot service
- Show status and logs

---

## Available Scripts

### 1. `deploy.ps1` - Full Deployment
**Use when:** You changed code and want to deploy

```powershell
.\deploy.ps1
```

Builds and deploys your bot to the server.

---

### 2. `deploy-quick.ps1` - Quick Deploy (Shortcut)
**Use when:** You want the shortest command

```powershell
.\deploy-quick.ps1
```

Same as `deploy.ps1` but shorter to type.

---

### 3. `deploy-watch.ps1` - Deploy + Watch Logs
**Use when:** You want to see live logs after deploying

```powershell
.\deploy-watch.ps1
```

Deploys and then streams live logs from the server.

---

### 4. `deploy-config.ps1` - Update Configuration Only
**Use when:** You only changed .env settings (no code changes)

```powershell
# Update from default .env
.\deploy-config.ps1

# Update from production config
.\deploy-config.ps1 -LocalEnvFile .env.production

# Update from custom file
.\deploy-config.ps1 -LocalEnvFile my-custom.env
```

Uploads new .env and restarts bot (much faster than full deployment).

---

### 5. `setup-production.ps1` - Get Production Credentials
**Use when:** You need to get production eBay API credentials

```powershell
.\setup-production.ps1
```

Interactive wizard that:
- Opens eBay Developer Portal
- Guides you through getting credentials
- Creates .env.production file
- Optionally deploys to production

---

## Common Workflows

### Making Code Changes
```powershell
# 1. Edit your code in VS Code
# 2. Deploy to server
.\deploy.ps1
```

### Switching to Production
```powershell
# 1. Get production credentials
.\setup-production.ps1

# 2. Bot is automatically updated
# 3. Authorize in Discord:
#    /ebay-authorize
#    /ebay-code code:<your_code>
#    /webhook-subscribe url:https://jacob.it.com/webhook/ebay/notification
```

### Updating Environment Variables
```powershell
# 1. Edit .env file
# 2. Upload changes
.\deploy-config.ps1
```

### View Live Logs
```powershell
# Option 1: Deploy and watch
.\deploy-watch.ps1

# Option 2: SSH directly
ssh jacob@192.168.0.12 "tail -f /home/jacob/ebay-bot/bot-error.log"
```

---

## Server Management Commands

### Check Bot Status
```powershell
ssh jacob@192.168.0.12 "sudo systemctl status ebay-bot"
```

### Restart Bot Manually
```powershell
ssh jacob@192.168.0.12 "sudo systemctl restart ebay-bot"
```

### Stop Bot
```powershell
ssh jacob@192.168.0.12 "sudo systemctl stop ebay-bot"
```

### Start Bot
```powershell
ssh jacob@192.168.0.12 "sudo systemctl start ebay-bot"
```

### View Logs (Last 50 Lines)
```powershell
ssh jacob@192.168.0.12 "tail -50 /home/jacob/ebay-bot/bot-error.log"
```

### Test Webhook
```powershell
Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health"
```

---

## File Locations

### On Windows (Your PC)
- Source code: Current directory
- Sandbox config: `.env.sandbox.backup`
- Production config: `.env.production` (created by setup-production.ps1)
- Deployment scripts: `deploy*.ps1`

### On Server (jacob.it.com)
- Bot directory: `/home/jacob/ebay-bot/`
- Executable: `/home/jacob/ebay-bot/ebaymanager-bot-linux`
- Configuration: `/home/jacob/ebay-bot/.env`
- Logs: `/home/jacob/ebay-bot/bot-error.log`
- Service: `/etc/systemd/system/ebay-bot.service`

---

## Troubleshooting

### Deployment fails with password prompt
**Solution:** Make sure you're entering your SSH password for jacob@192.168.0.12

### Bot won't start after deployment
```powershell
# Check logs
ssh jacob@192.168.0.12 "tail -50 /home/jacob/ebay-bot/bot-error.log"

# Check service status
ssh jacob@192.168.0.12 "sudo systemctl status ebay-bot"
```

### Webhook returns 502
```powershell
# Verify bot is running
ssh jacob@192.168.0.12 "systemctl is-active ebay-bot"

# Check if listening on port
ssh jacob@192.168.0.12 "ss -tlnp | grep 8081"
```

### Changes not appearing
Make sure you ran `.\deploy.ps1` to build and upload the new binary.

---

## Tips

1. **Always deploy after code changes** - The server won't see changes until you deploy
2. **Use deploy-config.ps1 for .env changes** - Much faster than full deployment
3. **Watch logs after deploying** - Use `deploy-watch.ps1` to see if bot started correctly
4. **Test in sandbox first** - Make sure features work before switching to production

---

## Need Help?

See these files for more info:
- [DEPLOYMENT_COMPLETE.md](DEPLOYMENT_COMPLETE.md) - Server deployment details
- [GET_PRODUCTION_CREDENTIALS.md](GET_PRODUCTION_CREDENTIALS.md) - Manual credential guide
- [NGINX_WEBHOOK_SETUP.md](NGINX_WEBHOOK_SETUP.md) - Webhook configuration
