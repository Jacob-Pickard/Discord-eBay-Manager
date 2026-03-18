# Update .env Configuration on Server
# Use this to update environment variables without redeploying the binary
# Configuration is loaded from deploy-config.env file

param(
    [string]$ConfigFile = "$PSScriptRoot/../deploy-config.env",
    [string]$LocalEnvFile = ".env"
)

# Load configuration from file
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

$ServerUser = $DEPLOY_SERVER_USER
$ServerIP = $DEPLOY_SERVER_IP
$ServerPath = $DEPLOY_SERVER_PATH

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
