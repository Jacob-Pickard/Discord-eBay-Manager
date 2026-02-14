# Getting Production eBay Credentials

## Step-by-Step Guide to Get Production Keys

### 1. Access eBay Developer Portal

**Go to:** https://developer.ebay.com/my/keys

**Sign in with:** Your **production eBay seller account** (not sandbox account)

---

### 2. Switch to Production Environment

**IMPORTANT:** At the top of the page, you'll see an environment selector:
```
[ SANDBOX ] [ PRODUCTION ] <-- Click PRODUCTION
```

**Make sure it shows "Production" highlighted!**

---

### 3. Get Your Production Application Keys

**You should see (or need to create):**

#### Option A: If You Already Have Production Keys

Look for:
```
Application Keys (Production)
├── App ID (Client ID): JacobPic-YourApp-PRD-xxxxxxxxx
├── Cert ID (Client Secret): PRD-xxxxxxxxx-xxxx-xxxx
└── Dev ID: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

**Copy these three values** - they're different from your sandbox keys!

#### Option B: If You Need to Create Production Keys

1. Click **"Create Application Keys"** or **"Get Access"**
2. Fill out the form:
   - Application Title: "Discord eBay Manager Production"
   - Application Type: "Web Application"
3. Submit and wait for approval (usually instant, can take up to 24 hours)
4. Once approved, you'll see your production keys

---

### 4. Get Your Production RuName (Redirect URI)

**This is the most important part!**

**In the same portal:**
1. Look for **"User Tokens"** section
2. Click **"Get a Token from eBay via Your Application"**
3. You'll see **"Auth Accepted URL"** or **"RuName"**

**If you don't have a Production RuName yet:**

1. Click **"Add RuName"**
2. Fill out:
   ```
   Your Company or Name: Jacob_Pickard
   Your App Name: Discord_eBay_Manager  
   Privacy Policy URL: https://yourdomain.com/privacy (or use placeholder)
   RuName: Jacob_Pickard-Discord_eB-Produc-xxxxx
   ```
3. Submit for approval
4. **Production RuNames require eBay approval** - can take 1-3 business days

**Your RuName will look like:**
```
Jacob_Pickard-JacobPic-Discor-prodxxx
```

**Key Difference from Sandbox:**
- Sandbox RuName: `Jacob_Pickard-JacobPic-Discor-mxtojv` (yours)
- Production RuName: `Jacob_Pickard-JacobPic-Discor-xxxxx` (**different!**)

---

### 5. Verify Your Production Scopes

**Make sure your production application has these scopes:**

**Required Scopes:**
- ✅ `https://api.ebay.com/oauth/api_scope`
- ✅ `https://api.ebay.com/oauth/api_scope/sell.account`
- ✅ `https://api.ebay.com/oauth/api_scope/sell.fulfillment`
- ✅ `https://api.ebay.com/oauth/api_scope/sell.inventory`
- ✅ `https://api.ebay.com/oauth/api_scope/sell.marketing`
- ✅ `https://api.ebay.com/oauth/api_scope/sell.finances`
- ✅ `https://api.ebay.com/oauth/api_scope/commerce.notification.subscription`

**To check/add scopes:**
1. In developer portal, go to your production app
2. Click "Edit Application"
3. Scroll to "OAuth Scopes"
4. Select all required scopes
5. Save

---

### 6. Check Your eBay Managed Payments Status

**Go to:** https://www.ebay.com/sh/fin

**Look for:**
- **"Payments"** tab in Seller Hub
- Should show your balance, payouts, etc.
- If you see financial info → You have Managed Payments ✅
- If it redirects you elsewhere → You may not be enrolled yet ❌

**If not enrolled:**
- Some finance APIs won't work (will show sample data)
- You can still use all other features (offers, orders, webhooks)
- Contact eBay to enroll in Managed Payments

---

### 7. Verify You Can Create Listings

**Go to:** https://www.ebay.com/sl/list

**Try creating a test listing:**
- Category: Choose any (e.g., "Toys & Hobbies")
- Condition: New/Used
- Price: $10 (test price)
- **Enable "Best Offer"** ← Important for testing `/accept-offer`
- Quantity: 1

**Don't publish yet** - just verify you *can* create listings

---

## Your Production Credentials Checklist

**Once you have all these, you're ready:**

```
✅ Production App ID (Client ID): ____________________
✅ Production Cert ID (Client Secret): _______________
✅ Production Dev ID: _______________________________
✅ Production RuName: _______________________________
✅ Verified Managed Payments enrolled
✅ Verified can create listings with Best Offer
```

---

## What If Production RuName Is Pending?

**eBay Production RuNames need approval (1-3 business days)**

**Options while waiting:**

**Option 1: Wait for approval (Recommended)**
- Continue testing in sandbox
- Check approval status daily at developer.ebay.com
- Once approved, switch to production

**Option 2: Request expedited approval**
- Contact eBay Developer Support
- Explain you need it for production deployment
- They sometimes expedite for active sellers

**Option 3: Use existing production app**
- If you have another eBay app with production access
- Use those credentials instead

---

## Common Issues

### "I don't see Production keys"
**Solution:** Your account may need verification
- Complete eBay Developer Agreement
- Verify email/phone
- May need to wait 24 hours after signing up

### "Production RuName stuck in pending"
**Solution:** Contact eBay Developer Support
- https://developer.ebay.com/support
- Usually approved within 1-3 business days
- Expedite requests sometimes granted

### "I can't find User Tokens section"
**Solution:** 
- Make sure you're on "Production" tab (not Sandbox)
- Look under "Application Keys (Production)"
- Click "Get OAuth Application Credentials"

---

## Once You Have Everything

**Create this file: `.env.production`**

```env
# Discord Configuration
DISCORD_BOT_TOKEN=<YOUR_DISCORD_BOT_TOKEN>

# eBay PRODUCTION API Configuration
EBAY_APP_ID=<YOUR_PRODUCTION_APP_ID>
EBAY_CERT_ID=<YOUR_PRODUCTION_CERT_ID>
EBAY_DEV_ID=<YOUR_PRODUCTION_DEV_ID>
EBAY_REDIRECT_URI=<YOUR_PRODUCTION_RUNAME>

# NO TOKENS - These will be generated after authorization
# Leave these commented out or empty
# EBAY_ACCESS_TOKEN=
# EBAY_REFRESH_TOKEN=

# Environment - PRODUCTION
EBAY_ENVIRONMENT=PRODUCTION

# Webhook Configuration (update after webhook setup)
WEBHOOK_PORT=8081
WEBHOOK_VERIFY_TOKEN=my_secure_verify_token_12345
NOTIFICATION_CHANNEL_ID=<your_discord_channel_id>
```

---

## Ready to Deploy?

**Once you have:**
✅ All production credentials  
✅ Webhook endpoint configured (see WEBHOOK_PRODUCTION_SETUP.md)

**Then we'll:**
1. Back up your sandbox config
2. Switch to production 
3. Authorize with production account
4. Test with one real listing
5. Go live! 🚀

**Tell me when you have the credentials and webhook sorted!**
