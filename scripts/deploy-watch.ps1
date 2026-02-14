# Deploy with Live Log Viewing
# Deploys and then shows live logs

param(
    [string]$ServerUser = "jacob",
    [string]$ServerIP = "192.168.0.12",
    [string]$ServerPath = "/home/jacob/ebay-bot"
)

# Run deployment
.\deploy.ps1

# Ask if user wants to view live logs
Write-Host ""
Write-Host "Would you like to view live logs? (y/n)" -ForegroundColor Yellow
$response = Read-Host
if ($response -eq 'y' -or $response -eq 'Y' -or $response -eq 'yes') {
    Write-Host ""
    Write-Host "📡 Showing live logs (Press Ctrl+C to exit)..." -ForegroundColor Cyan
    Write-Host ""
    ssh $ServerUser@$ServerIP "tail -f $ServerPath/bot-error.log"
}
