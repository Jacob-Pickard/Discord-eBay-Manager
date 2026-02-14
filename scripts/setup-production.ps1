# Get eBay Production Credentials - Interactive Guide
# This will open the pages you need and tell you exactly what to copy

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "   eBay Production Credentials Setup" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "This wizard will help you get your production credentials from eBay." -ForegroundColor White
Write-Host "I'll open the right pages and tell you exactly what to copy." -ForegroundColor White
Write-Host ""

# Step 1: Open eBay Developer Portal
Write-Host "STEP 1: Opening eBay Developer Portal..." -ForegroundColor Yellow
Write-Host ""
Write-Host "I'm opening: https://developer.ebay.com/my/keys" -ForegroundColor Gray
Start-Process "https://developer.ebay.com/my/keys"
Write-Host ""
Write-Host "▶️  Sign in with your eBay seller account" -ForegroundColor White
Write-Host "▶️  At the top of the page, click on 'PRODUCTION' (not Sandbox)" -ForegroundColor White
Write-Host ""
Read-Host "Press Enter when you're on the PRODUCTION keys page"

# Step 2: Get Application Keys
Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "STEP 2: Get Your Application Keys" -ForegroundColor Yellow
Write-Host ""
Write-Host "Look for a section called 'Application Keys (Production)'" -ForegroundColor White
Write-Host ""
Write-Host "You should see (or need to create):" -ForegroundColor White
Write-Host "  • App ID (Client ID)" -ForegroundColor Gray
Write-Host "  • Cert ID (Client Secret)" -ForegroundColor Gray
Write-Host "  • Dev ID" -ForegroundColor Gray
Write-Host ""
Write-Host "If you don't see these, click 'Create Application Keys'" -ForegroundColor Yellow
Write-Host ""

$hasKeys = Read-Host "Do you see your Production Application Keys? (y/n)"
if ($hasKeys -ne 'y' -and $hasKeys -ne 'Y') {
    Write-Host ""
    Write-Host "📝 You need to create production keys first:" -ForegroundColor Yellow
    Write-Host "   1. Click 'Create Application Keys' or 'Get Access'" -ForegroundColor White
    Write-Host "   2. Fill in the application details" -ForegroundColor White
    Write-Host "   3. Wait for approval (usually instant)" -ForegroundColor White
    Write-Host ""
    Write-Host "⚠️  Production keys may require eBay review (can take up to 24 hours)" -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Press Enter once you have production keys visible"
}

Write-Host ""
Write-Host "Now, let's collect your credentials..." -ForegroundColor Cyan
Write-Host ""

# Collect App ID
Write-Host "📋 PRODUCTION APP ID (Client ID):" -ForegroundColor Green
Write-Host "   It looks like: YourName-AppName-PRD-xxxxxxxxx" -ForegroundColor Gray
$appId = Read-Host "   Paste it here"

# Collect Cert ID
Write-Host ""
Write-Host "📋 PRODUCTION CERT ID (Client Secret):" -ForegroundColor Green
Write-Host "   It looks like: PRD-xxxxxxxx-xxxx-xxxx-xxxx" -ForegroundColor Gray
$certId = Read-Host "   Paste it here"

# Collect Dev ID
Write-Host ""
Write-Host "📋 PRODUCTION DEV ID:" -ForegroundColor Green
Write-Host "   It looks like: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" -ForegroundColor Gray
$devId = Read-Host "   Paste it here"

# Step 3: Get RuName
Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "STEP 3: Get Your Production RuName (Redirect URI)" -ForegroundColor Yellow
Write-Host ""
Write-Host "On the same page, look for 'User Tokens' section" -ForegroundColor White
Write-Host "Then find 'Get a Token from eBay via Your Application'" -ForegroundColor White
Write-Host ""
Write-Host "⚠️  IMPORTANT: You need a PRODUCTION RuName (different from sandbox)" -ForegroundColor Yellow
Write-Host ""

$hasRuName = Read-Host "Do you see a Production RuName listed? (y/n)"
if ($hasRuName -ne 'y' -and $hasRuName -ne 'Y') {
    Write-Host ""
    Write-Host "📝 You need to create a Production RuName:" -ForegroundColor Yellow
    Write-Host "   1. Click 'Add RuName' or 'Get OAuth Application Credentials'" -ForegroundColor White
    Write-Host "   2. Fill in:" -ForegroundColor White
    Write-Host "      - Your Company/Name" -ForegroundColor Gray
    Write-Host "      - App Name" -ForegroundColor Gray
    Write-Host "      - Privacy Policy URL (can use a placeholder)" -ForegroundColor Gray
    Write-Host "   3. Submit for approval" -ForegroundColor White
    Write-Host ""
    Write-Host "⏳ Production RuNames require eBay approval (1-3 business days)" -ForegroundColor Red
    Write-Host ""
    Write-Host "You can continue with sandbox for now and switch to production when approved." -ForegroundColor Cyan
    Write-Host ""
    $continue = Read-Host "Continue anyway? (y/n)"
    if ($continue -ne 'y' -and $continue -ne 'Y') {
        Write-Host ""
        Write-Host "Setup cancelled. Run this script again when you have your RuName." -ForegroundColor Yellow
        exit
    }
    $ruName = ""
} else {
    Write-Host ""
    Write-Host "📋 PRODUCTION RUNAME:" -ForegroundColor Green
    Write-Host "   It looks like: YourName-AppName-Produc-xxxxx" -ForegroundColor Gray
    $ruName = Read-Host "   Paste it here"
}

# Step 4: Create .env.production file
Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "STEP 4: Creating Production Configuration File..." -ForegroundColor Yellow
Write-Host ""

# Read current .env to get Discord token and channel ID
$currentEnv = Get-Content .env -Raw
$discordToken = if ($currentEnv -match 'DISCORD_BOT_TOKEN=(.+)') { $Matches[1].Trim() } else { "" }
$channelId = if ($currentEnv -match 'NOTIFICATION_CHANNEL_ID=(.+)') { $Matches[1].Trim() } else { "" }
$verifyToken = if ($currentEnv -match 'WEBHOOK_VERIFY_TOKEN=(.+)') { $Matches[1].Trim() } else { "my_secure_verify_token_12345" }

# Create .env.production
$productionEnv = @"
# Discord Configuration
DISCORD_BOT_TOKEN=$discordToken

# eBay PRODUCTION API Configuration
EBAY_APP_ID=$appId
EBAY_CERT_ID=$certId
EBAY_DEV_ID=$devId
EBAY_REDIRECT_URI=$ruName

# OAuth Tokens - Will be generated after authorization
# Leave these empty - they'll be filled automatically
EBAY_ACCESS_TOKEN=
EBAY_REFRESH_TOKEN=

# Environment - PRODUCTION
EBAY_ENVIRONMENT=PRODUCTION

# Webhook Configuration
WEBHOOK_PORT=8081
WEBHOOK_VERIFY_TOKEN=$verifyToken
NOTIFICATION_CHANNEL_ID=$channelId
"@

$productionEnv | Out-File -FilePath ".env.production" -Encoding UTF8 -NoNewline

Write-Host "✅ Created .env.production file" -ForegroundColor Green
Write-Host ""

# Step 5: Ask if they want to deploy now
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "STEP 5: Deploy to Production?" -ForegroundColor Yellow
Write-Host ""

if ($ruName -eq "") {
    Write-Host "⚠️  You don't have a Production RuName yet." -ForegroundColor Yellow
    Write-Host ""
    Write-Host "You can:" -ForegroundColor White
    Write-Host "  1. Wait for RuName approval (check developer.ebay.com daily)" -ForegroundColor Gray
    Write-Host "  2. Continue testing in sandbox mode" -ForegroundColor Gray
    Write-Host ""
    Write-Host "Once your RuName is approved:" -ForegroundColor Cyan
    Write-Host "  1. Edit .env.production and add your RuName" -ForegroundColor White
    Write-Host "  2. Run: .\deploy-config.ps1 -LocalEnvFile .env.production" -ForegroundColor White
    Write-Host "  3. In Discord: /ebay-authorize" -ForegroundColor White
    Write-Host ""
} else {
    Write-Host "Ready to switch to production!" -ForegroundColor Green
    Write-Host ""
    $deployNow = Read-Host "Deploy production config now? (y/n)"
    
    if ($deployNow -eq 'y' -or $deployNow -eq 'Y') {
        Write-Host ""
        Write-Host "🚀 Deploying production configuration..." -ForegroundColor Cyan
        .\deploy-config.ps1 -LocalEnvFile .env.production
        
        Write-Host ""
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
        Write-Host "   NEXT STEPS IN DISCORD" -ForegroundColor Green
        Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
        Write-Host ""
        Write-Host "1. Run this command in Discord:" -ForegroundColor White
        Write-Host "   /ebay-authorize" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "2. Click the authorization link and sign in" -ForegroundColor White
        Write-Host ""
        Write-Host "3. Copy the code eBay gives you" -ForegroundColor White
        Write-Host ""
        Write-Host "4. Run in Discord:" -ForegroundColor White
        Write-Host "   /ebay-code code:<paste_code_here>" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "5. Subscribe to webhooks:" -ForegroundColor White
        Write-Host "   /webhook-subscribe url:https://jacob.it.com/webhook/ebay/notification" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "6. Test with:" -ForegroundColor White
        Write-Host "   /ebay-status" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "🎉 You'll then be live in production!" -ForegroundColor Green
        Write-Host ""
    } else {
        Write-Host ""
        Write-Host "Configuration saved to .env.production" -ForegroundColor Yellow
        Write-Host "Deploy later with: .\deploy-config.ps1 -LocalEnvFile .env.production" -ForegroundColor White
        Write-Host ""
    }
}

Write-Host ""
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "   SETUP COMPLETE!" -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host ""
Write-Host "✅ Production credentials collected" -ForegroundColor Green
Write-Host "✅ .env.production file created" -ForegroundColor Green
Write-Host ""
Write-Host "Files created:" -ForegroundColor White
Write-Host "  📄 .env.production - Production configuration" -ForegroundColor Gray
Write-Host ""
