# 🚀 Quick Testing Start Guide

Now that your environment is set up, follow these steps to test all features:

## ✅ Your Current Status
- **Environment:** SANDBOX ✓
- **Discord Bot:** Configured ✓
- **eBay API:** Configured ✓  
- **OAuth Tokens:** Already generated ✓
- **Webhooks:** Ready to test ✓

---

## 📝 Step-by-Step Testing Process

### 1. Start the Bot (5 minutes)

```powershell
# From the project directory:
.\ebaymanager-bot.exe
```

**Expected output:**
```
Connecting to Discord...
Connected as: YourBotName#1234 (ID: ...)
Starting webhook server on :8080
eBay Manager Bot is now running. Press CTRL+C to exit.
```

✅ **Verify:** Bot shows as online in Discord server

---

### 2. Test Basic Commands (10 minutes)

Open Discord and test these commands in sequence:

#### A. Check Connection
```
/ebay-status
```
**Should show:** ✅ Connected to eBay API (SANDBOX mode)

#### B. View Orders
```
/get-orders
```
**Result:** Shows recent orders OR "No recent orders found"

#### C. View Offers
```
/get-offers
```
**Result:** Shows pending offers OR "No pending offers found"

#### D. View Balance
```
/get-balance
```
**Result:** Shows account balance (may be $0 in sandbox)

#### E. View Messages
```
/get-messages
```
**Result:** Shows buyer messages or "No messages"

---

### 3. Create Test Data in Sandbox (20 minutes)

You likely need test data. Follow these steps:

#### A. Create a Test Listing
1. Go to: https://sandbox.ebay.com
2. Sign in with your seller account
3. Click **"Sell"** → **"List an item"**
4. Fill out the form:
   - **Title:** "Test Gaming Headset - RGB LED"
   - **Category:** Electronics > Video Games
   - **Price:** $79.99
   - **Quantity:** 5
   - **✅ IMPORTANT:** Enable "Best Offer"
     - Auto-accept: $75.00
     - Auto-decline: $40.00
5. Click **"List Item"**

#### B. Create a Test Buyer Account
1. Open an **incognito/private browser window**
2. Go to: https://sandbox.ebay.com
3. Click **"Register"**
4. Create account:
   - Email: `testbuyer2026@test.com` (can be fake)
   - Username: `TestBuyer2026`
   - Password: Any secure password
5. Complete registration

#### C. Make Test Offers (as Buyer)
In the incognito window:
1. Search for your listing: "Test Gaming Headset"
2. Click on your item
3. Click **"Make Offer"**
4. Make these 3 offers one at a time:
   - **Offer 1:** $45.00 (low - will decline)
   - **Offer 2:** $65.00 (medium - will counter)
   - **Offer 3:** $77.00 (high - will accept)

---

### 4. Test Offer Management (15 minutes)

Back in Discord:

#### A. View the Offers
```
/get-offers
```
**Should show:** All 3 offers with IDs

#### B. Decline the Low Offer
```
/decline-offer offer-id:OFFER_1_ID
```
**Expected:** ✅ Offer declined successfully!

#### C. Counter the Medium Offer
```
/counter-offer offer-id:OFFER_2_ID price:70.00
```
**Expected:** ✅ Counter offer sent successfully!

**Verify:** In buyer's browser, go to "My eBay" → "Offers" to see the counter

#### D. Accept the High Offer
```
/accept-offer offer-id:OFFER_3_ID
```
**Expected:** ✅ Offer accepted successfully! Order created.

#### E. Verify Order Created
```
/get-orders
```
**Should show:** New order from the accepted offer

---

### 5. Test Webhooks (20 minutes)

Webhooks require a public URL. Use ngrok:

#### A. Install ngrok (if not installed)
Download from: https://ngrok.com/download

Or with Chocolatey:
```powershell
choco install ngrok
```

#### B. Start ngrok
In a **new terminal window**:
```powershell
ngrok http 8080
```

**Copy the HTTPS URL** (e.g., `https://abc123.ngrok.io`)

#### C. Subscribe to Webhooks
In Discord:
```
/webhook-subscribe url:https://YOUR-NGROK-URL.ngrok.io
```
**Expected:** ✅ Successfully subscribed to webhooks!

#### D. List Subscriptions
```
/webhook-list
```
**Should show:** Active webhook subscriptions

#### E. Trigger Test Events
Create events to trigger webhooks:

**New Offer:**
- As buyer, make another offer on your listing
- **Expected:** Discord notification appears in your notification channel

**New Order:**
- As buyer, purchase an item via "Buy It Now"
- **Expected:** Discord notification for new order

**New Message:**
- As buyer, ask a question on the item
- **Expected:** Discord notification for new message

---

### 6. Test Edge Cases (10 minutes)

#### A. Invalid Offer ID
```
/accept-offer offer-id:INVALID123
```
**Expected:** Error message (not crash)

#### B. Negative Price
```
/counter-offer offer-id:SOME_ID price:-50
```
**Expected:** Validation error

#### C. Already Handled Offer
Try accepting an offer you already declined
**Expected:** Appropriate error message

---

### 7. Test Token Refresh (Wait 2-3 minutes)

**Watch the console output** while bot is running

**Expected to see every 90 minutes:**
```
🔄 Refreshing eBay access token...
✅ Token refreshed successfully
```

**If you see this:** Token auto-refresh is working! ✅

---

## ✅ Testing Checklist Summary

Use this quick checklist while testing:

- [ ] Bot starts without errors
- [ ] `/ebay-status` shows connected
- [ ] `/get-orders` works
- [ ] `/get-offers` works
- [ ] `/get-balance` works
- [ ] `/get-messages` works
- [ ] Created test listing with Best Offer enabled
- [ ] Created test buyer account
- [ ] Made 3 test offers
- [ ] `/decline-offer` works correctly
- [ ] `/counter-offer` works correctly
- [ ] `/accept-offer` creates order
- [ ] Webhooks subscribed via `/webhook-subscribe`
- [ ] Webhook notifications received in Discord
- [ ] Token auto-refresh working (check logs)
- [ ] Error handling works (invalid inputs)
- [ ] No crashes or critical bugs

---

## 🐛 If You Find Issues

Document them in this format:

```
**Bug:** Brief description
**Steps to reproduce:** 
1. Step 1
2. Step 2
**Expected:** What should happen
**Actual:** What actually happened
**Severity:** Critical / High / Medium / Low
```

---

## ✅ Testing Complete!

Once you've checked all boxes above:

1. ✅ All commands working
2. ✅ Offers can be managed
3. ✅ Webhooks receiving notifications
4. ✅ No critical bugs

**Your eBay Manager Bot is ready for production use!** Test it with real listings and offers to ensure everything works as expected.

---

## 📞 Need Help?

**Common Issues:**

**Bot won't start:**
- Check `.env` file exists
- Run `go run tools/check_config.go`
- Verify Discord token is valid

**Commands not responding:**
- Check bot is online in Discord
- Verify bot has proper permissions
- Check console for errors

**No offers/orders showing:**
- Verify you're in SANDBOX mode
- Create test data following Step 3 above
- Check eBay sandbox account has listings

**Webhooks not working:**
- Ensure ngrok is running
- URL must be HTTPS
- Check WEBHOOK_VERIFY_TOKEN is set
- Verify NOTIFICATION_CHANNEL_ID is correct

---

## 🎯 Time Estimate

**Total testing time:** ~80 minutes
- Initial setup: 5 min
- Basic commands: 10 min
- Create test data: 20 min
- Offer management: 15 min
- Webhook testing: 20 min
- Edge cases: 10 min

**Pro tip:** Do it in one session to maintain context!
