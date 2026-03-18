#!/bin/bash

echo "=== Updating eBay Bot Service ==="

# Kill any manually running bot processes
echo "Stopping any manually running bot processes..."
pkill -f ebaymanager-bot-linux
sleep 2

# Copy service file (requires sudo)
echo "Updating systemd service file..."
sudo cp /tmp/ebay-bot.service /etc/systemd/system/ebay-bot.service
sudo chmod 644 /etc/systemd/system/ebay-bot.service

# Reload systemd
echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

# Enable and start the service
echo "Enabling and starting ebay-bot service..."
sudo systemctl enable ebay-bot.service
sudo systemctl restart ebay-bot.service

# Wait a moment for it to start
sleep 3

# Check status
echo ""
echo "=== Service Status ==="
sudo systemctl status ebay-bot.service --no-pager -l

echo ""
echo "=== Recent Logs ==="
sudo journalctl -u ebay-bot.service -n 20 --no-pager

echo ""
echo "=== Done ==="
echo "Service updated and restarted!"
echo "The bot should now be running with environment variables loaded from .env"
