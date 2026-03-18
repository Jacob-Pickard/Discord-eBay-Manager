# Test Scripts - Important Notes

The test scripts in this directory contain **EXAMPLE VALUES ONLY**.

## ⚠️ Before Using These Scripts

1. **Replace placeholder domains** - Change `yourdomain.com` or any domain references to your actual domain
2. **Set environment variables** - Make sure your `.env` file is properly configured
3. **Update webhook URLs** - Use your actual webhook endpoint URL

## 📝 Test Files

| File | Purpose | Notes |
|------|---------|-------|
| `test_webhook_subscription.go` | Test eBay webhook subscription | Replace domain placeholder |
| `test_webhook_for_support.sh` | Diagnostic test for eBay support | Example token values shown |
| `Test-Webhook-Simple.ps1` | Simple webhook tests | Update domain before running |
| `Test-TokenFormats.ps1` | Token format validation | Update domain before running |
| `check_config.go` | Validate environment config | Reads from .env |

## 🔧 Example Values vs Real Values

These test files may contain:
- ✅ **Example verification tokens** (like `my_secure_verify_token_12345...`) - These are just examples
- ✅ **Placeholder domains** like `yourdomain.com` - Replace with your domain
- ⚠️ **Hardcoded test values** - Clearly marked as examples

## 🚀 How to Use

1. **Copy the relevant test script**
2. **Update all placeholder values** with your actual configuration
3. **Run the test** to validate your setup

Or better yet: Set up your `.env` and `deploy-config.env` files, and the scripts will read from there automatically.
