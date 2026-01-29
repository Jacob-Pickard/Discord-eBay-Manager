# Webhook Setup Guide

## 🎣 Real-Time eBay Notifications

Your bot now includes a webhook server that can receive real-time notifications from eBay when orders, offers, and other events occur!

## ✅ What's Been Added

**New Discord Commands:**
- `/webhook-subscribe` - Set up eBay webhook notifications
- `/webhook-list` - View active webhook subscriptions
- `/webhook-test` - Test notification delivery to Discord

**Webhook Server:**
- Running on port 8080 (configurable via `WEBHOOK_PORT` in .env)
- Endpoints:
  - `POST /webhook/ebay/notification` - Receives eBay notifications
  - `GET /webhook/ebay/challenge` - Handles eBay endpoint verification
  - `GET /webhook/health` - Health check endpoint

**Supported Notifications:**
- 💰 New orders placed
- 💵 Payment received
- 💬 Buyer offers received
- 📝 Offer updates (countered, accepted, declined)
- 📦 Shipping updates
- ⏰ Listing ended

## 🚀 Setup Instructions

### Option 1: Testing Locally with ngrok (Easiest)

1. **Install ngrok** (if not already installed):
   - Download from: https://ngrok.com/download
   - Or use chocolatey: `choco install ngrok`

2. **Start ngrok tunnel:**
   ```powershell
   ngrok http 8080
   ```

3. **Copy the HTTPS URL** from ngrok (e.g., `https://abc123.ngrok-free.app`)

4. **Set notification channel in .env:**
   ```env
   NOTIFICATION_CHANNEL_ID=your_discord_channel_id
   ```
   (Right-click a channel in Discord with Developer Mode enabled to copy ID)

5. **Subscribe to webhooks** in Discord:
   ```
   /webhook-subscribe url:https://abc123.ngrok-free.app/webhook/ebay/notification
   ```

6. **Configure in eBay Developer Portal:**
   - Go to https://developer.ebay.com/my/subscriptions
   - Create a new subscription
   - Set the endpoint URL to your ngrok URL
   - eBay will send a challenge to verify ownership

### Option 2: Production Deployment

For production use, you need a publicly accessible server:

1. **Deploy bot to a server** (VPS, cloud instance, etc.)

2. **Set up a domain** with SSL (required by eBay):
   - Use Let's Encrypt for free SSL certificates
   - Point domain to your server
   - Set up reverse proxy (nginx/Apache) if needed

3. **Configure firewall:**
   ```bash
   # Allow port 8080 (or your WEBHOOK_PORT)
   sudo ufw allow 8080
   ```

4. **Update .env with your domain:**
   ```env
   WEBHOOK_PORT=8080
   NOTIFICATION_CHANNEL_ID=your_channel_id
   ```

5. **Subscribe via Discord:**
   ```
   /webhook-subscribe url:https://yourdomain.com/webhook/ebay/notification
   ```

## 📋 Configuration

**Environment Variables (.env):**

```env
# Webhook Configuration
WEBHOOK_PORT=8080
WEBHOOK_VERIFY_TOKEN=my_secure_verify_token_12345
NOTIFICATION_CHANNEL_ID=1234567890123456789
```

**How to get Discord Channel ID:**
1. Enable Developer Mode in Discord (User Settings → Advanced → Developer Mode)
2. Right-click the channel you want notifications in
3. Click "Copy Channel ID"
4. Paste into .env file

## 🔍 Testing

1. **Test webhook server is running:**
   ```
   curl http://localhost:8080/webhook/health
   ```
   Should return: `OK`

2. **Test Discord notifications:**
   ```
   /webhook-test
   ```
   You should see a test notification in the current channel

3. **Test eBay challenge response:**
   ```
   curl "http://localhost:8080/webhook/ebay/challenge?challenge_code=test123"
   ```
   Should return JSON with `challengeResponse`

## 📚 eBay Notification Topics

Available notification types from eBay:

| Topic | Description |
|-------|-------------|
| `marketplace.order.placed` | New order created |
| `marketplace.order.paid` | Payment received for order |
| `marketplace.offer.created` | Buyer submits an offer |
| `marketplace.offer.updated` | Offer countered/accepted/declined |
| `marketplace.listing.ended` | Listing has ended |
| `marketplace.item.sold` | Item sold (auction ended) |

## 🐛 Troubleshooting

**Webhook server not starting:**
- Check if port 8080 is already in use
- Change `WEBHOOK_PORT` in .env to a different port
- Make sure bot has permission to bind to the port

**eBay challenge failing:**
- Verify `WEBHOOK_VERIFY_TOKEN` matches what you configured in eBay portal
- Check ngrok tunnel is still active
- Ensure endpoint URL is publicly accessible

**Notifications not appearing in Discord:**
- Verify `NOTIFICATION_CHANNEL_ID` is set correctly
- Check bot has permission to send messages in that channel
- Use `/webhook-test` to verify Discord connectivity

**eBay says endpoint is unreachable:**
- Ensure your webhook URL is publicly accessible (test with curl from external server)
- Check firewall rules allow incoming connections
- Verify SSL certificate is valid (eBay requires HTTPS for production)

## 🔒 Security Notes

- Keep your `WEBHOOK_VERIFY_TOKEN` secret
- Use HTTPS in production (required by eBay)
- Validate incoming webhook signatures
- Consider IP whitelisting for eBay's webhook servers
- Don't expose internal server details in error messages

## 📈 Next Steps

Once webhooks are working, you can:

1. **Customize notification messages** in `internal/webhook/server.go`
2. **Add auto-responses** (e.g., auto-accept offers over a certain amount)
3. **Create notification filters** (only notify for high-value orders)
4. **Add analytics** (track order volume, offer conversion rates)
5. **Integrate with shipping APIs** (auto-generate labels when order is paid)

## 🆘 Need Help?

Check the bot logs for detailed error messages:
```
tail -f bot.log
```

Test individual components:
- Discord connection: `/ebay-status`
- Webhook server: `curl localhost:8080/webhook/health`
- Notifications: `/webhook-test`
