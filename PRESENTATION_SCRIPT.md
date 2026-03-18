# Final Presentation Script - Discord eBay Manager Bot
**Target Duration:** 10-15 minutes  
**Format:** Screen recording with narration

---

## 🎬 INTRODUCTION (1 minute)

**[SLIDE 1: Title Slide]**

**Script:**
> "Hello everyone! Today I'm presenting my automation project: the Discord eBay Manager Bot. This project automates eBay seller operations through Discord, allowing sellers to manage their business without constantly checking eBay.com."

**Visual:** Show bot logo or Discord bot running in server

---

## 🎯 SEGMENT 1: The Problem (1.5 minutes)

**[SLIDE 2: Problem Statement]**

**Script:**
> "As an eBay seller, I identified several pain points in the typical workflow:
> 
> **Problem 1:** Sellers must constantly refresh eBay.com to check for new orders or offers
> 
> **Problem 2:** Mobile notifications exist, but responding requires opening eBay's app or website
> 
> **Problem 3:** No centralized dashboard for quick actions - everything requires multiple clicks through eBay's interface
> 
> **Problem 4:** Missing time-sensitive offers can cost sales, especially when buyers send offers late at night
> 
> My solution was to bring eBay management into Discord - a platform I already have open all day - where I can view orders, respond to offers, and receive real-time notifications through webhooks."

**Visual:** Screenshots showing:
- Cluttered eBay seller hub
- vs. Clean Discord interface with bot commands

---

## 💻 SEGMENT 2: Technology Overview (2 minutes)

**[SLIDE 3: Technology Stack]**

**Script:**
> "For this project, I chose **Go** as my primary scripting language, along with several supporting technologies:
> 
> **Go (Golang)** - The core language for the bot and webhook server
> - Version: 1.21
> - Chosen for its performance, concurrency features, and strong standard library
> 
> **Key Libraries Used:**
> - `discordgo` - Discord API wrapper for bot commands and interactions
> - `godotenv` - Environment variable management
> - Go's built-in `net/http` - HTTP server for webhooks
> - Go's `crypto/sha256` - Webhook signature verification
> 
> **External APIs:**
> - Discord API - For bot messaging and commands
> - eBay API (4 different endpoints):
>   - Fulfillment API - For order management
>   - Trading API - For listings and legacy operations
>   - Browse API - For product images
>   - Identity API - For OAuth authentication
> 
> **Infrastructure:**
> - Linux VPS for production deployment
> - nginx as reverse proxy for webhook endpoints
> - systemd for service management"

**Visual:** Architecture diagram showing:
```
[Discord] ↔ [Go Bot] ↔ [eBay APIs]
                ↑
          [Webhook Server]
                ↑
          [eBay Notifications]
```

---

## 🔄 SEGMENT 3: Go vs Python Comparison (2 minutes)

**[SLIDE 4: Go vs Python]**

**Script:**
> "Let me compare Go with Python, which many of you are familiar with:
> 
> **Similarities:**
> - Both are high-level languages with garbage collection
> - Both have excellent standard libraries
> - Both support concurrent programming (Go: goroutines, Python: asyncio)
> - Both have strong community package ecosystems
> 
> **Where Go Excels:**
> - **Compiled language** - Produces a single binary with no dependencies
>   - My entire bot compiles to one executable that runs on the server
> - **Performance** - Typically 10-30x faster than Python for CPU-bound tasks
> - **Built-in concurrency** - Goroutines are lightweight and easy to use
>   - Example: I can handle webhook requests while the bot processes Discord commands
> - **Static typing** - Catches errors at compile time, not runtime
> - **Deployment** - Just copy one binary file vs. managing Python dependencies
> 
> **Where Python Excels:**
> - **Faster prototyping** - No compilation step
> - **More libraries** - Especially for data science and ML
> - **Dynamic typing** - More flexible (though also more error-prone)
> - **Simpler syntax** - Generally easier for beginners
> 
> **Why I chose Go for this project:**
> - Long-running server application benefits from Go's performance
> - Need concurrent handling of Discord commands and webhook events
> - Deployment simplicity - one binary to upload vs. managing pip packages
> - Learning opportunity - wanted to expand beyond Python"

**Visual:** Side-by-side code comparison (show same simple function in both languages)

---

## 📝 SEGMENT 4: Code Walkthrough (3-4 minutes)

**[SLIDE 5: Code Architecture]**

**Script:**
> "Now let's walk through the most important parts of the code. The project is organized into several key components:"

### Part A: Main Entry Point (30 seconds)

**[SCREEN: Show main.go]**

**Script:**
> "Main.go is our entry point. It initializes the eBay client, creates a Discord session, starts the webhook server in a goroutine, and registers command handlers. Notice how we use goroutines to run the webhook server concurrently with the Discord bot."

**Code to highlight:**
```go
webhookServer := webhook.NewServer(...)
go webhookServer.Start()  // Goroutine for concurrent execution
```

### Part B: OAuth Flow (45 seconds)

**[SCREEN: Show internal/ebay/oauth.go]**

**Script:**
> "The OAuth implementation was one of the trickiest parts. eBay requires OAuth 2.0 with automatic token refresh. Here's how I handle it:
> 
> 1. Generate authorization URL for user consent
> 2. Exchange authorization code for access token
> 3. Store tokens securely in environment
> 4. Automatically refresh tokens before they expire (tokens last 2 hours, auto-refresh at 1h 55m)
> 
> This struct handles all OAuth operations, and the `EnsureValidToken` function runs before every API call to ensure we always have valid credentials."

**Code to highlight:**
```go
func (tm *TokenManager) EnsureValidToken() error {
    if time.Now().Add(tm.refreshBefore).After(tm.expiresAt) {
        // Refresh token 5 minutes before expiry
        return tm.RefreshToken()
    }
    return nil
}
```

### Part C: Discord Command Handler (45 seconds)

**[SCREEN: Show internal/bot/handler.go]**

**Script:**
> "The command handler is where users actually interact with the bot. When someone types a command like `/get-orders`, this function processes it, calls the appropriate eBay API, and formats the response as a rich Discord embed.
> 
> Notice the error handling - if something goes wrong, we send a user-friendly message rather than crashing. This is important for production reliability."

**Code to highlight:**
- Command routing logic
- Error handling pattern
- Discord embed creation

### Part D: Webhook Security (45 seconds)

**[SCREEN: Show internal/webhook/server.go]**

**Script:**
> "Security is critical when accepting webhooks from external services. There are two verification steps:
> 
> 1. **Challenge Verification** (subscription setup): eBay sends a challenge code, and I hash it with SHA-256 along with my verification token and endpoint URL.
> 
> 2. **Signature Verification** (actual notifications): eBay sends an HMAC-SHA256 signature with each notification payload. If the signature doesn't match, the notification is rejected. This prevents malicious actors from sending fake notifications."

**Code to highlight:**
```go
// Challenge verification (SHA-256 hash)
hash := sha256.New()
hash.Write([]byte(challengeCode))
hash.Write([]byte(verifyToken))
hash.Write([]byte(endpointURL))
challengeResponse := base64.StdEncoding.EncodeToString(hash.Sum(nil))

// Notification signature verification (HMAC-SHA256)
func verifySignature(body []byte, signature string) bool {
    mac := hmac.New(sha256.New, []byte(verifyToken))
    mac.Write(body)
    expectedMAC := mac.Sum(nil)
    expectedSignature := base64.StdEncoding.EncodeToString(expectedMAC)
    return hmac.Equal([]byte(signature), []byte(expectedSignature))
}
```

---

## 🎥 SEGMENT 5: Live Demonstration (4-5 minutes)

**[SLIDE 6: Live Demo]**

**Script:**
> "Now let's see the bot in action!"

### Demo Flow:

1. **Show Discord Server Setup (20 seconds)**
   - "Here's my Discord server with the bot online and configured"
   - Show bot's online status

2. **Demonstrate `/get-orders` Command (45 seconds)**
   - Type `/get-orders` in Discord
   - "This fetches my recent eBay orders and displays them with rich embeds"
   - Point out: product images, buyer info, order status
   - "Notice how quickly it responds - usually under 2 seconds"

3. **Demonstrate Offer Management (1.5 minutes)**
   - "Let me show you the offer workflow. When a buyer submits a best offer, I receive it in Discord"
   - Show example offer notification (if available, or screenshot)  
   - Demo commands: `/accept-offer offer-id:<ID>`, `/counter-offer offer-id:<ID> price:<AMOUNT>`, `/decline-offer offer-id:<ID>`
   - "I can respond to offers directly from Discord without opening eBay"

4. **Show Webhook Notification (1 minute)**
   - "The most powerful feature is real-time webhooks"
   - If possible, trigger a webhook (or show recorded example)
   - "When an item sells or an offer comes in, eBay sends a webhook to my server"
   - "The webhook is verified for security, then posted to my notification channel"
   - Show the notification message format

5. **Show Deployment/Infrastructure (45 seconds)**
   - "Behind the scenes, this runs on a Linux VPS"
   - SSH into server (optional): `ssh user@server`
   - Show service status: `systemctl status ebay-bot`
   - "The bot runs as a systemd service, automatically restarts if it crashes"
   - Show logs briefly: `journalctl -u ebay-bot -f`

6. **Code Update Workflow (30 seconds)**
   - "When I make changes, I use my PowerShell deployment script"
   - Show: `.\scripts\deploy.ps1`
   - "This builds the Go binary, uploads it to the server, and restarts the service automatically"

---

## 📚 SEGMENT 6: Learning Resources (1 minute)

**[SLIDE 7: Recommended Resources]**

**Script:**
> "If you're interested in learning Go after seeing this project, here are my top recommended resources:
> 
> **For Beginners:**
> 1. **A Tour of Go** (tour.golang.org) - Free interactive tutorial
> 2. **Go by Example** (gobyexample.com) - Practical code examples
> 3. **Effective Go** (official documentation) - Best practices guide
> 
> **For Building Projects:**
> 4. **Go Web Examples** (gowebexamples.com) - HTTP servers, APIs, databases
> 5. **Awesome Go** (github.com/avelino/awesome-go) - Curated list of libraries
> 
> **For Discord Bots Specifically:**
> 6. **DiscordGo Documentation** (github.com/bwmarrin/discordgo)
> 7. **Discord Developer Portal** (discord.com/developers/docs)
> 
> **Books:**
> 8. 'The Go Programming Language' by Donovan & Kernighan - The definitive guide
> 
> All of these resources are free except the book. Go has an excellent learning curve compared to other systems programming languages."

---

## 🤖 SEGMENT 7: AI Coding Assistance (OPTIONAL - 1.5 minutes)

**[SLIDE 8: AI Tools Used]**

**Script:**
> "Since we have time, I'd like to share my experience using AI coding assistants on this project:
> 
> **Tools I Used:**
> - **GitHub Copilot** - Primary coding assistant
> - **Claude/ChatGPT** - For debugging and architectural questions
> 
> **What Worked Well:**
> - Generating boilerplate code (HTTP handlers, error handling)
> - Explaining complex Go concepts (channels, goroutines)
> - Debugging OAuth flow issues - I could paste error messages and get solutions
> - Writing documentation and comments
> - Suggesting idiomatic Go patterns I wasn't aware of
> 
> **Best Practices I Learned:**
> 1. **Be specific in prompts** - Instead of 'write a webhook handler,' say 'write a Go HTTP handler that verifies HMAC-SHA256 signatures'
> 2. **Verify all generated code** - AI sometimes produces code that compiles but has logic errors
> 3. **Use AI for learning** - Ask 'why' questions: 'Why use a pointer receiver here?'
> 4. **Iterate incrementally** - Get a basic version working, then ask for improvements
> 5. **Don't blindly copy-paste** - Understand what the code does
> 
> **Limitations I Hit:**
> - AI struggled with eBay's specific API requirements (needed to read docs)
> - Sometimes suggested outdated library versions
> - Couldn't help with environment-specific deployment issues
> 
> Overall, AI probably saved me 20-30% of development time, especially on syntax and boilerplate."

---

## 🎬 CONCLUSION (1 minute)

**[SLIDE 9: Summary & Thank You]**

**Script:**
> "To wrap up:
> 
> ✅ I built a Discord bot that automates eBay seller operations
> ✅ Used Go for its performance, concurrency, and deployment simplicity
> ✅ Integrated with eBay's Trading, Browse, and OAuth APIs
> ✅ Implemented secure webhook handling with signature verification
> ✅ Deployed to production with systemd and automated deployment scripts
> 
> **Key Takeaways:**
> - Automation doesn't have to be complex - identify a repetitive task and automate it
> - Go is excellent for long-running server applications
> - Security (OAuth, webhook verification) is critical for production services
> - Good documentation and deployment scripts make maintenance easier
> 
> Thank you for your time! I'm happy to answer any questions."

**Visual:** Final slide with:
- Project repository link (if public)
- Your contact info
- References (see next section)

---

## 📖 REFERENCE SLIDE

**[SLIDE 10: References - APA Format]**

```
References

Donovan, A. A., & Kernighan, B. W. (2015). The Go programming language. 
    Addison-Wesley Professional.

eBay Inc. (2024). eBay Developers Program API documentation. 
    https://developer.ebay.com/docs

GitHub, Inc. (2024). DiscordGo: Go bindings for Discord. 
    https://github.com/bwmarrin/discordgo

The Go Authors. (2024). The Go programming language documentation. 
    https://go.dev/doc/

The Go Authors. (2024). Effective Go. https://go.dev/doc/effective_go

Discord Inc. (2024). Discord Developer Portal. 
    https://discord.com/developers/docs

Calçado, P. (2024). Awesome Go: A curated list of awesome Go frameworks. 
    https://github.com/avelino/awesome-go
```

---

## 🎯 PRESENTATION TIPS

### Recording Setup:
- **Screen resolution:** 1920x1080 for best quality
- **Close unnecessary browser tabs** to avoid notifications
- **Use a script timer** - practice to stay within 10-15 minutes
- **Record audio separately** if needed for better quality

### Code Display:
- **Increase font size** to at least 16-18pt for readability
- **Use syntax highlighting** in your editor
- **Zoom in** when showing specific code blocks
- **Highlight important lines** with comments or cursor

### Visual Transitions:
- Use clear section transitions ("Now let's move to...")
- Keep slides simple with bullet points
- Show actual application/code more than slides

### Common Pitfalls to Avoid:
- ❌ Reading slides word-for-word
- ❌ Spending too long on any one section
- ❌ Getting stuck on technical issues during demo
- ❌ Using jargon without explanation

### Time Management:
- Introduction: **1 min**
- Problem Statement: **1.5 min**
- Technology Overview: **2 min**
- Go vs Python: **2 min**
- Code Walkthrough: **3-4 min**
- Live Demo: **4-5 min**
- Resources: **1 min**
- AI Tools (Optional): **1.5 min**
- Conclusion: **1 min**
- **TOTAL: 12-14 minutes** (leaving buffer time)

### Practice Run Checklist:
- [ ] Record a practice run
- [ ] Time each section
- [ ] Check audio levels
- [ ] Verify screen is readable when recorded
- [ ] Test any live demos work smoothly
- [ ] Have backup screenshots if live demo fails
- [ ] Review for clarity and pacing

---

## 🎤 Q&A PREPARATION

### Likely Questions to Prepare For:

**Technical:**
- "Why Go instead of Python or JavaScript?"
- "How do you handle rate limiting with eBay's API?"
- "What happens if the bot crashes?"
- "How do you secure your API keys?"

**Project-Specific:**
- "How long did this take to build?"
- "What was the hardest part?"
- "Could this scale to multiple eBay accounts?"
- "What would you add if you had more time?"

**Learning:**
- "How hard was it to learn Go?"
- "Would you recommend Go for other projects?"
- "How did you debug issues with the eBay API?"

### Have Ready:
- Link to GitHub repository (if making it public)
- Example of compiled binary size comparison (Go vs Python equivalent)
- Performance metrics if you have them (response times, etc.)

---

Good luck with your presentation! 🚀
