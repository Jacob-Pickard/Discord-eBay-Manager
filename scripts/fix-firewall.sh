#!/bin/bash
# Firewall Diagnostic and Fix Script for eBay Bot Webhook Access
# Run with: sudo bash fix-firewall.sh

echo "=== eBay Bot Webhook Firewall Diagnostic ==="
echo ""

# Check if UFW is installed and active
echo "1. Checking UFW firewall status..."
if command -v ufw &> /dev/null; then
    ufw status verbose
    echo ""
else
    echo "UFW not installed"
    echo ""
fi

# Check iptables rules
echo "2. Checking iptables rules..."
iptables -L -n -v | grep -E "(80|443|8081)"
echo ""

# Check if ports are blocked
echo "3. Checking which ports are currently allowed..."
if command -v ufw &> /dev/null; then
    ufw status numbered
    echo ""
fi

# Ask user if they want to fix
echo ""
read -p "Do you want to allow ports 80 (HTTP) and 443 (HTTPS) through the firewall? (y/n) " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "=== Fixing Firewall Rules ==="
    echo ""
    
    # Enable UFW if not enabled
    if command -v ufw &> /dev/null; then
        echo "Configuring UFW firewall..."
        
        # Allow SSH first (important!)
        ufw allow 22/tcp
        echo "✓ SSH (port 22) allowed"
        
        # Allow HTTP
        ufw allow 80/tcp
        echo "✓ HTTP (port 80) allowed"
        
        # Allow HTTPS
        ufw allow 443/tcp
        echo "✓ HTTPS (port 443) allowed"
        
        # Enable UFW if not already enabled
        ufw --force enable
        
        echo ""
        echo "=== New Firewall Status ==="
        ufw status verbose
        
        echo ""
        echo "✅ Firewall configured! External access should now work."
        echo ""
        echo "Test your webhook with:"
        echo "  curl -k https://jacob.it.com/webhook/health"
        echo ""
    else
        echo "UFW not installed. Checking iptables..."
        
        # Check if there are blocking rules in iptables
        if iptables -L INPUT -n | grep -q "DROP\|REJECT"; then
            echo "⚠️  Found blocking rules in iptables"
            echo "Current INPUT chain:"
            iptables -L INPUT -n --line-numbers
            echo ""
            echo "You may need to manually configure iptables to allow ports 80 and 443"
            echo "Example commands:"
            echo "  iptables -I INPUT -p tcp --dport 80 -j ACCEPT"
            echo "  iptables -I INPUT -p tcp --dport 443 -j ACCEPT"
        else
            echo "No obvious blocking rules found in iptables"
        fi
    fi
else
    echo "Skipping firewall changes."
fi

echo ""
echo "=== Additional Diagnostics ==="
echo ""

# Check if services are listening
echo "Services listening on webhook ports:"
ss -tuln | grep -E ':(80|443|8081)'
echo ""

# Check nginx status
echo "Nginx status:"
systemctl status nginx --no-pager | head -5
echo ""

# Check bot status
echo "eBay Bot status:"
systemctl status ebay-bot --no-pager | head -5
echo ""

echo "=== Router/ISP Check ==="
echo "Note: If the firewall is open but external access still fails,"
echo "check your router's firewall and port forwarding settings."
echo "Ports 80 and 443 must be forwarded to 192.168.0.199"
