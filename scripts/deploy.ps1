# eBay Discord Bot - Auto Deploy to Server
# Builds and deploys the bot to your production server
# Configuration is loaded from deploy-config.env file

param(
    [string]$ConfigFile = "$PSScriptRoot/../deploy-config.env",
    [switch]$Watch,
    [switch]$SkipBuild
)

Write-Host ""
Write-Host "╔════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "║   eBay Discord Bot - Deployment Script   ║" -ForegroundColor Cyan
Write-Host "╚════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""
if (Test-Path $ConfigFile) {
    Write-Host "Loading configuration from: $ConfigFile" -ForegroundColor Cyan
    Get-Content $ConfigFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.+)$') {
            $key = $Matches[1].Trim()
            $value = $Matches[2].Trim()
            Set-Variable -Name $key -Value $value -Scope Script
        }
    }
} else {
    Write-Host "❌ Configuration file not found: $ConfigFile" -ForegroundColor Red
    Write-Host "Please copy deploy-config.env.example to deploy-config.env and customize it" -ForegroundColor Yellow
    exit 1
}

# Validate required configuration
$required = @('DEPLOY_SERVER_IP', 'DEPLOY_SERVER_USER', 'DEPLOY_SERVER_PATH')
foreach ($var in $required) {
    if (-not (Get-Variable -Name $var -ErrorAction SilentlyContinue).Value) {
        Write-Host "❌ Required configuration missing: $var" -ForegroundColor Red
        Write-Host "Please set $var in $ConfigFile" -ForegroundColor Yellow
        exit 1
    }
}

$ServerIP = $DEPLOY_SERVER_IP
$ServerUser = $DEPLOY_SERVER_USER
$ServerPath = $DEPLOY_SERVER_PATH
$Domain = if ($DEPLOY_DOMAIN) { $DEPLOY_DOMAIN } else { "your-domain.com" }

Write-Host "🚀 Starting deployment to $ServerIP..." -ForegroundColor Cyan
if ($SkipBuild) {
    Write-Host "⏭️  Skipping build (using existing binary)" -ForegroundColor Yellow
}
if ($Watch) {
    Write-Host "👀 Live log viewing enabled" -ForegroundColor Yellow
}
Write-Host ""

# Step 1: Build Linux binary (unless skipped)
if (-not $SkipBuild) {
    Write-Host "📦 Building Linux binary..." -ForegroundColor Yellow
    $env:GOOS = 'linux'
    $env:GOARCH = 'amd64'
    Set-Location "$PSScriptRoot/.."
    go build -o bin/ebaymanager-bot-linux .

    if (-not (Test-Path bin/ebaymanager-bot-linux)) {
        Write-Host "❌ Build failed!" -ForegroundColor Red
        exit 1
    }

    $size = [math]::Round((Get-Item bin/ebaymanager-bot-linux).Length / 1MB, 2)
    Write-Host "✅ Build complete: $size MB" -ForegroundColor Green
    Write-Host ""
}

# Step 2: Stop bot service
Write-Host "⏸️  Stopping bot service..." -ForegroundColor Yellow
ssh $ServerUser@$ServerIP 'sudo systemctl stop ebay-bot 2>/dev/null || kill $(cat /home/jacob/ebay-bot/bot.pid 2>/dev/null) 2>/dev/null; true'
Start-Sleep -Seconds 2
Write-Host "✅ Service stopped" -ForegroundColor Green
Write-Host ""

# Step 3: Upload binary
Write-Host "⬆️  Uploading binary to server..." -ForegroundColor Yellow
scp bin/ebaymanager-bot-linux ${ServerUser}@${ServerIP}:${ServerPath}/ebaymanager-bot-linux
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Binary uploaded" -ForegroundColor Green
} else {
    Write-Host "❌ Upload failed!" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Step 4: Set permissions
Write-Host "🔑 Setting executable permissions..." -ForegroundColor Yellow
ssh $ServerUser@$ServerIP "chmod +x $ServerPath/ebaymanager-bot-linux"
Write-Host "✅ Permissions set" -ForegroundColor Green
Write-Host ""

# Step 5: Start bot service
Write-Host "▶️  Starting bot service..." -ForegroundColor Yellow
ssh $ServerUser@$ServerIP 'sudo systemctl start ebay-bot 2>/dev/null || (nohup /home/jacob/ebay-bot/ebaymanager-bot-linux > /home/jacob/ebay-bot/bot.log 2> /home/jacob/ebay-bot/bot-error.log & echo $! > /home/jacob/ebay-bot/bot.pid)'
Start-Sleep -Seconds 3
Write-Host "✅ Service started" -ForegroundColor Green
Write-Host ""

# Step 6: Check status
Write-Host "📊 Checking bot status..." -ForegroundColor Yellow
$status = ssh $ServerUser@$ServerIP "systemctl is-active ebay-bot"
if ($status -eq "active") {
    Write-Host "✅ Bot is running!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Bot status: $status" -ForegroundColor Yellow
}
Write-Host ""

# Step 7: Show recent logs
Write-Host "Recent logs:" -ForegroundColor Cyan
Write-Host "-----------------------------------------------------" -ForegroundColor Gray
ssh $ServerUser@$ServerIP "tail -10 $ServerPath/bot-error.log"
Write-Host "-----------------------------------------------------" -ForegroundColor Gray
Write-Host ""

# Step 8: Test webhook
Write-Host "🔍 Testing webhook endpoint..." -ForegroundColor Yellow
try {
    $webhookUrl = if ($DEPLOY_WEBHOOK_URL) { $DEPLOY_WEBHOOK_URL -replace '/ebay/notification', '/health' } else { "https://$Domain/webhook/health" }
    $response = Invoke-RestMethod -Uri $webhookUrl -TimeoutSec 5 -ErrorAction Stop
    Write-Host "✅ Webhook responding: $response" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Webhook test failed (this is okay if bot is still starting)" -ForegroundColor Yellow
}
Write-Host ""

Write-Host "============================================" -ForegroundColor Green
Write-Host "DEPLOYMENT COMPLETE!" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "Y

# Watch logs if requested
if ($Watch) {
    Write-Host ""
    Write-Host "📡 Showing live logs (Press Ctrl+C to exit)..." -ForegroundColor Cyan
    Write-Host "═══════════════════════════════════════════════" -ForegroundColor Gray
    Write-Host ""
    ssh $ServerUser@$ServerIP "tail -f $ServerPath/bot-error.log"
}our bot is now running on $Domain" -ForegroundColor Cyan
$webhookEndpoint = if ($DEPLOY_WEBHOOK_URL) { $DEPLOY_WEBHOOK_URL } else { "https://$Domain/webhook/ebay/notification" }
Write-Host "Webhook: $webhookEndpoint" -ForegroundColor Cyan
Write-Host ""
Write-Host "To view live logs:" -ForegroundColor White
Write-Host "  ssh $ServerUser@$ServerIP" -ForegroundColor Gray
Write-Host "  tail -f $ServerPath/bot-error.log" -ForegroundColor Gray
Write-Host ""
