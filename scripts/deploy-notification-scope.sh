#!/bin/bash
echo "Stopping bot..."
sudo systemctl stop ebay-bot

echo "Backing up old binary..."
cp /home/jacob/ebay-bot/ebaymanager-bot-linux /home/jacob/ebay-bot/ebaymanager-bot-linux.old

echo "Deploying new binary with notification scope..."
cp /home/jacob/ebay-bot/ebaymanager-bot-linux.new /home/jacob/ebay-bot/ebaymanager-bot-linux

echo "Starting bot..."
sudo systemctl start ebay-bot

sleep 2
echo ""
echo "Bot status:"
systemctl status ebay-bot --no-pager | head -15

echo ""
echo "✅ Deployment complete!"
echo ""
echo "NEXT STEPS:"
echo "1. Run /ebay-authorize in Discord to get new OAuth token with notification scope"
echo "2. Complete the OAuth flow"
echo "3. Run /webhook-subscribe to create webhook subscription"
