#!/bin/bash

# Load environment variables
source /home/jacob/ebay-bot/.env

# Determine base URL based on environment
if [ "$EBAY_ENVIRONMENT" = "PRODUCTION" ]; then
    BASE_URL="https://api.ebay.com"
else
    BASE_URL="https://api.sandbox.ebay.com"
fi

echo "=== Listing eBay Webhook Subscriptions ==="
echo "Environment: $EBAY_ENVIRONMENT"
echo "Base URL: $BASE_URL"
echo ""

# List subscriptions
curl -s -X GET "$BASE_URL/commerce/notification/v1/destination" \
  -H "Authorization: Bearer $EBAY_ACCESS_TOKEN" \
  -H "Content-Type: application/json" | jq '.'

echo ""
echo "=== End of Subscriptions ==="
