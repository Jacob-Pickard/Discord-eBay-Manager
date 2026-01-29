# eBay Sandbox Testing Guide

## 🧪 Creating Real Test Data in eBay Sandbox

Now that you have offer management commands, let's create real test orders and offers in your eBay Sandbox environment!

## 📋 Prerequisites

- eBay Developer Account (you already have this)
- Sandbox Seller Account (your main sandbox account)
- Sandbox Buyer Account (need to create one)

## Step 1: Create a Sandbox Buyer Account

1. **Go to eBay Sandbox**: https://sandbox.ebay.com
2. **Sign out** if logged in as seller
3. **Register new account** as a buyer:
   - Click "Register" 
   - Use a different email (can be fake: buyer123@test.com)
   - Username: TestBuyer2026
   - Complete registration

## Step 2: Create a Test Listing (As Seller)

1. **Log into Sandbox as Seller**: https://sandbox.ebay.com
   - Use your developer account credentials

2. **Create a Listing**:
   - Click "Sell" → "List an item"
   - Fill in details:
     - **Title**: "Test Gaming Headset - RGB Edition"
     - **Category**: Electronics > Video Games > Accessories
     - **Condition**: New
     - **Price**: $79.99
     - **Quantity**: 5
     - **Best Offer**: ✅ **ENABLE THIS** (crucial for testing offers!)
     - **Auto-accept price**: $75.00 (optional)
     - **Auto-decline price**: $50.00 (optional)

3. **Set Shipping**:
   - Flat rate: $5.99
   - Ships to: United States

4. **List the item**

## Step 3: Make a Test Offer (As Buyer)

1. **Log out and log back in as Buyer**: https://sandbox.ebay.com
   - Sign in with your TestBuyer2026 account

2. **Search for your listing**:
   - Search: "Test Gaming Headset"
   - Find your listing

3. **Make an offer**:
   - Click "Make Offer" button
   - Enter offer amount: **$65.00** (below your price)
   - Add message: "Very interested! Can you accept $65?"
   - Submit offer

## Step 4: View Offers in Discord

1. Go to your Discord server
2. Run: `/get-offers`
3. You should now see YOUR REAL OFFER:
   ```
   **Offer from TestBuyer2026**
   • Item: Test Gaming Headset - RGB Edition
   • Listed Price: $79.99 USD
   • Offer Amount: $65.00 USD
   • Status: PENDING
   ```

## Step 5: Test Offer Management

Now you can test the offer commands:

### Accept the Offer:
```
/accept-offer offer-id:<the-offer-id-from-get-offers>
```

### Counter the Offer:
```
/counter-offer offer-id:<offer-id> price:70.00
```

### Decline the Offer:
```
/decline-offer offer-id:<offer-id>
```

## Step 6: Create a Test Order (Buy It Now)

1. **As Buyer**, find your listing again
2. **Click "Buy It Now"** (not Make Offer)
3. **Complete checkout**:
   - Use Sandbox PayPal: https://www.sandbox.paypal.com
   - Credentials: Any test PayPal account
   - Complete payment

4. **Check Orders in Discord**:
   ```
   /get-orders
   ```
   You should see your real test order!

## 🎯 Quick Test Scenario

**Complete Test Flow:**

1. **Create Listing** (Seller):
   - Title: "Vintage Camera Collection"
   - Price: $500
   - Enable Best Offer

2. **Make 3 Offers** (Buyer):
   - Offer 1: $425 ("Interested! Can you do $425?")
   - Switch to 2nd buyer account
   - Offer 2: $450 ("Cash ready!")
   - Switch to 3rd buyer account  
   - Offer 3: $475 ("Great deal!")

3. **Manage Offers** (Discord):
   ```
   /get-offers
   /counter-offer offer-id:123 price:460
   /accept-offer offer-id:456
   /decline-offer offer-id:789
   ```

4. **Create Order** (Buyer buys):
   - Buy It Now → Complete Payment

5. **Check Order** (Discord):
   ```
   /get-orders
   ```

## 🔧 Useful Sandbox Tools

### Sandbox User Accounts
Create multiple test users here:
https://developer.ebay.com/sandbox/user-accounts

**Create:**
- 1 Seller account (you)
- 3+ Buyer accounts (for testing)

### Sandbox PayPal
Use test PayPal accounts:
- Buyer: buyer@test.com / password
- Seller: Your sandbox PayPal

### API Test Tool
Test API calls directly:
https://developer.ebay.com/api-explorer/

## 📊 Expected Results

After setup, your bot should show:

**`/get-offers`**:
```
💰 Pending Buyer Offers

**Offer from TestBuyer2026**
• Item: Test Gaming Headset - RGB Edition
• Listed Price: $79.99 USD
• Offer Amount: $65.00 USD (18% off)
• Message: "Very interested! Can you accept $65?"
• Status: PENDING
• Offer ID: 12345-67890

**Offer from CoolCollector**
• Item: Vintage Camera Collection
• Listed Price: $500.00 USD
• Offer Amount: $450.00 USD (10% off)
• Status: PENDING
• Offer ID: 98765-43210
```

**`/get-orders`**:
```
📋 Recent Orders

**Order #19-12345-67890**
• Buyer: TestBuyer2026
• Item: Test Gaming Headset - RGB Edition
• Total: $85.98 USD (includes shipping)
• Status: AWAITING_SHIPMENT
• Date: Jan 28, 2026
```

## 🎬 Video Tutorial

eBay provides a sandbox tutorial:
https://developer.ebay.com/DevZone/sandboxuser/default.aspx

## ❓ Troubleshooting

**No offers appearing?**
- Make sure "Best Offer" is enabled on listing
- Wait a few minutes for API sync
- Check offer wasn't auto-declined
- Verify API token has correct scopes

**Can't complete purchase?**
- Use sandbox PayPal: https://www.sandbox.paypal.com
- Don't use real PayPal in sandbox
- Create test PayPal if needed

**API errors?**
- Check `/ebay-status` in Discord
- Verify token hasn't expired (tokens last 2 hours)
- Re-authorize with `/ebay-authorize` if needed

## 🚀 Next Steps

Once you have real data:
1. Test all offer management commands
2. Set up webhook notifications for real-time alerts
3. Move to production when ready
4. Add automation (auto-accept offers over X amount)

## 📝 Notes

- Sandbox data resets periodically
- Sandbox transactions are fake (no real money)
- Production requires different credentials
- Always test thoroughly before going live!
