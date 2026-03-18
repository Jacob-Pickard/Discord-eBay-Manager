#!/bin/bash
cd /home/jacob/ebay-bot

# Kill any existing bot processes
pkill -f ebaymanager-bot-linux
sleep 2

# Use new log files to avoid permission issues
LOG_FILE="bot-$(date +%Y%m%d-%H%M%S).log"
ERR_FILE="bot-error-$(date +%Y%m%d-%H%M%S).log"

# Start the bot (it will load .env from current directory)
./ebaymanager-bot-linux >> "$LOG_FILE" 2>> "$ERR_FILE" &

BOT_PID=$!
echo "Bot started! PID: $BOT_PID"
echo "Log file: $LOG_FILE"
echo "Error file: $ERR_FILE"
echo "Waiting 5 seconds for startup..."
sleep 5

echo ""
echo "=== Checking if bot is running ==="
if ps -p $BOT_PID > /dev/null; then
    echo "✅ Bot is running (PID: $BOT_PID)"
    ps aux | grep ebaymanager-bot | grep -v grep
else
    echo "❌ Bot is not running - checking error log..."
    tail -20 "$ERR_FILE"
    exit 1
fi

echo ""
echo "=== Recent bot logs ==="
tail -30 "$LOG_FILE"

echo ""
echo "=== Bot is ready! ==="
echo "Try running /webhook-subscribe in Discord now"
