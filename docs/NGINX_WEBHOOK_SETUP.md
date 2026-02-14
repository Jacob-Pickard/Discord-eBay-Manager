# Nginx Configuration for jacob.it.com Webhook

## Your Setup
- **Domain:** jacob.it.com
- **Reverse Proxy:** Nginx
- **Bot Server:** 192.168.0.12:8081
- **Webhook Endpoints:** 
  - https://jacob.it.com/webhook/ebay/notification (main endpoint)
  - https://jacob.it.com/webhook/ebay/challenge (verification)
  - https://jacob.it.com/webhook/health (health check)

---

## Nginx Configuration to Add

### Location: `/etc/nginx/sites-available/jacob.it.com` (or your config file)

Add this location block to your existing server configuration:

```nginx
server {
    listen 443 ssl http2;
    server_name jacob.it.com;
    
    # Your existing SSL configuration
    ssl_certificate /path/to/your/cert.pem;
    ssl_certificate_key /path/to/your/key.pem;
    
    # ... your existing website configuration ...
    
    # eBay Bot Webhook Endpoints - ADD THIS
    location /webhook/ {
        proxy_pass http://192.168.0.12:8081;
        proxy_http_version 1.1;
        
        # Forward original headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support (if needed)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # Allow large payloads from eBay
        client_max_body_size 10M;
    }
    
    # ... rest of your existing configuration ...
}
```

---

## Step-by-Step Setup

### 1. Find Your Nginx Configuration File

**On your Proxmox server, SSH in and run:**

```bash
# Find your jacob.it.com config
ls /etc/nginx/sites-available/
# OR
ls /etc/nginx/conf.d/

# Common locations:
# - /etc/nginx/sites-available/jacob.it.com
# - /etc/nginx/sites-available/default
# - /etc/nginx/conf.d/jacob.it.com.conf
```

### 2. Edit the Configuration

```bash
# Edit your site config (use your actual file path)
sudo nano /etc/nginx/sites-available/jacob.it.com

# Or if using vi/vim
sudo vim /etc/nginx/sites-available/jacob.it.com
```

**Add the `/ebay-webhook` location block** shown above inside your `server { }` block

### 3. Test Nginx Configuration

```bash
# Test for syntax errors
sudo nginx -t

# Should show:
# nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### 4. Reload Nginx

```bash
# Reload to apply changes (no downtime)
sudo systemctl reload nginx

# Or restart if needed
sudo systemctl restart nginx
```

### 5. Test from Your Windows Machine

**From PowerShell on your Windows machine:**

```powershell
# Test if the health endpoint is accessible
Invoke-WebRequest -Uri "https://jacob.it.com/webhook/health" -Method GET

# You should get a response from your bot showing "OK" or health status
```

---

## Firewall Configuration

### Ensure Proxmox Firewall Allows Traffic

**On your Proxmox host/VM:**

```bash
# Check if firewall is blocking
sudo ufw status

# If firewall is active, allow connections from your reverse proxy
sudo ufw allow from <proxy_server_ip> to any port 8081

# Or if reverse proxy is on same machine
sudo ufw allow 8081/tcp
```

---

## Update Your Bot Configuration

**Once Nginx is configured, update your `.env.production.template`:**

The webhook endpoint will be: **`https://jacob.it.com/ebay-webhook`**
webhook/ebay/notification`**

Your production config already has
```env
# In .env.production.template (when you create it)
WEBHOOK_PORT=8081
NOTIFICATION_CHANNEL_ID=your_discord_channel_id
```

**The bot needs to know to listen on port 8081** (it already does this by default)

---

## Testing the Complete Flow

### Test 1: Direct Bot Connection (from Proxmox server)

```bash
# SSH into the server running the bot
curl http://192.168.0.12:8081/health

# Should respond with bot status
```

### Test 2: Through Nginx (from Proxmox server or Windows)

```bash
curl https://jacob.it.com/ebay-webhook
```webhook/health
```

```powershell
# From Windows PowerShell
Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health

### Test 3: From External Internet

**From your phone/another network:**
- Visit: `https://jacob.it.com/webhook/health`
- Should get a response showing bot health status

---

## Common Issues & Solutions

### Issue: "502 Bad Gateway"
**Cause:** Nginx can't reach the bot at 192.168.0.12:8081

**Fix:**
```bash
# 1. Check if bot is running
curl http://192.168.0.12:8081

# 2. Check firewall
sudo ufw status

# 3. Check bot logs
tail -f /path/to/bot/bot-error.log
```

### Issue: "Connection Refused"
**Cause:** Bot isn't listening on port 8081

**Fix:**
```bash
# Check if bot is running and listening
netstat -tlnp | grep 8081

# Or
ss -tlnp | grep 8081

# Should show the bot process
```

### Issue: "SSL Certificate Error"
**Cause:** Your existing SSL cert doesn't cover this endpoint

**Fix:**
- The /ebay-webhook path uses your existing SSL cert
- Should work automatically if your site already has HTTPS
- If using Let's Encrypt, cert should cover all paths

### Issue: "404 Not Found"
**Cause:** Nginx config not loaded

**Fix:**
```bash
# Test config
sudo nginx -t

# Reload
sudo systemctl reload nginx

# Check Nginx error log
sudo tail -f /var/log/nginx/error.log
```

---

## Verification Checklist

Before proceeding to production:

- [ ] Nginx configuration added to jacob.it.com config file
- [ ] `sudo nginx -t` passes with no errors
- [ ] Nginx reloaded: `sudo systemctl reload nginx`
- [ ] Bot is running on 192.168.0.12:8081
- [ ] `https://jacob.it.com/webhook/health` is accessible from Windows
- [ ] Firewall allows traffic between Nginx and bot
- [ ] SSL certificate works (should - same as your site)

---

## Next: Update eBay Subscription

**Once this is working, you'll configure eBay webhooks with:**

**Webhook Destination URL:** `https://jacob.it.com/webhook/ebay/notification`

**In Discord bot commands:**
```
/webhook-subscribe url:https://jacob.it.com/webhook/ebay/notification
```

You'll use this URL when subscribing to eBay notifications.

---

## Quick Copy-Paste for Your Server

**SSH into your Proxmox server and run:**

```bash
# 1. Edit Nginx config (adjust path if needed)
sudo nano /etc/nginx/sites-available/jacob.it.com

# 2. Add the location block provided above

# 3. Test configuration
sudo nginx -t

# 4. Reload Nginx
sudo systemctl reload nginx

# 5. Test from server
curl https://jacob.it.com/ebay-webhook
```

**Then from Windows PowerShell:**

```powershell
# Test external access
Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health"
```

---

## Ready?

Once you've added the Nginx configuration and tested it, tell me:
- ✅ "Nginx configured and tested"

Then we'll move on to getting your production eBay credentials!
