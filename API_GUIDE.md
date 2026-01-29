# eBay Manager Bot - API Setup Guide

## What's Been Implemented

✅ **Completed Features:**
- OAuth 2.0 authentication flow
- Token management with auto-refresh
- Get Orders API (Fulfillment API)
- Get Offers API (Inventory API)  
- Create Listings API (Inventory API)
- Respond to Offers API (Negotiation API)
- Discord bot with slash commands
- Deferred responses for long API calls

## API Architecture

### OAuth Flow

The bot uses eBay's OAuth 2.0 for authentication. Three types of tokens:

1. **Application Token**: For public APIs (no user context)
2. **User Token**: For user-specific operations (requires authorization)
3. **Refresh Token**: To get new access tokens when they expire

### Implemented eBay APIs

#### 1. **Fulfillment API** - Orders Management
- **Endpoint**: `/sell/fulfillment/v1/order`
- **Purpose**: Retrieve order details, shipping info
- **Usage**: `ebayClient.GetOrders(limit)`
- **Returns**: List of orders with buyer info, prices, shipping addresses

#### 2. **Inventory API** - Listings & Offers
- **Endpoints**: 
  - `/sell/inventory/v1/inventory_item/{sku}` - Create inventory
  - `/sell/inventory/v1/offer` - Create/manage offers
- **Purpose**: Create and manage listings
- **Usage**: `ebayClient.CreateListing(listing)`
- **Process**: Creates inventory item → Creates offer → Publishes listing

#### 3. **Negotiation API** - Offer Responses
- **Endpoint**: `/sell/negotiation/v1/offer/{offerId}/respond`
- **Purpose**: Accept, decline, or counter buyer offers
- **Usage**: `ebayClient.RespondToOffer(offerID, "ACCEPT", 0)`
- **Actions**: ACCEPT, DECLINE, COUNTER

## Getting Your eBay Credentials

### Step 1: Create eBay Developer Account

1. Go to https://developer.ebay.com/
2. Sign in with your eBay seller account
3. Accept the developer terms

### Step 2: Create Application Keys

1. Navigate to "My Account" → "Application Keys"
2. Click "Create Application Key Set"
3. Fill in application details:
   - **Application Title**: "eBay Manager Bot"
   - **Select Sandbox Keys** for testing first
4. Save these credentials:
   - **App ID (Client ID)**: Used for OAuth
   - **Cert ID (Client Secret)**: Used for OAuth
   - **Dev ID**: Required for some APIs

### Step 3: Configure OAuth Redirect URI

1. In your application settings, set **Redirect URI**
2. For local testing: `http://localhost:3000/callback`
3. For production: Your actual callback URL
4. This must match exactly in your `.env` file

### Step 4: Generate User Access Token

**Option A: Using eBay's Token Generator (Quickest)**

1. Go to https://developer.ebay.com/my/auth/?env=sandbox&index=0
2. Select these scopes:
   - `https://api.ebay.com/oauth/api_scope`
   - `https://api.ebay.com/oauth/api_scope/sell.inventory`
   - `https://api.ebay.com/oauth/api_scope/sell.fulfillment`
   - `https://api.ebay.com/oauth/api_scope/sell.account`
   - `https://api.ebay.com/oauth/api_scope/sell.marketing`
3. Click "Get a Token from eBay via Your Application"
4. Sign in and authorize
5. Copy the **User Token** (access token) - valid for 2 hours
6. Copy the **Refresh Token** - valid for 18 months

**Option B: Implement Full OAuth Flow**

The bot has `GetUserAuthorizationURL()` and `ExchangeCodeForToken()` methods:

```go
// Generate authorization URL
authURL := ebayClient.GetUserAuthorizationURL("random-state-string")
// Send user to this URL

// After user authorizes, exchange code for tokens
token, err := ebayClient.ExchangeCodeForToken(authorizationCode)
```

### Step 5: Set Up Environment Variables

Create `.env` file:

```env
# Discord Bot Token
DISCORD_BOT_TOKEN=your_discord_bot_token_here

# eBay API Credentials
EBAY_APP_ID=your_app_id_here
EBAY_CERT_ID=your_cert_id_here
EBAY_DEV_ID=your_dev_id_here
EBAY_REDIRECT_URI=http://localhost:3000/callback

# eBay OAuth Tokens
EBAY_ACCESS_TOKEN=v^1.1#i^1#...your_full_token...
EBAY_REFRESH_TOKEN=v^1.1#i^1#...your_full_refresh_token...

# Environment
EBAY_ENVIRONMENT=SANDBOX
```

## API Implementation Details

### Token Management

The bot includes a `TokenManager` that:
- Automatically refreshes tokens before expiry
- Handles token storage and updates
- Provides `GetValidToken()` for API calls

Example usage:
```go
tokenMgr := ebay.NewTokenManager(client, accessToken, refreshToken, expiresIn)
validToken, err := tokenMgr.GetValidToken()
```

### Error Handling

All API calls return proper errors:
- **401 Unauthorized**: Token expired or invalid
- **400 Bad Request**: Invalid parameters
- **429 Too Many Requests**: Rate limit exceeded
- **500 Server Error**: eBay API issues

### Rate Limiting

eBay APIs have rate limits:
- **Sandbox**: 5,000 calls/day
- **Production**: Varies by API (10,000-5,000,000/day)

Consider implementing rate limiting logic if you expect high volume.

## Discord Commands Reference

### `/ebay-status`
Checks eBay API connection status and token validity.

### `/list-item`
Creates a new eBay listing.
- **Parameters**: title, price, description
- **Process**: 
  1. Creates inventory item with SKU
  2. Creates offer with pricing
  3. Publishes to eBay marketplace

### `/get-orders`
Fetches recent orders (last 10).
- Shows order ID, buyer, total, status
- Displays creation date

### `/get-offers`
Fetches pending offers.
- Shows offer details
- Displays buyer offer vs. list price

## What You Need to Implement Next

### 1. **Fulfillment Policies**
Before creating listings, set up:
- Shipping policies
- Payment policies  
- Return policies

Create these in eBay Seller Hub, then get their IDs via API:
```
GET /sell/account/v1/fulfillment_policy
GET /sell/account/v1/payment_policy
GET /sell/account/v1/return_policy
```

Update `CreateListing()` to use real policy IDs instead of "default".

### 2. **Category Selection**
Currently hardcoded to category "1". Get proper category IDs:
```
GET /commerce/taxonomy/v1/category_tree/{categoryTreeId}
```

### 3. **Image Upload**
For listings with images:
1. Upload images to eBay's image server
2. Get image URLs
3. Add to listing's `ImageURLs` array

### 4. **Webhooks for Real-time Notifications**
Set up eBay webhooks to get notified about:
- New orders
- New offers
- Messages from buyers

Configure at: https://developer.ebay.com/my/subscriptions

### 5. **Shipping Label Purchase**
Implement using Buy API or third-party services like Shippo.

### 6. **Offer Management Commands**
Add Discord commands to respond to offers:
```
/accept-offer [offer-id]
/decline-offer [offer-id]
/counter-offer [offer-id] [price]
```

## Testing in Sandbox

1. Use **Sandbox** environment first
2. Sandbox eBay: https://sandbox.ebay.com
3. Create test listings and orders
4. Test all functionality before going to production

## Moving to Production

1. Get **Production** credentials from eBay
2. Update `.env`:
   ```env
   EBAY_ENVIRONMENT=PRODUCTION
   ```
3. Generate new production OAuth tokens
4. Test thoroughly with small volume
5. Monitor error rates and API responses

## Common Issues & Solutions

### "No access token available"
- Generate user token using eBay's token generator
- Add to `.env` file
- Restart bot

### "Token expired"
- If refresh token is set, bot auto-refreshes
- Otherwise, generate new tokens manually

### "Invalid category ID"
- Use eBay's Category API to find valid IDs
- Different categories per marketplace (US, UK, etc.)

### "Missing fulfillment policy"
- Create policies in eBay Seller Hub first
- Get policy IDs via Account API
- Update bot code with real policy IDs

## API Documentation Links

- **eBay Developer Docs**: https://developer.ebay.com/docs
- **API Explorer**: https://developer.ebay.com/my/api_test_tool
- **OAuth Guide**: https://developer.ebay.com/api-docs/static/oauth-tokens.html
- **Fulfillment API**: https://developer.ebay.com/api-docs/sell/fulfillment/overview.html
- **Inventory API**: https://developer.ebay.com/api-docs/sell/inventory/overview.html

## Code Structure

```
internal/ebay/
├── client.go     # Main API client, request handling
├── oauth.go      # OAuth flows, token management  
└── types.go      # Data structures for API responses

internal/bot/
└── handler.go    # Discord command handlers

internal/config/
└── config.go     # Configuration management
```

## Next Steps

1. ✅ Go installed and project built
2. ⏳ Get Discord bot token
3. ⏳ Get eBay developer credentials
4. ⏳ Generate OAuth tokens
5. ⏳ Configure `.env` file
6. ⏳ Run the bot: `./ebaymanager-bot.exe`
7. ⏳ Test commands in Discord
8. ⏳ Implement fulfillment policies
9. ⏳ Add more commands and features
