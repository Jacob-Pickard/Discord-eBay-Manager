# Security & Sensitive Information Guide

This guide explains how sensitive information is managed in this project to keep it safe for GitHub while still being usable.

## 🔐 Sensitive Information Protected

The following types of sensitive information are **NEVER** committed to the repository:

### Credentials & Tokens
- Discord bot tokens
- eBay API credentials (App ID, Cert ID, Dev ID)
- OAuth access tokens and refresh tokens
- Webhook verification tokens

### Server & Deployment Info
- Server IP addresses
- SSH usernames
- Domain names
- Server file paths

## 📁 File Organization

### Configuration Files

| File | Purpose | Status |
|------|---------|--------|
| `.env` | Your actual secrets | ❌ Gitignored |
| `.env.example` | Template with placeholders | ✅ Committed |
| `deploy-config.env` | Your server details | ❌ Gitignored |
| `deploy-config.env.example` | Template for deployment | ✅ Committed |
| `config/.env.production.template` | Production template | ✅ Committed (placeholders) |
| `config/ebay-bot.service` | Your service file | ❌ Gitignored |
| `config/ebay-bot.service.example` | Service template | ✅ Committed |
| `config/*.conf` | Your nginx configs | ❌ Gitignored |
| `config/*.example` | Example configs | ✅ Committed |

### What's in `.gitignore`

```
# Environment files with real credentials
.env
.env.local
.env.production
config/.env.*
!config/.env.example
!config/.env.production.template

# Deployment configuration
deploy-config.env

# Service and server configs
config/ebay-bot.service
config/*.conf
!config/*.example
```

## 🚀 Setup for New Users

### 1. Create Your `.env` File

```bash
cp .env.example .env
# Edit .env with your actual values
```

### 2. Create Deployment Config (if deploying to a server)

```bash
cp deploy-config.env.example deploy-config.env
# Edit deploy-config.env with your server details
```

### 3. Create Service File (on your server)

```bash
cp config/ebay-bot.service.example /etc/systemd/system/ebay-bot.service
# Edit the service file with your username and paths
```

## 📝 For Maintainers

### Before Committing

1. **Never commit real credentials** - Always use placeholders in example files
2. **Check for hardcoded values** - Search for:
   - IP addresses (e.g., `192.168.x.x`)
   - Domain names (use `yourdomain.com` instead)
   - Usernames (use `youruser` instead)
   - API keys or tokens

### Placeholders to Use

| Type | Example Placeholder |
|------|-------------------|
| Domain | `yourdomain.com` |
| IP Address | `192.168.1.100` or `YOUR_SERVER_IP` |
| Username | `youruser` or `YOUR_USERNAME` |
| Path | `/home/youruser/ebay-bot` |
| Discord Token | `your_discord_bot_token_here` |
| eBay App ID | `your_ebay_app_id_here` |
| Webhook Token | `generate_secure_random_token_32_to_80_chars` |
| Channel ID | `your_discord_channel_id_here` |

### Documentation Guidelines

When writing documentation:
- Use placeholder values, not real ones
- Include comments like `(Replace with your actual domain)`
- Reference the`.example` files for structure

## 🛡️ Security Best Practices

1. **Rotate tokens if exposed** - If you accidentally commit sensitive data, rotate all affected tokens immediately
2. **Use strong verification tokens** - Generate with: `openssl rand -base64 48 | tr -d "=+/" | cut -c1-60`
3. **Keep .env files secure** - Never share or commit them
4. **Use environment-specific configs** - Keep production and sandbox separate
5. **Review before pushing** - Always check `git diff` before committing

## ❓ FAQ

**Q: I accidentally committed a `.env` file. What do I do?**
A: 
1. Remove it from Git: `git rm --cached .env`
2. Rotate all credentials in that file immediately
3. Commit the removal: `git commit -m "Remove accidentally committed .env"`

**Q: Can I commit my deployment scripts with my server info?**
A: No! Keep your server details in `deploy-config.env` (gitignored). The scripts will read from that file.

**Q: What if I want to share configurations between machines?**
A: Use a secure password manager or encrypted cloud storage for your `.env` and `deploy-config.env` files. Never commit them to Git.

## 📚 Related Documentation

- [Deployment Scripts](docs/DEPLOYMENT_SCRIPTS.md)
- [Webhook Setup](docs/WEBHOOK_SETUP.md)
- [Quick Start Guide](QUICKSTART.md)
