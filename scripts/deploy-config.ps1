# Update .env Configuration on Server
# Use this to update environment variables without redeploying the binary

param(
    [string]$ServerUser = "jacob",
    [string]$ServerIP = "192.168.0.12",
    [string]$ServerPath = "/home/jacob/ebay-bot",
    [string]$LocalEnvFile = ".env"
)

if (-not (Test-Path $LocalEnvFile)) {
    Write-Host "❌ Error: $LocalEnvFile not found!" -ForegroundColor Red
    Write-Host "Create a .env file or specify the path with -LocalEnvFile" -ForegroundColor Yellow
    exit 1
}

Write-Host "📝 Updating .env configuration on server..." -ForegroundColor Cyan
Write-Host ""

# Upload .env file
Write-Host "⬆️  Uploading .env file..." -ForegroundColor Yellow
scp $LocalEnvFile ${ServerUser}@${ServerIP}:${ServerPath}/.env

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Configuration uploaded" -ForegroundColor Green
} else {
    Write-Host "❌ Upload failed!" -ForegroundColor Red
    exit 1
}
Write-Host ""

# Restart bot to apply changes
Write-Host "🔄 Restarting bot to apply changes..." -ForegroundColor Yellow
ssh -t $ServerUser@$ServerIP "sudo systemctl restart ebay-bot" 2>$null
Start-Sleep -Seconds 3

$status = ssh $ServerUser@$ServerIP "systemctl is-active ebay-bot" 2>$null
if ($status -eq "active") {
    Write-Host "✅ Bot restarted successfully!" -ForegroundColor Green
} else {
    Write-Host "⚠️  Bot status: $status" -ForegroundColor Yellow
}
Write-Host ""

# Show recent logs
Write-Host "📋 Recent logs:" -ForegroundColor Cyan
ssh $ServerUser@$ServerIP "tail -15 $ServerPath/bot-error.log" 2>$null
Write-Host ""
Write-Host "✅ Configuration update complete!" -ForegroundColor Green
