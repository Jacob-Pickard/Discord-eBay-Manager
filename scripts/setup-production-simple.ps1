# Simple Production Setup Script
Write-Host ""
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host " eBay Production Credentials Setup" -ForegroundColor Cyan  
Write-Host "===============================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Opening eBay Developer Portal..." -ForegroundColor Yellow
Start-Process "https://developer.ebay.com/my/keys"
Write-Host ""
Write-Host "INSTRUCTIONS:" -ForegroundColor White
Write-Host "1. Sign in with your eBay seller account" -ForegroundColor Gray
Write-Host "2. Click PRODUCTION tab at the top" -ForegroundColor Gray
Write-Host "3. Find your Application Keys (Production)" -ForegroundColor Gray
Write-Host ""
Read-Host "Press Enter when you're viewing the PRODUCTION keys page"

Write-Host ""
Write-Host "Enter your PRODUCTION credentials:" -ForegroundColor Cyan
Write-Host ""

$appId = Read-Host "App ID (Client ID)"
$certId = Read-Host "Cert ID (Client Secret)"  
$devId = Read-Host "Dev ID"
$ruName = Read-Host "RuName (Redirect URI) - Leave blank if waiting for approval"

# Get Discord settings from current .env
$currentEnv = Get-Content .env -Raw
$discordToken = if ($currentEnv -match 'DISCORD_BOT_TOKEN=(.+)') { $Matches[1].Trim() } else { "" }
$channelId = if ($currentEnv -match 'NOTIFICATION_CHANNEL_ID=(.+)') { $Matches[1].Trim() } else { "" }

# Create production config
$config = @"
# Discord Configuration
DISCORD_BOT_TOKEN=$discordToken

# eBay PRODUCTION API Configuration
EBAY_APP_ID=$appId
EBAY_CERT_ID=$certId
EBAY_DEV_ID=$devId
EBAY_REDIRECT_URI=$ruName

# OAuth Tokens - Will be generated after authorization
EBAY_ACCESS_TOKEN=
EBAY_REFRESH_TOKEN=

# Environment - PRODUCTION
EBAY_ENVIRONMENT=PRODUCTION

# Webhook Configuration
WEBHOOK_PORT=8081
WEBHOOK_VERIFY_TOKEN=my_secure_verify_token_12345
NOTIFICATION_CHANNEL_ID=$channelId
"@

$config | Out-File -FilePath ".env.production" -Encoding UTF8

Write-Host ""
Write-Host "SUCCESS! Created .env.production" -ForegroundColor Green
Write-Host ""

if ([string]::IsNullOrWhiteSpace($ruName)) {
    Write-Host "NOTE: No RuName provided" -ForegroundColor Yellow
    Write-Host "Your production RuName may still be pending approval from eBay" -ForegroundColor Yellow
    Write-Host "Check https://developer.ebay.com/my/keys daily for approval" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Once approved, edit .env.production and add your RuName" -ForegroundColor White
    Write-Host "Then deploy with: .\deploy-config.ps1 -LocalEnvFile .env.production" -ForegroundColor White
} else {
    Write-Host "Ready to deploy to production!" -ForegroundColor Green
    $deploy = Read-Host "Deploy now? (y/n)"
    if ($deploy -eq 'y') {
        Write-Host ""
        Write-Host "Deploying production configuration..." -ForegroundColor Cyan
        .\deploy-config.ps1 -LocalEnvFile .env.production
        
        Write-Host ""
        Write-Host "NEXT STEPS IN DISCORD:" -ForegroundColor Green
        Write-Host "1. /ebay-authorize" -ForegroundColor White
        Write-Host "2. Click link and authorize" -ForegroundColor White  
        Write-Host "3. /ebay-code code:YOUR_CODE" -ForegroundColor White
        Write-Host "4. /webhook-subscribe url:https://jacob.it.com/webhook/ebay/notification" -ForegroundColor White
    }
}

Write-Host ""
Write-Host "Setup complete!" -ForegroundColor Cyan
