# eBay Bot Production Deployment - Complete! ✅

## Deployment Summary

**Date:** February 12, 2026  
**Server:** jacob.it.com (192.168.0.12)  
**Status:** ✅ **DEPLOYED AND RUNNING**

---

## What Was Deployed

### 1. Bot Application
- **Location:** `/home/jacob/ebay-bot/`
- **Executable:** `ebaymanager-bot-linux`
- **Config:** `.env` (sandbox configuration)
- **Status:** Running as systemd service

### 2. Systemd Service
- **Service Name:** `ebay-bot.service`
- **Status:** Active and running
- **Auto-start:** Enabled (starts on boot)
- **User:** jacob
- **Logs:**
  - Standard output: `/home/jacob/ebay-bot/bot.log`
  - Error output: `/home/jacob/ebay-bot/bot-error.log`

### 3. Nginx Configuration
- **Config File:** `/etc/nginx/conf.d/jacob.it.com.conf`
- **Webhook Endpoint:** `https://jacob.it.com/webhook/`
- **Proxy Target:** `http://127.0.0.1:8081`
- **SSL:** Enabled (Let's Encrypt)
- **Status:** Active and configured

### 4. Bot Status
- **Discord Bot:** Connected as "eBay Manager#7897"
- **Commands:** 14 registered successfully
- **Webhook Server:** Listening on port 8081
- **Health Check:** ✅ `https://jacob.it.com/webhook/health` responding

---

## Webhook Endpoints

Your bot is now accessible via these public HTTPS URLs:

1. **Main webhook:** `https://jacob.it.com/webhook/ebay/notification`
2. **Challenge verification:** `https://jacob.it.com/webhook/ebay/challenge`
3. **Health check:** `https://jacob.it.com/webhook/health`

---

## Useful Commands

### Check Bot Status
```bash
ssh jacob@192.168.0.12 "sudo systemctl status ebay-bot"
```

### View Live Logs
```bash
ssh jacob@192.168.0.12 "tail -f /home/jacob/ebay-bot/bot-error.log"
```

### Restart Bot
```bash
ssh jacob@192.168.0.12 "sudo systemctl restart ebay-bot"
```

### Stop Bot
```bash
ssh jacob@192.168.0.12 "sudo systemctl stop ebay-bot"
```

### Start Bot
```bash
ssh jacob@192.168.0.12 "sudo systemctl start ebay-bot"
```

### Check if Bot is Listening
```bash
ssh jacob@192.168.0.12 "ss -tlnp | grep 8081"
```

### Test Webhook from Windows
```powershell
Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health"
```

---

## Current Configuration Status

### Environment
- **Current Mode:** SANDBOX
- **eBay Environment:** Sandbox testing
- **Discord:** Connected
- **OAuth Tokens:** Valid (sandbox)

### What's Working
- ✅ Bot deployed and running 24/7
- ✅ HTTPS webhook accessible publicly
- ✅ All Discord commands registered
- ✅ Auto-restarts on failure
- ✅ Starts automatically on server reboot

---

## Next Steps for Production

### Option A: Switch to Production NOW
If you have production credentials ready:

1. **SSH into server:**
   ```bash
   ssh jacob@192.168.0.12
   ```

2. **Edit .env file:**
   ```bash
   cd /home/jacob/ebay-bot
   nano .env
   ```

3. **Update these values:**
   ```env
   EBAY_ENVIRONMENT=PRODUCTION
   EBAY_APP_ID=<production_app_id>
   EBAY_CERT_ID=<production_cert_id>
   EBAY_DEV_ID=<production_dev_id>
   EBAY_REDIRECT_URI=<production_runame>
   # Clear tokens (will be regenerated)
   EBAY_ACCESS_TOKEN=
   EBAY_REFRESH_TOKEN=
   ```

4. **Restart bot:**
   ```bash
   sudo systemctl restart ebay-bot
   ```

5. **In Discord:**
   ```
   /ebay-authorize
   /ebay-code code:<your_code>
   /webhook-subscribe url:https://jacob.it.com/webhook/ebay/notification
   ```

### Option B: Keep Testing in Sandbox
Your bot is running in sandbox mode and ready for testing:
- All commands work
- Webhooks are configured
- Can test with sandbox offers/orders

---

## Production Credential Checklist

**Still need:**
- [ ] Production App ID from eBay Developer Portal
- [ ] Production Cert ID from eBay Developer Portal
- [ ] Production Dev ID from eBay Developer Portal
- [ ] Production RuName (may take 1-3 days approval)

**Get these from:** https://developer.ebay.com/my/keys (PRODUCTION tab)

**See:** [GET_PRODUCTION_CREDENTIALS.md](GET_PRODUCTION_CREDENTIALS.md) for detailed instructions

---

## Troubleshooting

### Bot Won't Start
```bash
# Check logs
ssh jacob@192.168.0.12 "tail -50 /home/jacob/ebay-bot/bot-error.log"

# Check service status
ssh jacob@192.168.0.12 "sudo systemctl status ebay-bot"
```

### Webhook Returns 502
```bash
# Check if bot is running
ssh jacob@192.168.0.12 "systemctl is-active ebay-bot"

# Check if listening on port
ssh jacob@192.168.0.12 "ss -tlnp | grep 8081"
```

### Need to Update Config
```bash
# SSH in
ssh jacob@192.168.0.12

# Edit config
cd /home/jacob/ebay-bot
nano .env

# Restart bot
sudo systemctl restart ebay-bot
```

---

## Important Notes

1. **Bot runs 24/7** - No need to keep Windows on
2. **Auto-restarts** - Will restart if it crashes
3. **Logs are persistent** - Check `/home/jacob/ebay-bot/bot-error.log`
4. **Sandbox mode active** - Currently using sandbox eBay API
5. **Ready for production** - Just need production credentials

---

## File Locations

### On Server (192.168.0.12)
- Bot directory: `/home/jacob/ebay-bot/`
- Executable: `/home/jacob/ebay-bot/ebaymanager-bot-linux`
- Config: `/home/jacob/ebay-bot/.env`
- Logs: `/home/jacob/ebay-bot/bot-error.log`
- Service file: `/etc/systemd/system/ebay-bot.service`
- Nginx config: `/etc/nginx/conf.d/jacob.it.com.conf`

### On Windows (Backup)
- Sandbox backup: `.env.sandbox.backup`
- Production template: `.env.production.template`
- Linux binary: `ebaymanager-bot-linux`
- Windows binary: `ebaymanager-bot.exe`

---

## Success! 🎉

Your eBay Discord bot is now:
- ✅ Running on your Proxmox server
- ✅ Accessible via HTTPS at jacob.it.com
- ✅ Set to auto-start on boot
- ✅ Auto-restarts on failures
- ✅ Ready for production (pending credentials)

**The bot is LIVE and fully operational in sandbox mode!**
