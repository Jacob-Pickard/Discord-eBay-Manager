#!/bin/bash
ENV_FILE=/home/jacob/ebay-bot/.env
TOKEN=$(grep EBAY_ACCESS_TOKEN "$ENV_FILE" | cut -d= -f2-)

echo "================================================"
echo "eBay API Diagnostic Test - $(date)"
echo "================================================"
echo ""

echo "=== TOKEN INFO ==="
echo "Token (first 30 chars): ${TOKEN:0:30}..."
echo ""

echo "=== TEST 1: Orders API (sell/fulfillment/v1/order) ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.ebay.com/sell/fulfillment/v1/order?limit=5"
echo ""

echo "=== TEST 2: Finances - Balance (sell/finances/v1/seller_funds_summary) ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.ebay.com/sell/finances/v1/seller_funds_summary"
echo ""

echo "=== TEST 3: Finances - Payouts (sell/finances/v1/payout) ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.ebay.com/sell/finances/v1/payout?limit=5"
echo ""

echo "=== TEST 4: Notification API - List Destinations ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.ebay.com/commerce/notification/v1/destination"
echo ""

echo "=== TEST 5: Webhook Challenge Verification ==="
curl -s -w "\nHTTP_STATUS:%{http_code}\n" \
  "https://jacob.it.com/webhook/ebay/notification?challenge_code=diagnostic_test_12345"
echo ""

echo "================================================"
echo "Diagnostic complete."
echo "================================================"
