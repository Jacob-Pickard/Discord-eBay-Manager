# Webhook Production Setup Guide

## 🚨 Critical Issue: Private IP Address

**You mentioned**: `192.168.0.12` for your Proxmox server  
**Problem**: This is a **private/local IP address** - eBay cannot reach it from the internet

eBay webhooks need to POST to a **public** endpoint accessible from anywhere on the internet.

---

## Finding Your Current Website Setup

Since you mentioned you're **"hosting my website"** on this Proxmox server, let's figure out how it's currently accessible:

### Check Your Public Setup

**Your website is likely accessible via:**
1. A public domain name (e.g., `yourdomain.com`)
2. Your public IP address 
3. A reverse proxy/SSL setup (likely Nginx or Caddy)

**To find your public domain/IP:**

```powershell
# Check what domain points to your server
nslookup yourdomain.com

# Or find your public IP
Invoke-RestMethod -Uri "https://api.ipify.org?format=text"
```

---

## Three Solutions for Production Webhooks

### ✅ Option 1: Use Your Existing Website Setup (RECOMMENDED)

**If your website is already publicly accessible with HTTPS:**

You already have everything you need! Just add a webhook endpoint to your server.

**Example Setup:**
- Your website: `https://yourdomain.com`
- Webhook endpoint: `https://yourdomain.com/ebay-webhook`
- Bot runs on Proxmox at `192.168.0.12:8081`
- Your reverse proxy (Nginx/Caddy) forwards `/ebay-webhook` → `192.168.0.12:8081`

**Nginx Configuration Example:**
```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;
    
    # Your existing website config...
    
    # Add webhook endpoint
    location /ebay-webhook {
        proxy_pass http://192.168.0.12:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

**What you need:**
- ✅ Your existing domain name with SSL
- ✅ Add reverse proxy rule
- ✅ Ensure firewall allows your proxy → bot (192.168.0.12:8081)

---

### ✅ Option 2: Cloudflare Tunnel (FREE, No Port Forwarding)

**Perfect if you want to keep things simple and secure.**

**Benefits:**
- ✅ Free static URL
- ✅ Free SSL certificate  
- ✅ No port forwarding needed
- ✅ No exposing internal ports
- ✅ DDoS protection included

**Setup (5 minutes):**

```powershell
# 1. Download cloudflared
Invoke-WebRequest -Uri "https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-windows-amd64.exe" -OutFile "cloudflared.exe"

# 2. Login to Cloudflare (opens browser)
.\cloudflared.exe tunnel login

# 3. Create tunnel
.\cloudflared.exe tunnel create ebay-webhook

# 4. Configure tunnel (creates config.yml)
# Point to your bot at 192.168.0.12:8081

# 5. Run tunnel
.\cloudflared.exe tunnel run ebay-webhook
```

**Result:** You get `https://youruniqueid.cfargotunnel.com` that routes to your bot

**Installation on Proxmox:**
```bash
# SSH into your Proxmox LXC/VM running the bot
apt-get install cloudflared
cloudflared tunnel login
cloudflared tunnel create ebay-webhook
cloudflared tunnel route dns ebay-webhook ebay-webhook.yourdomain.com
# Configure and run as systemd service
```

---

### ✅ Option 3: Direct Port Forwarding + Let's Encrypt

**If you want the bot directly accessible:**

**Requirements:**
- Public IP address (from your ISP)
- Port forwarding on your router
- Dynamic DNS (if your IP changes)
- SSL certificate (Let's Encrypt)

**Steps:**
1. **Router Configuration:**
   - Forward external port 443 → 192.168.0.12:8081
   - Or forward external port 8081 → 192.168.0.12:8081

2. **SSL Certificate:**
   - Use Let's Encrypt with certbot
   - Or use Caddy (auto SSL)

3. **Firewall:**
   - Allow incoming on port 443/8081
   - Proxmox firewall rules

**Less Recommended Because:**
- ❌ Exposes bot directly to internet
- ❌ More security risk
- ❌ Requires SSL certificate management
- ❌ Dynamic IP issues

---

## Recommended Solution for Your Setup

**Based on "hosting my website":**

### Use Your Existing Website! (Option 1)

**Why this is best:**
- ✅ You already have HTTPS working
- ✅ You already have a domain
- ✅ Just add a reverse proxy rule
- ✅ No additional services needed
- ✅ Professional setup

**What you need to tell me:**
1. Your public domain name (e.g., `example.com`)
2. What reverse proxy you're using (Nginx, Caddy, Traefik, Apache?)
3. Where your reverse proxy config is located

**Then I can:**
- Give you the exact config to add
- Test the webhook endpoint
- Configure the bot with your domain

---

## Quick Test: Is Your Setup Already Ready?

**Run this to check if your website is accessible:**

```powershell
# Replace with your actual domain
$domain = "yourdomain.com"
Invoke-WebRequest -Uri "https://$domain" -UseBasicParsing
```

**If this works, you just need to:**
1. Add webhook path to your reverse proxy
2. Use `https://yourdomain.com/ebay-webhook` for eBay webhooks
3. Done! ✅

---

## What Information I Need from You

**To help you configure this, please tell me:**

1. **Your public domain name** (the one hosting your website)
   - Example: `mysite.com` or `subdomain.mysite.com`

2. **Your reverse proxy software**
   - Nginx, Caddy, Traefik, Apache, or something else?

3. **OR** if you prefer Option 2 (Cloudflare Tunnel):
   - Just say "Let's use Cloudflare Tunnel"
   - I'll walk you through it step-by-step

---

## Important Notes

- **192.168.0.12 only works on your local network** - eBay cannot reach it
- **Your website is already accessible publicly** - we just need to route webhooks to the bot
- **This is the final piece** - once webhooks work, you're fully production ready
- **ngrok free tier won't work** - URLs change daily, eBay needs static endpoints

---

## Next Steps

**Choice 1: Use Existing Website**
→ Tell me your domain and reverse proxy software  
→ I'll give you exact config to add
→ 5 minute setup

**Choice 2: Cloudflare Tunnel**  
→ Say "Cloudflare Tunnel"
→ I'll guide you through installation
→ 10 minute setup, completely free

**Once webhooks are configured, we'll:**
1. Get your production eBay credentials
2. Switch environment to PRODUCTION
3. Test with one real listing
4. You're live! 🚀
