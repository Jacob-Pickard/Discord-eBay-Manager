# Deploy with Live Log Viewing
# Deploys the bot and automatically shows live logs
# This is a convenience wrapper around deploy.ps1 -Watch

param(
    [string]$ConfigFile = "$PSScriptRoot/../deploy-config.env"
)

# Call deploy.ps1 with -Watch parameter
& "$PSScriptRoot\deploy.ps1" -ConfigFile $ConfigFile -Watch
