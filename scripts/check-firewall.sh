#!/bin/bash
# Firewall Diagnostic Script (No sudo required)
# Run with: bash check-firewall.sh

echo "=== eBay Bot Webhook Firewall Diagnostic ==="
echo ""

# Check if UFW is installed
echo "1. Checking UFW firewall status..."
if command -v ufw &> /dev/null; then
    echo "UFW is installed"
    ufw status 2>&1 || echo "   (Need sudo to check status)"
else
    echo "UFW not installed"
fi
echo ""

# Check if services are listening
echo "2. Services listening on webhook ports:"
ss -tuln | grep -E ':(80|443|8081)'
echo ""

# Check nginx status
echo "3. Nginx status:"
systemctl status nginx --no-pager 2>&1 | head -5
echo ""

# Check bot status
echo "4. eBay Bot status:"
systemctl status ebay-bot --no-pager 2>&1 | head -5
echo ""

# Check nginx config
echo "5. Nginx webhook configuration:"
cat /etc/nginx/conf.d/jacob.it.com.conf 2>&1 | grep -A 20 "location /webhook"
echo ""

# Test local connection
echo "6. Testing local webhook endpoint:"
curl -s http://127.0.0.1:8081/webhook/health || echo "   Failed to connect locally"
echo ""

echo "7. Testing via nginx (localhost):"
curl -s -k http://localhost/webhook/health 2>&1 || echo "   Failed"
echo ""

echo ""
echo "=== Result ==="
echo ""
echo "To fix firewall issues, run these commands:"
echo ""
echo "  sudo ufw allow 22/tcp    # SSH"
echo "  sudo ufw allow 80/tcp    # HTTP"
echo "  sudo ufw allow 443/tcp   # HTTPS"
echo "  sudo ufw enable"
echo "  sudo ufw status verbose"
echo ""
echo "To check if it's a router issue:"
echo "  Test from external network: curl -k https://jacob.it.com/webhook/health"
