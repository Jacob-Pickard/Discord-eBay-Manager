# Quick Test Script - Try Different Token Formats
# This tests if the issue is with the verification token format

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Testing Different Verification Tokens" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Get OAuth token
$AccessToken = Read-Host "Enter your eBay OAuth access token (from /ebay-authorize in Discord)"

if (-not $AccessToken) {
    Write-Host "ERROR: Access token required" -ForegroundColor Red
    exit
}

Write-Host "`nTesting with different verification token formats...`n" -ForegroundColor Yellow

# Test 1: Minimum length (32 chars)
Write-Host "TEST 1: Minimum length token (32 chars)" -ForegroundColor Yellow
$token1 = "abcdefghij1234567890ABCDEFGHIJ"
Write-Host "Token: $token1 (Length: $($token1.Length))`n"

$body1 = @{
    name = "Test_32_Chars"
    status = "ENABLED"
    deliveryConfig = @{
        endpoint = "https://jacob.it.com/webhook/ebay/notification"
        verifyToken = $token1
    }
    topics = @(@{ topicName = "MARKETPLACE_OFFER" })
} | ConvertTo-Json -Depth 10

try {
    $response1 = Invoke-RestMethod `
        -Uri "https://api.ebay.com/commerce/notification/v1/destination" `
        -Method Post `
        -Headers @{
            "Authorization" = "Bearer $AccessToken"
            "Content-Type" = "application/json"
        } `
        -Body $body1 `
        -ErrorAction Stop
    
    Write-Host "SUCCESS! Subscription created with 32-char token" -ForegroundColor Green
    $response1 | ConvertTo-Json
} catch {
    $error1 = $_.ErrorDetails.Message
    Write-Host "FAILED: $error1`n" -ForegroundColor Red
}

Start-Sleep -Seconds 2

# Test 2: Simple alphanumeric only (40 chars)
Write-Host "TEST 2: Simple alphanumeric token (40 chars)" -ForegroundColor Yellow
$token2 = "simpletoken1234567890ABCDEFGHIJ1234567"
Write-Host "Token: $token2 (Length: $($token2.Length))`n"

$body2 = @{
    name = "Test_Simple"
    status = "ENABLED"
    deliveryConfig = @{
        endpoint = "https://jacob.it.com/webhook/ebay/notification"
        verifyToken = $token2
    }
    topics = @(@{ topicName = "MARKETPLACE_OFFER" })
} | ConvertTo-Json -Depth 10

try {
    $response2 = Invoke-RestMethod `
        -Uri "https://api.ebay.com/commerce/notification/v1/destination" `
        -Method Post `
        -Headers @{
            "Authorization" = "Bearer $AccessToken"
            "Content-Type" = "application/json"
        } `
        -Body $body2 `
        -ErrorAction Stop
    
    Write-Host "SUCCESS! Subscription created with simple token" -ForegroundColor Green
    $response2 | ConvertTo-Json
} catch {
    $error2 = $_.ErrorDetails.Message
    Write-Host "FAILED: $error2`n" -ForegroundColor Red
}

Start-Sleep -Seconds 2

# Test 3: Your current token (55 chars with underscores)
Write-Host "TEST 3: Current token format (55 chars, with underscores)" -ForegroundColor Yellow
$token3 = "my_secure_verify_token_12345_67890_abcdefghij_klmnopqrs"
Write-Host "Token: $token3 (Length: $($token3.Length))`n"

$body3 = @{
    name = "Test_Current"
    status = "ENABLED"
    deliveryConfig = @{
        endpoint = "https://jacob.it.com/webhook/ebay/notification"
        verifyToken = $token3
    }
    topics = @(@{ topicName = "MARKETPLACE_OFFER" })
} | ConvertTo-Json -Depth 10

try {
    $response3 = Invoke-RestMethod `
        -Uri "https://api.ebay.com/commerce/notification/v1/destination" `
        -Method Post `
        -Headers @{
            "Authorization" = "Bearer $AccessToken"
            "Content-Type" = "application/json"
        } `
        -Body $body3 `
        -ErrorAction Stop
    
    Write-Host "SUCCESS! Subscription created with current token" -ForegroundColor Green
    $response3 | ConvertTo-Json
} catch {
    $error3 = $_.ErrorDetails.Message
    Write-Host "FAILED: $error3`n" -ForegroundColor Red
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "Test Results Summary" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan
Write-Host "If ALL tests failed with error 195019:"
Write-Host "  -> Issue is NOT with token format"
Write-Host "  -> Issue is likely API access/enablement"
Write-Host "  -> Submit support ticket with results`n"
Write-Host "If any test SUCCEEDED:"
Write-Host "  -> Issue IS with token format"
Write-Host "  -> Update your config to use the working token format`n"
