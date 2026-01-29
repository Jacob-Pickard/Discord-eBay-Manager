# eBay Listing Management via Discord

## Overview
You can now create, view, and manage your eBay sandbox listings directly through Discord commands! This makes testing much easier - no need to use the eBay Sandbox interface.

## New Commands

### 1. `/publish-listing` - Create a Full Listing
Creates a complete published eBay listing with optional Best Offer enabled.

**Parameters:**
- `title` (required): Item title
- `price` (required): Buy It Now price in USD
- `description` (required): Item description
- `enable-offers` (optional): Enable Best Offer feature (true/false)
- `min-offer` (optional): Minimum acceptable offer price (auto-accept threshold)

**Examples:**
```
/publish-listing title:"Vintage Watch" price:"99.99" description:"Great condition vintage timepiece"

/publish-listing title:"Leather Jacket" price:"150.00" description:"Brand new leather jacket" enable-offers:true min-offer:"120"
```

**What it does:**
1. Creates inventory item with SKU
2. Creates an offer with pricing
3. Publishes the listing to eBay Sandbox
4. Enables Best Offer if requested
5. Sets auto-accept price if min-offer specified
6. Returns SKU and Listing ID

### 2. `/view-listings` - View Active Listings
See all your currently active eBay listings.

**Parameters:**
- `limit` (optional): Number of listings to show (default: 5)

**Examples:**
```
/view-listings

/view-listings limit:10
```

**What it shows:**
- Item title
- Current price
- SKU (for use with other commands)
- Best Offer status (enabled/disabled)
- Direct link to sandbox listing

### 3. `/end-listing` - End a Listing
Remove a listing from eBay (unpublish and delete).

**Parameters:**
- `sku` (required): The SKU of the listing to end

**Example:**
```
/end-listing sku:ITEM-1738092039
```

**What it does:**
1. Withdraws (unpublishes) any active offers
2. Declines any pending offers automatically
3. Deletes the offer/listing
4. Removes the inventory item

## Complete Testing Workflow

### Step 1: Create a Test Listing
```
/publish-listing 
  title:"Test Item - Gaming Mouse" 
  price:"49.99" 
  description:"Wireless gaming mouse with RGB lighting" 
  enable-offers:true 
  min-offer:"40"
```

Response will include:
- ✅ SKU: `ITEM-1738092039`
- 🆔 Listing ID: `110123456789`
- 🌐 Direct link to view

### Step 2: View Your Listings
```
/view-listings
```

You'll see all active listings with:
- Title, price, SKU
- Best Offer status
- Direct links

### Step 3: Make an Offer (as Buyer)
1. Copy the sandbox URL from `/publish-listing` or `/view-listings`
2. Open in browser
3. Log in as a sandbox buyer account
4. Click "Make Offer" and submit an offer

### Step 4: Check Offers in Discord
```
/get-offers
```

You'll see pending offers with offer IDs.

### Step 5: Respond to Offers
```
/accept-offer offer-id:12345
/counter-offer offer-id:12345 price:45.00
/decline-offer offer-id:12345
```

### Step 6: Clean Up
```
/end-listing sku:ITEM-1738092039
```

## Best Offer Settings

When creating listings, you can configure Best Offer in different ways:

### No Offers
```
/publish-listing title:"Item" price:"100" description:"Desc"
```
Buyers can only Buy It Now at full price.

### Offers Enabled (Manual Review)
```
/publish-listing title:"Item" price:"100" description:"Desc" enable-offers:true
```
Buyers can make offers, you review each one manually.

### Offers with Auto-Accept
```
/publish-listing title:"Item" price:"100" description:"Desc" enable-offers:true min-offer:"85"
```
- Offers ≥ $85: Automatically accepted by eBay
- Offers < $85: Sent to you for review

## Tips & Best Practices

### For Testing:
1. **Use descriptive titles** so you can identify items easily
2. **Set realistic prices** ($10-$200 range works well)
3. **Enable Best Offer** on test items to test the offer workflow
4. **Set min-offer slightly below price** (e.g., price=$100, min=$80)

### Category Notes:
Listings are created in category **88433** (Fashion Jewelry) by default. This is a safe category for testing that doesn't require additional specifications.

### Default Policies:
The bot uses default sandbox policies for:
- **Fulfillment Policy**: Standard shipping (ID: 6367426000)
- **Payment Policy**: Standard payment (ID: 6367427000)
- **Return Policy**: Standard returns (ID: 6367428000)

These are pre-configured in the sandbox environment.

## Troubleshooting

### "Failed to publish listing"
- Check that you're authorized: `/ebay-status`
- Verify your access token is valid
- Make sure price is a valid number

### "No Active Listings"
- Create your first listing with `/publish-listing`
- Check that previous listings weren't deleted

### Can't make offers as buyer
- You need a separate sandbox buyer account
- Log out and create a new test user at: https://developer.ebay.com/sandbox/register
- Make sure you're on the sandbox site (URL should contain "sandbox")

### Offers not showing in `/get-offers`
- Wait 1-2 minutes after making offer for webhook notification
- Check that webhook is subscribed: `/webhook-list`
- Manually refresh: `/get-offers` (will show if offer exists)

## Integration with Other Features

### Webhook Notifications
When a listing is published with Best Offer enabled, you'll automatically receive notifications when:
- 💬 OFFER_RECEIVED: Buyer makes an offer
- 💰 OFFER_COUNTERED: Buyer responds to your counter
- ✅ ITEM_SOLD: Offer accepted and item sells

### Order Management
When an offer is accepted and item sells:
- Check order details: `/get-orders`
- Order will show up with payment status
- Use for shipping label generation (future feature)

## Example Full Session

```
# 1. Check status
/ebay-status

# 2. Create listing with Best Offer
/publish-listing 
  title:"Test Watch - Silver" 
  price:"79.99" 
  description:"Silver watch for testing" 
  enable-offers:true 
  min-offer:"65"

# 3. View all listings
/view-listings

# 4. (As buyer) Go to URL and make offer for $60

# 5. Check for offers
/get-offers

# 6. Counter the offer
/counter-offer offer-id:ABC123 price:70

# 7. (As buyer) Accept counter offer

# 8. Check orders
/get-orders

# 9. Clean up
/end-listing sku:ITEM-1738092039
```

## Command Count
Your bot now has **15 total commands**:
1. `/list-item` - Create inventory item (without publishing)
2. `/publish-listing` - Create and publish full listing ⭐ NEW
3. `/view-listings` - View active listings ⭐ NEW
4. `/end-listing` - End/delete listing ⭐ NEW
5. `/get-orders` - View orders
6. `/get-offers` - View pending offers
7. `/accept-offer` - Accept offer
8. `/counter-offer` - Counter offer
9. `/decline-offer` - Decline offer
10. `/ebay-status` - Check connection
11. `/ebay-authorize` - Get auth URL
12. `/ebay-code` - Submit auth code
13. `/webhook-subscribe` - Subscribe to notifications
14. `/webhook-list` - List subscriptions
15. `/webhook-test` - Test webhook

Now you can manage your entire eBay testing workflow without leaving Discord! 🎉
