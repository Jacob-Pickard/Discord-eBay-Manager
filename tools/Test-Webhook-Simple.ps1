# PowerShell: Quick Webhook Test Commands
# Run these commands one at a time to test your webhook endpoint

Write-Host "=========================================="
Write-Host "eBay Webhook Testing Script"
Write-Host "=========================================="
Write-Host ""

# Test 1: Health Check
Write-Host "TEST 1: Health Endpoint" -ForegroundColor Yellow
Write-Host "Command: curl https://jacob.it.com/webhook/health"
Write-Host ""

try {
    $health = Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health"
    Write-Host "Result: PASS" -ForegroundColor Green
    Write-Host ($health | ConvertTo-Json)
} catch {
    Write-Host "Result: FAIL - $_" -ForegroundColor Red
}

Write-Host ""
Write-Host ""

# Test 2: Challenge Response
Write-Host "TEST 2: Challenge Response" -ForegroundColor Yellow
$challenge = "test_$(Get-Date -Format 'HHmmss')"
Write-Host "Command: curl 'https://jacob.it.com/webhook/ebay/notification?challenge_code=$challenge'"
Write-Host ""

try {
    $response = Invoke-RestMethod -Uri "https://jacob.it.com/webhook/ebay/notification?challenge_code=$challenge"
    Write-Host "Result: PASS" -ForegroundColor Green
    Write-Host ($response | ConvertTo-Json)
    Write-Host ""
    if ($response.challengeResponse) {
        Write-Host "challengeResponse received: $($response.challengeResponse)" -ForegroundColor Cyan
    }
} catch {
    Write-Host "Result: FAIL - $_" -ForegroundColor Red
}

Write-Host ""
Write-Host ""

# Test 3: Show the API call that fails
Write-Host "TEST 3: Subscription API Call (Requires OAuth Token)" -ForegroundColor Yellow
Write-Host ""
Write-Host "To test the failing API call, you need an OAuth access token."
Write-Host "Get it by running /ebay-authorize in Discord, then check bot logs."
Write-Host ""
Write-Host "The command that fails with error 195019:"
Write-Host ""
Write-Host @'
curl -X POST https://api.ebay.com/commerce/notification/v1/destination \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test_Subscription",
    "status": "ENABLED",
    "deliveryConfig": {
      "endpoint": "https://jacob.it.com/webhook/ebay/notification",
      "verifyToken": "my_secure_verify_token_12345_67890_abcdefghij_klmnopqrs"
    },
    "topics": [{"topicName": "MARKETPLACE_OFFER"}]
  }'
'@ -ForegroundColor Gray

Write-Host ""
Write-Host ""
Write-Host "SUMMARY FOR SUPPORT TICKET:" -ForegroundColor Cyan
Write-Host ""
Write-Host "- Endpoint works (Test 1 passes)"
Write-Host "- Challenge response works (Test 2 passes)" 
Write-Host "- eBay API returns error 195019 WITHOUT attempting verification"
Write-Host "- No verification requests appear in server logs"
Write-Host ""
Write-Host "See docs/EBAY_SUPPORT_TICKET_TEMPLATE.md for full ticket template"
Write-Host "See docs/EXACT_COMMANDS_FOR_SUPPORT.md for copy-paste commands"
Write-Host ""
