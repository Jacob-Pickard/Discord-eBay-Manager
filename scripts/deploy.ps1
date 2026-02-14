# eBay Discord Bot - Auto Deploy to Server
param(
    [string]$ServerIP = "192.168.0.12",
    [string]$ServerUser = "jacob",
    [string]$ServerPath = "/home/jacob/ebay-bot"
)

Write-Host ""
Write-Host "🚀 Starting deployment to $ServerIP..." -ForegroundColor Cyan
Write-Host ""

# Step 1: Build Linux binary
Write-Host "📦 Building Linux binary..." -ForegroundColor Yellow
$env:GOOS = 'linux'
$env:GOARCH = 'amd64'
Set-Location ..
go build -o bin/ebaymanager-bot-linux

if (-not (Test-Path bin/ebaymanager-bot-linux)) {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}

$size = [math]::Round((Get-Item bin/ebaymanager-bot-linux).Length / 1MB, 2)
Write-Host "✅ Build complete: $size MB" -ForegroundColor Green
Write-Host ""

# Step 2: Stop bot service
Write-Host "⏸️  Stopping bot service..." -ForegroundColor Yellow
ssh $ServerUser@$ServerIP "sudo systemctl stop ebay-bot"
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
ssh $ServerUser@$ServerIP "sudo systemctl start ebay-bot"
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
Write-Host "📋 Recent logs (last 10 lines):" -ForegroundColor Cyan
Write-Host "─────────────────────────────────────────────────────" -ForegroundColor Gray
ssh $ServerUser@$ServerIP "tail -10 $ServerPath/bot-error.log"
Write-Host "─────────────────────────────────────────────────────" -ForegroundColor Gray
Write-Host ""

# Step 8: Test webhook
Write-Host "🔍 Testing webhook endpoint..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "https://jacob.it.com/webhook/health" -TimeoutSec 5 -ErrorAction Stop
    Write-Host "✅ Webhook responding: $response" -ForegroundColor Green
} catch {
    Write-Host "⚠️  Webhook test failed (this is okay if bot is still starting)" -ForegroundColor Yellow
}
Write-Host ""

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
Write-Host "🎉 DEPLOYMENT COMPLETE!" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Green
Write-Host ""
Write-Host "Your bot is now running on jacob.it.com" -ForegroundColor Cyan
Write-Host "Webhook: https://jacob.it.com/webhook/ebay/notification" -ForegroundColor Cyan
Write-Host ""
Write-Host "To view live logs:" -ForegroundColor White
Write-Host "  ssh jacob@192.168.0.12" -ForegroundColor Gray
Write-Host "  tail -f /home/jacob/ebay-bot/bot-error.log" -ForegroundColor Gray
Write-Host ""
