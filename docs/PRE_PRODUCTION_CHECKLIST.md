# Pre-Production Deployment Checklist

## CRITICAL: Complete Before Switching to Production

### 1. Production Credentials Verification

**Do you have production credentials from eBay Developer Portal?**

You need **SEPARATE** production credentials (different from sandbox):

- [ ] **Production App ID (Client ID)** - Different from sandbox
- [ ] **Production Cert ID (Client Secret)** - Different from sandbox  
- [ ] **Production Dev ID** - May be same as sandbox
- [ ] **Production RuName (Redirect URI)** - MUST be different from sandbox

**Where to get these:**
1. Go to https://developer.ebay.com/my/keys
2. Switch to **"Production"** environment at top
3. Create a Production Application Key Set if you haven't
4. Get the Production RuName from the User Tokens section

---

### 2. Production Account Requirements

**Your eBay seller account MUST have:**

- [ ] **Active seller account** - Not a test/sandbox account
- [ ] **Managed Payments enrolled** - Required for /get-balance, /get-payouts
- [ ] **Selling history** - At least one listing
- [ ] **OAuth grant flow permissions** - Approved in eBay Dev Portal

**Verify at:**
- Seller Hub: https://www.ebay.com/sh/ovw
- Payments: https://www.ebay.com/sh/fin

---

### 3. Webhook Endpoint Requirements

**For production webhooks, you need:**

- [ ] **Public HTTPS endpoint** - No HTTP, no localhost
- [ ] **Valid SSL certificate** - Not self-signed
- [ ] **Static public IP or domain** - ngrok free tier changes URL

**Options:**
1. **VPS/Cloud Server** (Recommended)
   - AWS EC2, DigitalOcean, Linode
   - Static IP with SSL
   - Runs 24/7

2. **Cloudflare Tunnel** (Free alternative)
   - Free static URL
   - SSL included
   - Can run from home

3. **ngrok Pro** ($10/month)
   - Static domain
   - Reliable for testing

---

### 4. Backup Current Configuration

```powershell
# Run this to backup your current working sandbox configuration
Copy-Item .env .env.sandbox.backup
```

This preserves your working sandbox setup in case you need to revert.

---

### 5. Production Environment File

Create a **NEW** `.env.production` file with your production credentials:

```env
# Discord Configuration (same as sandbox)
DISCORD_BOT_TOKEN=your_discord_token

# eBay PRODUCTION API Configuration
EBAY_APP_ID=YOUR_PRODUCTION_APP_ID
EBAY_CERT_ID=YOUR_PRODUCTION_CERT_ID
EBAY_DEV_ID=YOUR_PRODUCTION_DEV_ID
EBAY_REDIRECT_URI=YOUR_PRODUCTION_RUNAME

# NO TOKENS YET - Will be generated after authorization
# EBAY_ACCESS_TOKEN=
# EBAY_REFRESH_TOKEN=

# Environment - SET TO PRODUCTION
EBAY_ENVIRONMENT=PRODUCTION

# Webhook Configuration
WEBHOOK_PORT=8081
WEBHOOK_VERIFY_TOKEN=your_secure_verify_token
NOTIFICATION_CHANNEL_ID=your_discord_channel_id
```

---

### 6. Critical Differences: Sandbox vs Production

| Aspect | Sandbox | Production |
|--------|---------|------------|
| **App Credentials** | Test keys | Separate production keys |
| **RuName** | Sandbox RuName | Production RuName |
| **eBay Account** | Sandbox test account | Real seller account |
| **Listings** | Test listings | Real listings |
| **Money** | Fake money | Real money |
| **Buyers** | Test buyers | Real buyers |
| **Webhooks** | Can use ngrok | Need stable endpoint |
| **API Limits** | Generous | Strict rate limits |
| **Errors** | Test friendly | Production critical |

---

### 7. Deployment Steps (When Ready)

**Step 1: Backup & Prepare**
```powershell
# Backup sandbox config
Copy-Item .env .env.sandbox.backup

# Copy production config to .env
Copy-Item .env.production .env
```

**Step 2: Stop Bot**
```powershell
# Stop the running bot (if running)
# Press Ctrl+C in the terminal or close the bot window
```

**Step 3: Start with Production Config**
```powershell
# Start bot with production environment
.\ebaymanager-bot.exe
```

**Step 4: Authorize Production Account**
```
1. In Discord: /ebay-status
   - Should show "Not authorized" or "Token expired"
   
2. In Discord: /ebay-authorize
   - Bot returns PRODUCTION authorization URL
   
3. Click link, sign in with PRODUCTION eBay seller account
   - NOT sandbox account
   
4. Authorize app, get code
   
5. In Discord: /ebay-code code:<the_code_from_ebay>
   - Bot saves tokens to .env
   
6. In Discord: /ebay-status
   - Should show "Authorized" with production account
```

**Step 5: Verify Production Webhooks**
```
1. In Discord: /webhook-subscribe
   - Subscribe to production notifications
   
2. Wait for confirmation message
   
3. In Discord: /webhook-list
   - Should show active subscriptions
```

**Step 6: Test with Real Listing**
```
1. Create a test listing on real eBay
   - Enable "Best Offer"
   - List at $100, accept offers of $75+
   
2. Use a friend/alt account to make a test offer
   
3. In Discord: /get-offers
   - Should show the real offer
   
4. Test: /accept-offer or /counter-offer
   - This creates a REAL transaction
   
5. Verify webhook notification arrives in Discord
```

---

### 8. Post-Deployment Monitoring

**Watch for these in first 24 hours:**
- [ ] OAuth token refresh works automatically
- [ ] Webhook notifications arrive consistently  
- [ ] All commands respond without errors
- [ ] No "403 Forbidden" or "401 Unauthorized" errors
- [ ] Real financial data shows correctly (/get-balance)

**Check logs:**
```powershell
Get-Content bot-error.log -Tail 50
```

---

### 9. Emergency Rollback Plan

**If something goes wrong:**

```powershell
# Stop production bot
# Kill the bot process

# Restore sandbox configuration
Copy-Item .env.sandbox.backup .env

# Restart with sandbox
.\ebaymanager-bot.exe

# Your sandbox setup is intact
```

---

### 10. Important Production Warnings

⚠️ **REAL MONEY**: All transactions are real in production
⚠️ **REAL BUYERS**: Commands affect real customer orders
⚠️ **NO UNDO**: Accepting/declining offers cannot be undone
⚠️ **API LIMITS**: Production has stricter rate limits
⚠️ **ACCOUNT RISK**: API violations can suspend your account

**Best Practice:**
- Test all offer commands on ONE test listing first
- Use low-value items for initial testing ($5-10)
- Monitor closely for the first week
- Keep sandbox environment as testing backup

---

## FINAL VERIFICATION BEFORE PROCEEDING

**I have:**
- [ ] Production App ID, Cert ID, Dev ID, RuName from eBay
- [ ] Real eBay seller account with Managed Payments
- [ ] Public HTTPS webhook endpoint (not ngrok free tier)
- [ ] Backed up current .env file
- [ ] Created .env.production with production credentials
- [ ] Tested all features in sandbox successfully
- [ ] Understood that production affects real money/orders
- [ ] Read emergency rollback procedure

**If you checked ALL boxes above, you are ready for production deployment.**

**If ANY box is unchecked, DO NOT proceed yet.**

---

## Questions to Answer Before Proceeding

1. **Do you have production credentials?** (separate from sandbox)
2. **Is your eBay account enrolled in Managed Payments?**
3. **Do you have a stable public HTTPS webhook endpoint?** (not ngrok free)
4. **Have you successfully tested offer management in sandbox?**
5. **Do you understand production affects real money and orders?**

**All answers must be YES to proceed safely.**
